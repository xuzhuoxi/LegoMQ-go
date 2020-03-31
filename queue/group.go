package queue

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
	"github.com/xuzhuoxi/infra-go/lang/collectionx"
	"strconv"
	"sync"
)

var (
	ErrQueueIdUnknown  = errors.New("MessageQueueGroup: QueueId Unknown. ")
	ErrQueueIdExists   = errors.New("MessageQueueGroup: QueueId Exists. ")
	ErrQueueNil        = errors.New("MessageQueueGroup: Queue is nil. ")
	ErrQueueIndexRange = errors.New("MessageQueueGroup: Index out of range. ")
)

// 消息队列组
// 配置接口
type IMessageQueueGroupConfig interface {
	// 队列数量
	QueueSize() int
	// 队列标识信息
	QueueIds() []string
	// 创建一个队列，id使用默认规则创建
	// err:
	// 		ErrQueueModeUnregister:	QueueMode未注册
	CreateQueue(mode QueueMode, size int) (queue IMessageContextQueue, err error)
	// 创建一个队列，id使用默认规则创建
	// err:
	// 		ErrQueueModeUnregister:	QueueMode未注册
	CreateQueues(modes []QueueMode, size int) (queues []IMessageContextQueue, err []error)
	// 加入一个队列
	// err:
	// 		ErrQueueNil:		queue为nil时
	//		ErrQueueIdExists:	QueueId重复时
	AddQueue(queue IMessageContextQueue) error
	// 加入多个队列
	// count: 成功加入的队列数量
	// err1:
	// 		ErrQueueNil:		queues中单元为nil时
	// err2:
	//		ErrQueueIdExists:	QueueId重复时
	AddQueues(queues []IMessageContextQueue) (count int, failArr []IMessageContextQueue, err []error)
	// 移除一个队列
	// queue: 返回被移除的队列
	// err:
	//		ErrQueueIdUnknown:	QueueId不存在时
	RemoveQueue(queueId string) (queue IMessageContextQueue, err error)
	// 移除多个队列
	// queues: 返回被移除的队列数组
	// err:
	//		ErrQueueIdUnknown:	QueueId不存在时
	RemoveQueues(queueIdArr []string) (queues []IMessageContextQueue, err []error)
	// 替换一个队列
	// 根据QueueId进行替换，如果找不到相同QueueId，直接加入
	// err:
	// 		ErrQueueNil:	queue为nil时
	UpdateQueue(queue IMessageContextQueue) (err error)
	// 替换一个队列
	// 根据QueueId进行替换，如果找不到相同QueueId，直接加入
	// err:
	// 		ErrQueueNil:	queue为nil时
	UpdateQueues(queues []IMessageContextQueue) (err []error)
	// 使用配置初始化队列组，覆盖旧配置
	InitQueueGroup(settings []QueueSetting) (queues []IMessageContextQueue, err error)
	// 路由信息
	RoutingElements() []support.IRoutingTarget
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
	//		ErrQueueIdUnknown:	QueueId不存在时返回
	//		写入失败时返回相应错误
	WriteMessageTo(msg message.IMessageContext, queueId string) error
	// 往指定消息队列写入一个消息
	// count: 成功写入的队列数量
	// failArr: 未写入的队列标识数组
	// err1:
	//		ErrQueueIdUnknown:	queueIdArr中QueueId不存在时返回
	// err2:
	//		写入失败时返回相应错误
	WriteMessageToMulti(msg message.IMessageContext, queueIdArr []string) (count int, failArr []string, err1 error, err2 error)
	// 往指定消息队列写入多个消息
	// 返回说明：
	// count: 成功写入消息的数量
	// err:
	//		ErrQueueIdUnknown:	QueueId不存在时返回
	//		写入失败时返回相应错误
	WriteMessagesTo(msg []message.IMessageContext, queueId string) (count int, err error)
	// 往指定消息队列写入多个消息
	// 返回说明：
	// err1:
	//		ErrQueueIdUnknown:	queueIdArr中QueueId不存在时返回
	// err2:
	//		写入失败时返回相应错误
	WriteMessagesToMulti(msg []message.IMessageContext, queueIdArr []string) (err1 error, err2 error)
}

type IMessageQueueGroup interface {
	IMessageQueueGroupReader
	IMessageQueueGroupWriter
	// 取生成者
	// err:
	// 		ErrQueueIdUnknown:	QueueId不存在
	GetQueue(queueId string) (IMessageContextQueue, error)
	// 取生成者
	// err:
	// 		ErrQueueIndexRange:	index越界
	GetQueueAt(index int) (IMessageContextQueue, error)
	// 配置入口
	Config() IMessageQueueGroupConfig
	// 遍历元素
	ForEachElement(f func(index int, ele IMessageContextQueue) (stop bool))
}

func NewMessageQueueGroup() (config IMessageQueueGroupConfig, group IMessageQueueGroup) {
	rs := &queueGroup{}
	rs.group = *collectionx.NewOrderHashGroup()
	return rs, rs
}

//---------------------

type queueGroup struct {
	group  collectionx.OrderHashGroup
	autoId int
	mu     sync.RWMutex
}

func (g *queueGroup) QueueSize() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.group.Size()
}

func (g *queueGroup) QueueIds() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.group.Ids()
}

func (g *queueGroup) CreateQueue(mode QueueMode, size int) (queue IMessageContextQueue, err error) {
	queue, err = mode.NewContextQueue(size)
	if nil != err {
		return nil, err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	queue.SetId(strconv.Itoa(g.autoId))
	g.autoId += 1
	err = g.group.Add(queue)
	return
}

func (g *queueGroup) CreateQueues(modes []QueueMode, size int) (queues []IMessageContextQueue, err []error) {
	if len(modes) == 0 {
		return nil, nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, mode := range modes {
		queue, e := mode.NewContextQueue(size)
		if nil != e {
			err = append(err, e)
			continue
		}
		queue.SetId(strconv.Itoa(g.autoId))
		g.autoId += 1
		e = g.group.Add(queue)
		if nil != e {
			err = append(err, e)
			continue
		}
		queues = append(queues, queue)
	}
	return
}

func (g *queueGroup) AddQueue(queue IMessageContextQueue) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.group.Add(queue)

}

func (g *queueGroup) AddQueues(queues []IMessageContextQueue) (count int, failArr []IMessageContextQueue, err []error) {
	if len(queues) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for idx, _ := range queues {
		e := g.group.Add(queues[idx])
		if nil != e {
			err = append(err, e)
			failArr = append(failArr, queues[idx])
		} else {
			count += 1
		}
	}
	return
}

func (g *queueGroup) RemoveQueue(queueId string) (queue IMessageContextQueue, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	ele, e := g.group.Remove(queueId)
	if nil != e {
		return nil, e
	}
	return ele.(IMessageContextQueue), nil
}

func (g *queueGroup) RemoveQueues(queueIdArr []string) (queues []IMessageContextQueue, err []error) {
	if len(queueIdArr) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, eleId := range queueIdArr {
		ele, e := g.group.Remove(eleId)
		if nil != e {
			err = append(err, e)
		} else {
			queues = append(queues, ele.(IMessageContextQueue))
		}
	}
	return
}

func (g *queueGroup) UpdateQueue(queue IMessageContextQueue) (err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	_, err = g.group.Update(queue)
	return
}

func (g *queueGroup) UpdateQueues(queues []IMessageContextQueue) (err []error) {
	if len(queues) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for idx, _ := range queues {
		_, e := g.group.Update(queues[idx])
		if nil != e {
			err = append(err, e)
		}
	}
	return
}

func (g *queueGroup) InitQueueGroup(settings []QueueSetting) (queues []IMessageContextQueue, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	group := *collectionx.NewOrderHashGroup()
	if len(settings) == 0 {
		return nil, nil
	}
	for idx, _ := range settings {
		queue, err := NewContextQueue(settings[idx])
		err = group.Add(queue)
		if nil != err {
			return nil, err
		}
		queues = append(queues, queue)
	}
	g.group = group
	return queues, nil
}

func (g *queueGroup) RoutingElements() []support.IRoutingTarget {
	g.mu.RLock()
	defer g.mu.RUnlock()
	rs := make([]support.IRoutingTarget, 0, g.QueueSize())
	g.group.ForEachElement(func(_ int, ele collectionx.IOrderHashElement) (stop bool) {
		rs = append(rs, ele.(support.IRoutingTarget))
		return false
	})
	return rs
}

//-------------------------

func (g *queueGroup) ReadMessageFrom(queueId string) (msg message.IMessageContext, err error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if q, err := g.getQueue(queueId); nil == err {
		return q.ReadContext()
	}
	return nil, ErrQueueIdUnknown
}

func (g *queueGroup) ReadMessagesFrom(queueId string, count int) (msg []message.IMessageContext, err error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if q, err := g.getQueue(queueId); nil == err {
		return q.ReadContexts(count)
	}
	return nil, ErrQueueIdUnknown
}

func (g *queueGroup) ReadMessagesTo(queueId string, msg []message.IMessageContext) (count int, err error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if q, err := g.getQueue(queueId); nil == err {
		return q.ReadContextsTo(msg)
	}
	return -1, ErrQueueIdUnknown
}

func (g *queueGroup) WriteMessageTo(msg message.IMessageContext, queueId string) error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if q, err := g.getQueue(queueId); err == nil {
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
	var q IMessageContextQueue
	var err error
	for idx, _ := range queueIdArr {
		q, err = g.getQueue(queueIdArr[idx])
		if nil != err {
			failArr = append(failArr, queueIdArr[idx])
			continue
		}
		err = q.WriteContext(msg)
		if nil != err {
			err2 = err
			failArr = append(failArr, queueIdArr[idx])
			continue
		}
		count++
	}
	return
}

func (g *queueGroup) WriteMessagesTo(msg []message.IMessageContext, queueId string) (count int, err error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var q IMessageContextQueue
	if q, err = g.getQueue(queueId); err == nil {
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
	var err error
	var q IMessageContextQueue
	for idx, _ := range queueIdArr {
		if q, err = g.getQueue(queueIdArr[idx]); nil == err {
			_, err = q.WriteContexts(msg)
			if nil != err {
				err2 = err
			}
		} else {
			err1 = ErrQueueIdUnknown
		}
	}
	return
}

func (g *queueGroup) GetQueue(queueId string) (IMessageContextQueue, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.getQueue(queueId)
}

func (g *queueGroup) GetQueueAt(index int) (IMessageContextQueue, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.getQueueAt(index)
}

func (g *queueGroup) Config() IMessageQueueGroupConfig {
	return g
}

func (g *queueGroup) ForEachElement(f func(index int, ele IMessageContextQueue) (stop bool)) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	g.group.ForEachElement(func(idx int, ele collectionx.IOrderHashElement) (stop bool) {
		return f(idx, ele.(IMessageContextQueue))
	})
}

//-------------------

func (g *queueGroup) getQueue(queueId string) (IMessageContextQueue, error) {
	if ele, ok := g.group.Get(queueId); ok {
		return ele.(IMessageContextQueue), nil
	}
	return nil, ErrQueueIdUnknown
}

func (g *queueGroup) getQueueAt(index int) (IMessageContextQueue, error) {
	if ele, ok := g.group.GetAt(index); ok {
		return ele.(IMessageContextQueue), nil
	}
	return nil, ErrQueueIndexRange
}
