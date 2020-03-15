package consumer

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/lang/collectionx"
	"strconv"
	"sync"
)

var (
	ErrConsumerNil       = errors.New("MessageConsumer: Message is nil. ")
	ErrConsumerIdExists  = errors.New("MessageConsumer: Message is nil. ")
	ErrConsumerIdUnknown = errors.New("MessageConsumer: Message is nil. ")
)

type IMessageConsumerGroupConfig interface {
	// 消费者数量
	ConsumerSize() int
	// 消费者标识信息
	ConsumerIds() []string

	// 创建一个消息，id使用默认规则创建
	// err:
	// 		ErrConsumerModeUnregister:	ConsumerMode未注册
	CreateConsumer(mode ConsumerMode) (consumer IMessageConsumer, err error)
	// 创建一个消费者，id使用默认规则创建
	// err:
	// 		ErrConsumerModeUnregister:	ConsumerMode未注册
	CreateConsumers(modes []ConsumerMode) (consumers []IMessageConsumer, err []error)
	// 加入一个消费者
	// err:
	// 		ErrConsumerNil:			Consumer为nil
	//		ErrConsumerIdExists:	ConsumerId重复
	AddConsumer(consumer IMessageConsumer) error
	// 加入多个消费者
	// count: 成功加入的消费者数量
	// err:
	// 		ErrConsumerNil:			Consumer为nil
	//		ErrConsumerIdExists:	ConsumerId重复
	AddConsumers(consumers []IMessageConsumer) (count int, failArr []IMessageConsumer, err []error)
	// 移除一个消费者
	// consumer: 返回被移除的消费者
	// err:
	//		ErrConsumerIdUnknown:	ConsumerId不存在
	RemoveConsumer(consumerId string) (consumer IMessageConsumer, err error)
	// 移除多个消费者
	// consumers: 返回被移除的消费者数组
	// err:
	//		ErrConsumerIdUnknown:	ConsumerId不存在
	RemoveConsumers(consumerIdArr []string) (consumers []IMessageConsumer, err []error)
	// 替换一个消费者
	// 根据ConsumerId进行替换，如果找不到相同ConsumerId，直接加入
	// err:
	// 		ErrConsumerNil:			Consumer为nil
	UpdateConsumer(consumer IMessageConsumer) (err error)
	// 替换一个消费者
	// 根据ConsumerId进行替换，如果找不到相同ConsumerId，直接加入
	// err:
	// 		ErrConsumerNil:			Consumer为nil
	UpdateConsumers(consumers []IMessageConsumer) (err []error)
	// 使用配置初始化消费者组，覆盖旧配置
	InitConsumerGroup(settings []ConsumerSetting) error
}

type IMessageConsumerGroup interface {
	// 被动接收一个消息并进行处理
	// err:
	// 		ErrConsumerMessageNil: msg=nil时
	// 		ErrConsumerIdUnknown: ConsumerId不存在
	ConsumeMessage(msg message.IMessageContext, consumerId string) error
	// 被动接收多个消息并进行处理
	ConsumeMessages(msg []message.IMessageContext, consumerId string) error
}

func NewMessageConsumerGroup() IMessageConsumerGroup {
	rs := &consumerGroup{}
	return rs
}

//---------------------

type consumerGroup struct {
	group  collectionx.OrderHashGroup
	autoId int
	mu     sync.RWMutex
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
	if len(consumers) == 0 {
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
	return g.group.Update(consumer)
}

func (g *consumerGroup) UpdateConsumers(consumers []IMessageConsumer) (err []error) {
	if len(consumers) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for idx, _ := range consumers {
		e := g.group.Update(consumers[idx])
		if nil != e {
			err = append(err, e)
		}
	}
	return
}

func (g *consumerGroup) InitConsumerGroup(settings []ConsumerSetting) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	group := *collectionx.NewOrderHashGroup()
	if len(settings) == 0 {
		return nil
	}
	for idx, _ := range settings {
		consumer, err := settings[idx].Mode.NewMessageConsumer()
		if nil != err {
			return err
		}
		consumer.SetId(settings[idx].Id)
		err = group.Add(consumer)
		if nil != err {
			return err
		}
	}
	g.group = group
	return nil
}

//----------------------

func (g *consumerGroup) ConsumeMessage(msg message.IMessageContext, consumerId string) error {
	if nil == msg {
		return ErrConsumerMessageNil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	ele, ok := g.group.Get(consumerId)
	if !ok {
		return ErrConsumerIdUnknown
	}
	return ele.(IMessageConsumer).ConsumeMessage(msg)
}

func (g *consumerGroup) ConsumeMessages(msg []message.IMessageContext, consumerId string) error {
	if 0 == len(msg) {
		return ErrConsumerMessageNil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	ele, ok := g.group.Get(consumerId)
	if !ok {
		return ErrConsumerIdUnknown
	}
	return ele.(IMessageConsumer).ConsumeMessages(msg)
}
