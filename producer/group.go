package producer

import (
	"errors"
	"strconv"
	"sync"

	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/lang/collectionx"
)

var (
	ErrProducerIdExists   = errors.New("MessageProducer: ProducerId exists. ")
	ErrProducerIdUnknown  = errors.New("MessageProducer: ProducerId unknown. ")
	ErrProducerIndexRange = errors.New("MessageProducer: Index out of range. ")
)

type FuncOnMessageProduced func(msg message.IMessageContext, locateId string)
type FuncOnMessagesProduced func(msg []message.IMessageContext, locateId string)

// IMessageProducerGroupConfig
// 消息生产者组
// 配置接口
type IMessageProducerGroupConfig interface {
	// ProducerSize
	// 生产者数量
	ProducerSize() int
	// ProducerIds
	// 生产者标识信息
	ProducerIds() []string

	// CreateProducer
	// 创建一个生产者并加入组，id使用默认规则创建
	// err:
	// 		ErrProducerModeUnregister:	ProducerMode未注册
	CreateProducer(mode ProducerMode) (producer IMessageProducer, err error)
	// CreateProducers
	// 创建多个生产者并加入组，id使用默认规则创建
	// err:
	// 		ErrProducerModeUnregister:	ProducerMode未注册
	CreateProducers(modes []ProducerMode) (producers []IMessageProducer, err []error)
	// AddProducer
	// 加入一个生产者
	// err:
	// 		ErrProducerNil:			Producer为nil
	//		ErrProducerIdExists:	ProducerId重复
	AddProducer(producer IMessageProducer) error
	// AddProducers
	// 加入多个生产者
	// count: 成功加入的生产者数量
	// err:
	// 		ErrProducerNil:			Producer为nil
	//		ErrProducerIdExists:	ProducerId重复
	AddProducers(producers []IMessageProducer) (count int, failArr []IMessageProducer, err []error)
	// RemoveProducer
	// 移除一个生产者
	// producer: 返回被移除的生产者
	// err:
	//		ErrProducerIdUnknown:	ProducerId不存在
	RemoveProducer(producerId string) (producer IMessageProducer, err error)
	// RemoveProducers
	// 移除多个生产者
	// producers: 返回被移除的生产者数组
	// err:
	//		ErrProducerIdUnknown:	ProducerId不存在
	RemoveProducers(producerIdArr []string) (producers []IMessageProducer, err []error)
	// ClearProducers
	// 清空
	ClearProducers()
	// UpdateProducer
	// 替换一个生产者
	// 根据ProducerId进行替换，如果找不到相同ProducerId，直接加入
	// err:
	// 		ErrProducerNil:			Producer为nil
	UpdateProducer(producer IMessageProducer) (err error)
	// UpdateProducers
	// 替换一个生产者
	// 根据ProducerId进行替换，如果找不到相同ProducerId，直接加入
	// err:
	// 		ErrProducerNil:			Producer为nil
	UpdateProducers(producers []IMessageProducer) (err []error)
	// InitProducerGroup
	// 使用配置初始化生产者组，覆盖旧配置
	InitProducerGroup(settings []ProducerSetting) (producers []IMessageProducer, err error)
	// SetProducedFunc
	// 设置响应行为
	SetProducedFunc(f FuncOnMessageProduced, fs FuncOnMessagesProduced)
}

// IMessageProducerGroup
// 消息生产者组
// 操作接口
type IMessageProducerGroup interface {
	// Config
	// 配置入口
	Config() IMessageProducerGroupConfig
	// GetProducer
	// 取生成者
	// err:
	// 		ErrProducerIdUnknown:	ProducerId不存在
	GetProducer(producerId string) (IMessageProducer, error)
	// GetProducerAt
	// 取生成者
	// err:
	// 		ErrProducerIndexRange:	index越界
	GetProducerAt(index int) (IMessageProducer, error)
	// NotifyMessageProduced
	// 生产消息
	// 抛出事件 EventMessageOnProducer
	// err:
	//		ErrProducerMessageNil:	msg为nil时
	NotifyMessageProduced(msg message.IMessageContext, producerId string) error
	// NotifyMessagesProduced
	// 生产消息
	// 抛出事件 EventMultiMessageOnProducer
	// err:
	//		ErrProducerMessagesEmpty:	msg长度为0时
	// 		ErrProducerIdUnknown: 		ProducerId不存在
	//		其它错误
	NotifyMessagesProduced(msg []message.IMessageContext, producerId string) error
}

// NewMessageProducerGroup
// 实例化消息生产组
// config: 	配置接口
// group: 	操作接口
func NewMessageProducerGroup() (config IMessageProducerGroupConfig, group IMessageProducerGroup) {
	rs := &producerGroup{}
	rs.group = *collectionx.NewOrderHashGroup()
	return rs, rs
}

//---------------------

type producerGroup struct {
	group            collectionx.OrderHashGroup
	autoId           int
	funcMsgProduced  FuncOnMessageProduced
	funcMsgsProduced FuncOnMessagesProduced

	mu sync.RWMutex
}

func (g *producerGroup) Config() IMessageProducerGroupConfig {
	return g
}

func (g *producerGroup) ProducerSize() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.group.Size()
}

func (g *producerGroup) ProducerIds() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.group.Ids()
}

func (g *producerGroup) CreateProducer(mode ProducerMode) (producer IMessageProducer, err error) {
	producer, err = mode.NewMessageProducer()
	if nil != err {
		return nil, err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	producer.SetId(strconv.Itoa(g.autoId))
	g.autoId += 1
	err = g.group.Add(producer)
	if nil == err {
		g.addListeners(producer)
	}
	return
}

func (g *producerGroup) CreateProducers(modes []ProducerMode) (producers []IMessageProducer, err []error) {
	if len(modes) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, mode := range modes {
		producer, e := mode.NewMessageProducer()
		if nil != e {
			err = append(err, e)
			continue
		}
		producer.SetId(strconv.Itoa(g.autoId))
		g.autoId += 1
		e = g.group.Add(producer)
		if nil != e {
			err = append(err, e)
			continue
		}
		g.addListeners(producer)
		producers = append(producers, producer)
	}
	return
}

func (g *producerGroup) AddProducer(producer IMessageProducer) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	err := g.group.Add(producer)
	if nil != err {
		return err
	}
	g.addListeners(producer)
	return nil
}

func (g *producerGroup) AddProducers(producers []IMessageProducer) (count int, failArr []IMessageProducer, err []error) {
	if len(producers) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for idx, _ := range producers {
		e := g.group.Add(producers[idx])
		if nil != e {
			err = append(err, e)
			failArr = append(failArr, producers[idx])
		} else {
			g.addListeners(producers[idx])
			count += 1
		}
	}
	return
}

func (g *producerGroup) RemoveProducer(producerId string) (producer IMessageProducer, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	ele, e := g.group.Remove(producerId)
	if nil != e {
		return nil, e
	}
	producer = ele.(IMessageProducer)
	g.removeListeners(producer)
	return
}

func (g *producerGroup) RemoveProducers(producerIdArr []string) (producers []IMessageProducer, err []error) {
	if len(producerIdArr) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, eleId := range producerIdArr {
		ele, e := g.group.Remove(eleId)
		if nil != e {
			err = append(err, e)
		} else {
			producer := ele.(IMessageProducer)
			g.removeListeners(producer)
			producers = append(producers, producer)
		}
	}
	return
}

func (g *producerGroup) ClearProducers() {
	g.mu.Lock()
	defer g.mu.Unlock()
	removes := g.group.RemoveAll()
	if len(removes) == 0 {
		return
	}
	for idx, _ := range removes {
		g.removeListeners(removes[idx].(IMessageProducer))
	}
}

func (g *producerGroup) UpdateProducer(producer IMessageProducer) (err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	replaced, e := g.group.Update(producer)
	if e != nil {
		return e
	}
	g.removeListeners(replaced.(IMessageProducer))
	g.addListeners(producer)
	return nil
}

func (g *producerGroup) UpdateProducers(producers []IMessageProducer) (err []error) {
	if len(producers) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for idx, _ := range producers {
		r, e := g.group.Update(producers[idx])
		if nil != e {
			err = append(err, e)
		} else {
			g.removeListeners(r.(IMessageProducer))
			g.addListeners(producers[idx])
		}
	}
	return
}

func (g *producerGroup) InitProducerGroup(settings []ProducerSetting) (producers []IMessageProducer, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	group := *collectionx.NewOrderHashGroup()
	if len(settings) == 0 {
		return nil, nil
	}
	for idx, _ := range settings {
		producer, err := settings[idx].NewMessageProducer()
		if nil != err {
			return nil, err
		}
		err = group.Add(producer)
		if nil != err {
			return nil, err
		}
		g.addListeners(producer)
		producers = append(producers, producer)
	}
	g.group = group
	return producers, nil
}

func (g *producerGroup) SetProducedFunc(f FuncOnMessageProduced, fs FuncOnMessagesProduced) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.funcMsgProduced, g.funcMsgsProduced = f, fs
}

func (g *producerGroup) addListeners(producer IMessageProducer) {
	producer.AddEventListener(EventMessageOnProducer, g.onProduced)
	producer.AddEventListener(EventMultiMessageOnProducer, g.onMultiProduced)
}

func (g *producerGroup) removeListeners(producer IMessageProducer) {
	producer.RemoveEventListener(EventMultiMessageOnProducer, g.onMultiProduced)
	producer.RemoveEventListener(EventMessageOnProducer, g.onProduced)
}

func (g *producerGroup) onProduced(evt *eventx.EventData) {
	if nil != g.funcMsgProduced {
		ct := evt.CurrentTarget.(IMessageProducer)
		msg := evt.Data.(message.IMessageContext)
		g.funcMsgProduced(msg, ct.LocateId())
	}
}

func (g *producerGroup) onMultiProduced(evt *eventx.EventData) {
	if nil != g.funcMsgsProduced {
		ct := evt.CurrentTarget.(IMessageProducer)
		msg := evt.Data.([]message.IMessageContext)
		g.funcMsgsProduced(msg, ct.LocateId())
	}
}

//----------------------

func (g *producerGroup) GetProducer(producerId string) (IMessageProducer, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if ele, ok := g.group.Get(producerId); ok {
		return ele.(IMessageProducer), nil
	} else {
		return nil, ErrProducerIdUnknown
	}
}

func (g *producerGroup) GetProducerAt(index int) (IMessageProducer, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if ele, ok := g.group.GetAt(index); ok {
		return ele.(IMessageProducer), nil
	} else {
		return nil, ErrProducerIndexRange
	}
}

func (g *producerGroup) NotifyMessageProduced(msg message.IMessageContext, producerId string) error {
	if nil == msg {
		return ErrProducerMessageNil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	ele, ok := g.group.Get(producerId)
	if !ok {
		return ErrProducerIdUnknown
	}
	return ele.(IMessageProducer).NotifyMessageProduced(msg)
}

func (g *producerGroup) NotifyMessagesProduced(msg []message.IMessageContext, producerId string) error {
	if 0 == len(msg) {
		return ErrProducerMessagesEmpty
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	ele, ok := g.group.Get(producerId)
	if !ok {
		return ErrProducerIdUnknown
	}
	return ele.(IMessageProducer).NotifyMessagesProduced(msg)
}
