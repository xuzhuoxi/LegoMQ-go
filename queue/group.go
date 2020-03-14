package queue

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"strconv"
	"sync"
)

var (
	ErrQueueIdUnknown = errors.New("QueueId Unknown. ")
	ErrQueueIdExists  = errors.New("QueueId Exists. ")
	ErrQueueNil       = errors.New("Queue is nil. ")
)

type QueueGroupSetting struct {
	Queues []QueueSetting
}

// 消息队列组
// 配置接口
type IMessageQueueGroupConfig interface {
	// 队列数量
	QueueSize() int
	// 队列标识信息
	QueueIds() []string
	// 创建一个队列，id使用默认规则创建
	// err:
	// 		QueueMode未注册时返回 ErrQueueModeRegister
	CreateQueue(mode QueueMode, size int) (queue IMessageContextQueue, err error)
	// 创建一个队列，id使用默认规则创建
	// err:
	// 		QueueMode未注册时返回 ErrQueueModeRegister
	CreateQueues(modes []QueueMode, size int) (queues []IMessageContextQueue, err error)
	// 加入一个队列
	// err:
	// 		queue为nil时返回ErrQueueNil
	//		QueueId重复时返回 ErrQueueIdExists
	AddQueue(queue IMessageContextQueue) error
	// 加入多个队列
	// count: 成功加入的队列数量
	// err1:
	// 		queues中单元为nil时返回ErrQueueNil
	// err2:
	//		QueueId重复时返回 ErrQueueIdExists
	AddQueues(queues []IMessageContextQueue) (count int, failArr []IMessageContextQueue, err1 error, err2 error)
	// 移除一个队列
	// queue: 返回被移除的队列
	// err:
	//		QueueId不存在时返回 ErrQueueIdUnknown
	RemoveQueue(queueId string) (queue IMessageContextQueue, err error)
	// 移除多个队列
	// queues: 返回被移除的队列数组
	// err:
	//		QueueId不存在时返回 ErrQueueIdUnknown
	RemoveQueues(queueIdArr []string) (queues []IMessageContextQueue, err error)
	// 替换一个队列
	// 根据QueueId进行替换，如果找不到相同QueueId，直接加入
	// err:
	// 		queue为nil时返回ErrQueueNil
	UpdateQueue(queue IMessageContextQueue) (err error)
	// 替换一个队列
	// 根据QueueId进行替换，如果找不到相同QueueId，直接加入
	// err:
	// 		queue为nil时返回ErrQueueNil
	UpdateQueues(queues []IMessageContextQueue) (err error)
	// 使用配置初始化队列组，覆盖旧配置
	InitQueueGroup(settings QueueGroupSetting) error
}

// 消息队列组
// 读取接口
type IMessageQueueGroupReader interface {
	// 从某个队列中读取一条消息
	// err:
	//		QueueId不存在时返回 ErrQueueIdUnknown
	ReadMessageFrom(queueId string) (msg message.IMessageContext, err error)
	// 从某个队列中读取多条消息
	// err:
	//		QueueId不存在时返回 ErrQueueIdUnknown
	ReadMessagesFrom(queueId string, count int) (msg []message.IMessageContext, err error)
	// 从某个队列中读取多条消息并放入到指定数组中
	// count: 成功读取的信息数量
	// err:
	//		QueueId不存在时返回 ErrQueueIdUnknown
	ReadMessagesTo(queueId string, msg []message.IMessageContext) (count int, err error)
}

// 消息队列组
// 写入接口
type IMessageQueueGroupWriter interface {
	// 往指定消息队列写入一个消息
	// err:
	//		QueueId不存在时返回 ErrQueueIdUnknown
	//		写入失败时返回相应错误
	WriteMessageTo(msg message.IMessageContext, queueId string) error
	// 往指定消息队列写入一个消息
	// count: 成功写入的队列数量
	// failArr: 未写入的队列标识数组
	// err1:
	//		QueueId不存在时返回 ErrQueueIdUnknown
	// err2:
	//		写入失败时返回相应错误
	WriteMessageToMulti(msg message.IMessageContext, queueIdArr []string) (count int, failArr []string, err1 error, err2 error)
	// 往指定消息队列写入多个消息
	// 返回说明：
	// count: 成功写入消息的数量
	// err:
	//		QueueId不存在时返回 ErrQueueIdUnknown
	//		写入失败时返回相应错误
	WriteMessagesTo(msg []message.IMessageContext, queueId string) (count int, err error)
	// 往指定消息队列写入多个消息
	// 返回说明：
	// err1:
	//		QueueId不存在时返回 ErrQueueIdUnknown
	// err2:
	//		写入失败时返回相应错误
	WriteMessagesToMulti(msg []message.IMessageContext, queueIdArr []string) (err1 error, err2 error)
}

type IMessageQueueGroup interface {
	IMessageQueueGroupReader
	IMessageQueueGroupWriter
}

func NewMessageQueueGroup() IMessageQueueGroup {
	rs := &queueGroup{}
	rs.queueMap = make(map[string]IMessageContextQueue)
	rs.queues = make([]IMessageContextQueue, 0, 32)
	return rs
}

//---------------------

type queueGroup struct {
	queues   []IMessageContextQueue
	queueMap map[string]IMessageContextQueue
	mu       sync.RWMutex
}

func (g *queueGroup) QueueSize() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.queues)
}

func (g *queueGroup) QueueIds() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	ln := len(g.queues)
	rs := make([]string, ln, ln)
	for index, q := range g.queues {
		rs[index] = q.Id()
	}
	return rs
}

func (g *queueGroup) CreateQueue(mode QueueMode, size int) (queue IMessageContextQueue, err error) {
	queue, err = mode.NewContextQueue(size)
	if nil != err {
		return nil, err
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	index := len(g.queues)
	qId := strconv.Itoa(index)
	queue.SetId(qId)

	g.addQueue(queue)
	return
}

func (g *queueGroup) CreateQueues(modes []QueueMode, size int) (queues []IMessageContextQueue, err error) {
	if len(modes) == 0 {
		return nil, nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, mode := range modes {
		queue, e := mode.NewContextQueue(size)
		if nil != e {
			err = e
			continue
		}
		index := len(g.queues)
		qId := strconv.Itoa(index)
		queue.SetId(qId)
		g.addQueue(queue)
		queues = append(queues, queue)
	}
	return
}

func (g *queueGroup) AddQueue(queue IMessageContextQueue) error {
	if nil != queue {
		return ErrQueueNil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.isIdExists(queue.Id()) {
		return ErrQueueIdExists
	}
	g.addQueue(queue)
	return nil

}

func (g *queueGroup) AddQueues(queues []IMessageContextQueue) (count int, failArr []IMessageContextQueue, err1 error, err2 error) {
	if len(queues) == 0 {
		return 0, nil, nil, nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, queue := range queues {
		if nil == queue {
			err1 = ErrQueueNil
			continue
		}
		if g.isIdExists(queue.Id()) {
			err2 = ErrQueueIdExists
			failArr = append(failArr, queue)
			continue
		}
		g.addQueue(queue)
		count += 1
	}
	return
}

func (g *queueGroup) RemoveQueue(queueId string) (queue IMessageContextQueue, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.isIdExists(queueId) {
		return nil, ErrQueueIdUnknown
	}
	queue = g.removeQueueBy(queueId)
	index := g.findQueueIndex(queueId)
	g.removeQueueAt(index)
	return
}

func (g *queueGroup) RemoveQueues(queueIdArr []string) (queues []IMessageContextQueue, err error) {
	if len(queues) == 0 {
		return nil, nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, qId := range queueIdArr {
		if !g.isIdExists(qId) {
			err = ErrQueueIdUnknown
			continue
		}
		queue := g.removeQueueBy(qId)
		index := g.findQueueIndex(qId)
		g.removeQueueAt(index)
		queues = append(queues, queue)
	}
	return
}

func (g *queueGroup) UpdateQueue(queue IMessageContextQueue) (err error) {
	if nil != queue {
		return ErrQueueNil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	qId := queue.Id()
	if g.isIdExists(qId) {
		index := g.findQueueIndex(qId)
		g.queues[index] = queue
	} else {
		g.queues = append(g.queues, queue)
	}
	g.queueMap[qId] = queue
	return nil
}

func (g *queueGroup) UpdateQueues(queues []IMessageContextQueue) (err error) {
	if len(queues) == 0 {
		return nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, queue := range queues {
		if nil != queue {
			err = ErrQueueNil
		}
		qId := queue.Id()
		if g.isIdExists(qId) {
			index := g.findQueueIndex(qId)
			g.queues[index] = queue
		} else {
			g.queues = append(g.queues, queue)
		}
		g.queueMap[qId] = queue
	}
	return
}

func (g *queueGroup) InitQueueGroup(settings QueueGroupSetting) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	var queues []IMessageContextQueue
	var queueMap = make(map[string]IMessageContextQueue)
	for index, _ := range settings.Queues {
		queue, err := NewContextQueue(settings.Queues[index])
		if nil != err {
			return nil
		}
		queues = append(queues, queue)
		queueMap[queue.Id()] = queue
	}
	g.queues = queues
	g.queueMap = queueMap
	return nil
}

func (g *queueGroup) addQueue(queue IMessageContextQueue) {
	g.queues = append(g.queues, queue)
	g.queueMap[queue.Id()] = queue
}

func (g *queueGroup) removeQueueAt(index int) (queue IMessageContextQueue) {
	if index >= 0 && index < len(g.queues) {
		queue = g.queues[index]
		g.queues = append(g.queues[:index], g.queues[index:]...)
	}
	return
}

func (g *queueGroup) removeQueueBy(queueId string) (queue IMessageContextQueue) {
	if q, ok := g.queueMap[queueId]; ok {
		delete(g.queueMap, queueId)
		return q
	} else {
		return nil
	}
}

func (g *queueGroup) findQueueIndex(queueId string) int {
	for index, q := range g.queues {
		if q.Id() == queueId {
			return index
		}
	}
	return -1
}

func (g *queueGroup) isIdExists(queueId string) bool {
	if _, ok := g.queueMap[queueId]; ok {
		return true
	}
	return false
}

//-------------------------

func (g *queueGroup) ReadMessageFrom(queueId string) (msg message.IMessageContext, err error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if q, ok := g.queueMap[queueId]; ok {
		return q.ReadContext()
	}
	return nil, ErrQueueIdUnknown
}

func (g *queueGroup) ReadMessagesFrom(queueId string, count int) (msg []message.IMessageContext, err error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if q, ok := g.queueMap[queueId]; ok {
		return q.ReadContexts(count)
	}
	return nil, ErrQueueIdUnknown
}

func (g *queueGroup) ReadMessagesTo(queueId string, msg []message.IMessageContext) (count int, err error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if q, ok := g.queueMap[queueId]; ok {
		return q.ReadContextsTo(msg)
	}
	return -1, ErrQueueIdUnknown
}

func (g *queueGroup) WriteMessageTo(msg message.IMessageContext, queueId string) error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if q, ok := g.queueMap[queueId]; ok {
		return q.WriteContext(msg)
	}
	return ErrQueueIdUnknown
}

func (g *queueGroup) WriteMessageToMulti(msg message.IMessageContext, queueIdArr []string) (count int, failArr []string, err1 error, err2 error) {
	if len(queueIdArr) == 0 {
		return
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, queueId := range queueIdArr {
		if q, ok := g.queueMap[queueId]; ok {
			err := q.WriteContext(msg)
			if nil != err {
				err2 = err
			} else {
				count += 1
			}
		} else {
			err1 = ErrQueueIdUnknown
		}
	}
	return
}

func (g *queueGroup) WriteMessagesTo(msg []message.IMessageContext, queueId string) (count int, err error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if q, ok := g.queueMap[queueId]; ok {
		return q.WriteContexts(msg)
	}
	return -1, ErrQueueIdUnknown
}

func (g *queueGroup) WriteMessagesToMulti(msg []message.IMessageContext, queueIdArr []string) (err1 error, err2 error) {
	if len(queueIdArr) == 0 {
		return
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, queueId := range queueIdArr {
		if q, ok := g.queueMap[queueId]; ok {
			_, err := q.WriteContexts(msg)
			if nil != err {
				err2 = err
			}
		} else {
			err1 = ErrQueueIdUnknown
		}
	}
	return
}
