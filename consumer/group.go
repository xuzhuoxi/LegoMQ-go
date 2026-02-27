package consumer

import (
	"errors"
	"strconv"
	"sync"

	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/routing"
	"github.com/xuzhuoxi/infra-go/lang/collectionx"
)

var (
	ErrConsumerIdExists   = errors.New("MessageConsumer: ConsumerId exists. ")
	ErrConsumerIdUnknown  = errors.New("MessageConsumer: ConsumerId unknown. ")
	ErrConsumerIndexRange = errors.New("MessageConsumer: Index out of range. ")
)

// IMessageConsumerGroupConfig
// 消息消费者组
// 配置接口
type IMessageConsumerGroupConfig interface {
	// ConsumerSize
	// 消费者数量
	ConsumerSize() int
	// ConsumerIds
	// 消费者标识信息
	ConsumerIds() []string

	// CreateConsumer
	// 创建一个消费者，id使用默认规则创建
	// err:
	// 		ErrConsumerModeUnregister:	ConsumerMode未注册
	CreateConsumer(mode ConsumerMode) (consumer IMessageConsumer, err error)
	// CreateConsumers
	// 创建一个消费者，id使用默认规则创建
	// err:
	// 		ErrConsumerModeUnregister:	ConsumerMode未注册
	CreateConsumers(modes []ConsumerMode) (consumers []IMessageConsumer, err []error)
	// AddConsumer
	// 加入一个消费者
	// err:
	// 		ErrConsumerNil:			Consumer为nil
	//		ErrConsumerIdExists:	ConsumerId重复
	AddConsumer(consumer IMessageConsumer) error
	// AddConsumers
	// 加入多个消费者
	// count: 成功加入的消费者数量
	// err:
	// 		ErrConsumerNil:			Consumer为nil
	//		ErrConsumerIdExists:	ConsumerId重复
	AddConsumers(consumers []IMessageConsumer) (count int, failArr []IMessageConsumer, err []error)
	// RemoveConsumer
	// 移除一个消费者
	// consumer: 返回被移除的消费者
	// err:
	//		ErrConsumerIdUnknown:	ConsumerId不存在
	RemoveConsumer(consumerId string) (consumer IMessageConsumer, err error)
	// RemoveConsumers
	// 移除多个消费者
	// consumers: 返回被移除的消费者数组
	// err:
	//		ErrConsumerIdUnknown:	ConsumerId不存在
	RemoveConsumers(consumerIdArr []string) (consumers []IMessageConsumer, err []error)
	// UpdateConsumer
	// 替换一个消费者
	// 根据ConsumerId进行替换，如果找不到相同ConsumerId，直接加入
	// err:
	// 		ErrConsumerNil:			Consumer为nil
	UpdateConsumer(consumer IMessageConsumer) (err error)
	// UpdateConsumers
	// 替换一个消费者
	// 根据ConsumerId进行替换，如果找不到相同ConsumerId，直接加入
	// err:
	// 		ErrConsumerNil:			Consumer为nil
	UpdateConsumers(consumers []IMessageConsumer) (err []error)
	// InitConsumerGroup
	// 使用配置初始化消费者组，覆盖旧配置
	InitConsumerGroup(settings []ConsumerSetting) (consumers []IMessageConsumer, err error)
	// RoutingElements
	// 路由元素
	RoutingElements() []routing.IRoutingTarget
}

// IMessageConsumerGroup
// 消息消费者组
// 操作接口
type IMessageConsumerGroup interface {
	// Config
	// 配置入口
	Config() IMessageConsumerGroupConfig
	// GetConsumer
	// 取生成者
	// err:
	// 		ErrConsumerIdUnknown:	ConsumerId不存在
	GetConsumer(producerId string) (IMessageConsumer, error)
	// GetConsumerAt
	// 取生成者
	// err:
	// 		ErrConsumerIndexRange:	index越界
	GetConsumerAt(index int) (IMessageConsumer, error)
	// ConsumeMessage
	// 被动接收一个消息并进行处理
	// err:
	// 		ErrConsumerMessageNil: msg=nil时
	// 		ErrConsumerIdUnknown: ConsumerId不存在
	ConsumeMessage(msg message.IMessageContext, consumerId string) error
	// ConsumeMessageMulti
	// 被动接收一个消息并进行处理
	// err:
	// 		ErrConsumerMessageNil: msg=nil时
	// 		ErrConsumerIdUnknown: ConsumerId不存在
	ConsumeMessageMulti(msg message.IMessageContext, consumerIds []string) (err []error)
	// ConsumeMessages
	// 被动接收多个消息并进行处理
	// err:
	//		ErrConsumerMessagesEmpty: msg长度为0
	// 		ErrConsumerIdUnknown: ConsumerId不存在
	//		其它错误
	ConsumeMessages(msg []message.IMessageContext, consumerId string) error
	// ForEachElement
	// 遍历元素
	ForEachElement(f func(index int, ele IMessageConsumer) (stop bool))
}

// NewMessageConsumerGroup
// 实例化消息消费者组
// config: 	配置接口
// group: 	操作接口
func NewMessageConsumerGroup() (config IMessageConsumerGroupConfig, group IMessageConsumerGroup) {
	rs := &consumerGroup{}
	rs.group = *collectionx.NewOrderHashGroup()
	return rs, rs
}

//---------------------

type consumerGroup struct {
	group  collectionx.OrderHashGroup
	autoId int
	mu     sync.RWMutex
}

func (g *consumerGroup) Config() IMessageConsumerGroupConfig {
	return g
}

func (g *consumerGroup) ConsumerSize() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.group.Size()
}

func (g *consumerGroup) ConsumerIds() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.group.Ids()
}

func (g *consumerGroup) CreateConsumer(mode ConsumerMode) (consumer IMessageConsumer, err error) {
	consumer, err = mode.NewMessageConsumer()
	if nil != err {
		return nil, err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	consumer.SetId(strconv.Itoa(g.autoId))
	g.autoId += 1
	err = g.group.Add(consumer)
	return
}

func (g *consumerGroup) CreateConsumers(modes []ConsumerMode) (consumers []IMessageConsumer, err []error) {
	if len(modes) == 0 {
		return
	}
	for _, mode := range modes {
		consumer, e := mode.NewMessageConsumer()
		if nil != e {
			err = append(err, e)
			continue
		}
		consumer.SetId(strconv.Itoa(g.autoId))
		g.autoId += 1
		e = g.group.Add(consumer)
		if nil != e {
			err = append(err, e)
			continue
		}
		consumers = append(consumers, consumer)
	}
	return
}

func (g *consumerGroup) AddConsumer(consumer IMessageConsumer) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.group.Add(consumer)
}

func (g *consumerGroup) AddConsumers(consumers []IMessageConsumer) (count int, failArr []IMessageConsumer, err []error) {
	if len(consumers) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for idx, _ := range consumers {
		e := g.group.Add(consumers[idx])
		if nil != e {
			err = append(err, e)
			failArr = append(failArr, consumers[idx])
		} else {
			count += 1
		}
	}
	return
}

func (g *consumerGroup) RemoveConsumer(consumerId string) (consumer IMessageConsumer, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	ele, e := g.group.Remove(consumerId)
	if nil != e {
		return nil, e
	}
	return ele.(IMessageConsumer), nil
}

func (g *consumerGroup) RemoveConsumers(consumerIdArr []string) (consumers []IMessageConsumer, err []error) {
	if len(consumerIdArr) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, eleId := range consumerIdArr {
		ele, e := g.group.Remove(eleId)
		if nil != e {
			err = append(err, e)
		} else {
			consumers = append(consumers, ele.(IMessageConsumer))
		}
	}
	return
}

func (g *consumerGroup) UpdateConsumer(consumer IMessageConsumer) (err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	_, err = g.group.Update(consumer)
	return
}

func (g *consumerGroup) UpdateConsumers(consumers []IMessageConsumer) (err []error) {
	if len(consumers) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for idx, _ := range consumers {
		_, e := g.group.Update(consumers[idx])
		if nil != e {
			err = append(err, e)
		}
	}
	return
}

func (g *consumerGroup) InitConsumerGroup(settings []ConsumerSetting) (consumers []IMessageConsumer, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	group := *collectionx.NewOrderHashGroup()
	if len(settings) == 0 {
		return nil, nil
	}
	for idx, _ := range settings {
		consumer, err := settings[idx].NewMessageConsumer()
		if nil != err {
			return nil, err
		}
		err = group.Add(consumer)
		if nil != err {
			return nil, err
		}
		consumers = append(consumers, consumer)
	}
	g.group = group
	return consumers, nil
}

func (g *consumerGroup) RoutingElements() []routing.IRoutingTarget {
	g.mu.RLock()
	defer g.mu.RUnlock()
	rs := make([]routing.IRoutingTarget, 0, g.ConsumerSize())
	g.group.ForEachElement(func(_ int, ele collectionx.IOrderHashElement) (stop bool) {
		rs = append(rs, ele.(routing.IRoutingTarget))
		return false
	})
	return rs
}

//----------------------

func (g *consumerGroup) GetConsumerAt(index int) (IMessageConsumer, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if ele, ok := g.group.GetAt(index); ok {
		return ele.(IMessageConsumer), nil
	} else {
		return nil, ErrConsumerIdUnknown
	}
}

func (g *consumerGroup) GetConsumer(producerId string) (IMessageConsumer, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if ele, ok := g.group.Get(producerId); ok {
		return ele.(IMessageConsumer), nil
	} else {
		return nil, ErrConsumerIdUnknown
	}
}

func (g *consumerGroup) ConsumeMessage(msg message.IMessageContext, consumerId string) error {
	if nil == msg {
		return ErrConsumerMessageNil
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	ele, ok := g.group.Get(consumerId)
	if !ok {
		return ErrConsumerIdUnknown
	}
	return ele.(IMessageConsumer).ConsumeMessage(msg)
}

func (g *consumerGroup) ConsumeMessageMulti(msg message.IMessageContext, consumerIds []string) (err []error) {
	if 0 == len(consumerIds) {
		return nil
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	for idx, _ := range consumerIds {
		ele, ok := g.group.Get(consumerIds[idx])
		if !ok {
			err = append(err, ErrConsumerIdUnknown)
			continue
		}
		e := ele.(IMessageConsumer).ConsumeMessage(msg)
		if nil != e {
			err = append(err, e)
		}
	}
	return
}

func (g *consumerGroup) ConsumeMessages(msg []message.IMessageContext, consumerId string) error {
	if 0 == len(msg) {
		return ErrConsumerMessagesEmpty
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	ele, ok := g.group.Get(consumerId)
	if !ok {
		return ErrConsumerIdUnknown
	}
	return ele.(IMessageConsumer).ConsumeMessages(msg)
}

func (g *consumerGroup) ForEachElement(f func(index int, ele IMessageConsumer) (stop bool)) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	g.group.ForEachElement(func(idx int, ele collectionx.IOrderHashElement) (stop bool) {
		return f(idx, ele.(IMessageConsumer))
	})
}
