package producer

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/lang/collectionx"
	"strconv"
	"sync"
)

var (
	ErrProducerNil        = errors.New("MessageProducer: Producer is nil. ")
	ErrProducerIdExists   = errors.New("MessageProducer: ProducerId exists. ")
	ErrProducerIdUnknown  = errors.New("MessageProducer: ProducerId unknown. ")
	ErrProducerIndexRange = errors.New("MessageProducer: Index out of range. ")
)

type IMessageProducerGroupConfig interface {
	// 生产者数量
	ProducerSize() int
	// 生产者标识信息
	ProducerIds() []string

	// 创建一个生产者，id使用默认规则创建
	// err:
	// 		ErrProducerModeUnregister:	ProducerMode未注册
	CreateProducer(mode ProducerMode) (producer IMessageProducer, err error)
	// 创建一个生产者，id使用默认规则创建
	// err:
	// 		ErrProducerModeUnregister:	ProducerMode未注册
	CreateProducers(modes []ProducerMode) (producers []IMessageProducer, err []error)
	// 加入一个生产者
	// err:
	// 		ErrProducerNil:			Producer为nil
	//		ErrProducerIdExists:	ProducerId重复
	AddProducer(producer IMessageProducer) error
	// 加入多个生产者
	// count: 成功加入的生产者数量
	// err:
	// 		ErrProducerNil:			Producer为nil
	//		ErrProducerIdExists:	ProducerId重复
	AddProducers(producers []IMessageProducer) (count int, failArr []IMessageProducer, err []error)
	// 移除一个生产者
	// producer: 返回被移除的生产者
	// err:
	//		ErrProducerIdUnknown:	ProducerId不存在
	RemoveProducer(producerId string) (producer IMessageProducer, err error)
	// 移除多个生产者
	// producers: 返回被移除的生产者数组
	// err:
	//		ErrProducerIdUnknown:	ProducerId不存在
	RemoveProducers(producerIdArr []string) (producers []IMessageProducer, err []error)
	// 替换一个生产者
	// 根据ProducerId进行替换，如果找不到相同ProducerId，直接加入
	// err:
	// 		ErrProducerNil:			Producer为nil
	UpdateProducer(producer IMessageProducer) (err error)
	// 替换一个生产者
	// 根据ProducerId进行替换，如果找不到相同ProducerId，直接加入
	// err:
	// 		ErrProducerNil:			Producer为nil
	UpdateProducers(producers []IMessageProducer) (err []error)
	// 使用配置初始化生产者组，覆盖旧配置
	InitProducerGroup(settings []ProducerSetting) (producers []IMessageProducer, err error)
}

type IMessageProducerGroup interface {
	// 取生成者
	// err:
	// 		ErrProducerIdUnknown:	ProducerId不存在
	GetProducer(producerId string) (IMessageProducer, error)
	// 取生成者
	// err:
	// 		ErrProducerIndexRange:	index越界
	GetProducerAt(index int) (IMessageProducer, error)
	// 生产消息
	// 抛出事件 EventMessageOnProducer
	// err:
	//		ErrProducerMessageNil:	msg为nil时
	NotifyMessageProduced(msg message.IMessageContext, producerId string) error
	// 生产消息
	// 抛出事件 EventMultiMessageOnProducer
	// err:
	//		ErrProducerMessagesEmpty:	msg长度为0时
	// 		ErrProducerIdUnknown: 		ProducerId不存在
	//		其它错误
	NotifyMessagesProduced(msg []message.IMessageContext, producerId string) error
}

func NewMessageProducerGroup() (config IMessageProducerGroupConfig, group IMessageProducerGroup) {
	rs := &producerGroup{}
	return rs, rs
}

//---------------------

type producerGroup struct {
	group  collectionx.OrderHashGroup
	autoId int
	mu     sync.RWMutex
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
	return
}

func (g *producerGroup) CreateProducers(modes []ProducerMode) (producers []IMessageProducer, err []error) {
	if len(modes) == 0 {
		return
	}
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
		producers = append(producers, producer)
	}
	return
}

func (g *producerGroup) AddProducer(producer IMessageProducer) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.group.Add(producer)
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
	return ele.(IMessageProducer), nil
}

func (g *producerGroup) RemoveProducers(producerIdArr []string) (producers []IMessageProducer, err []error) {
	if len(producers) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, eleId := range producerIdArr {
		ele, e := g.group.Remove(eleId)
		if nil != e {
			err = append(err, e)
		} else {
			producers = append(producers, ele.(IMessageProducer))
		}
	}
	return
}

func (g *producerGroup) UpdateProducer(producer IMessageProducer) (err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.group.Update(producer)
}

func (g *producerGroup) UpdateProducers(producers []IMessageProducer) (err []error) {
	if len(producers) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for idx, _ := range producers {
		e := g.group.Update(producers[idx])
		if nil != e {
			err = append(err, e)
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
		producer, err := settings[idx].Mode.NewMessageProducer()
		if nil != err {
			return nil, err
		}
		producer.SetId(settings[idx].Id)
		err = group.Add(producer)
		if nil != err {
			return nil, err
		}
		producers = append(producers, producer)
	}
	g.group = group
	return producers, nil
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
