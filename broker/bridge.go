package broker

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/consumer"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/producer"
	"github.com/xuzhuoxi/LegoMQ-go/queue"
	"github.com/xuzhuoxi/LegoMQ-go/routing"
	"github.com/xuzhuoxi/infra-go/eventx"
	"sync"
	"time"
)

var (
	ErrBridgeProducerNil = errors.New("BrokerBridge: Producer is nil. ")
	ErrBridgeQueueNil    = errors.New("BrokerBridge: Queue is nil. ")
	ErrBridgeConsumerNil = errors.New("BrokerBridge: Consumer is nil. ")
	ErrBridgeRoutingNil  = errors.New("BrokerBridge: Routing strategy is nil. ")

	ErrBridgeLinked    = errors.New("BrokerBridge: Bridge is linked. ")
	ErrBridgeNotLinked = errors.New("BrokerBridge: Bridge is not linked. ")

	ErrBridgeDriverDuration = errors.New("BrokerBridge: Duration < 0. ")
	ErrBridgeDriverQuantity = errors.New("BrokerBridge: Quantity < 1. ")
)

type iBridge interface {
	// 设置路由模式
	// err:
	//		ErrBridgeLinked: 		已经连接
	// 		ErrRougingRegister:		路由模式未注册
	SetRoutingMode(mode routing.RoutingMode) error
	// 设置路由策略
	// err:
	//		ErrBridgeLinked: 		已经连接
	// 		ErrBridgeRoutingNil:	strategy=nil
	SetRoutingStrategy(strategy routing.IRoutingStrategy) error
	// 连接
	// err:
	//		ErrBridgeLinked: 		已经连接
	Link() error
	// 取消连接
	// err:
	//		ErrBridgeNotLinked: 	未连接
	Unlink() error
}

// Producer到Queue桥接器
// 提供功能：
// 		1.Producer消息捕捉
//		2.消息路由
// 		3.消息加入Queue
type IBridgeProducer2Queue interface {
	iBridge
	// 设置桥接点
	// err:
	//		ErrBridgeLinked: 		已经连接
	//		ErrBridgeProducerNil: 	pierIn=nil
	//		ErrBridgeQueueNil: 		pierOut=nil
	SetBridgePier(pierIn producer.IMessageProducerGroup, pierOut queue.IMessageQueueGroup) error
}

// Queue到Consumer桥接器
// 提供功能：
// 		1.按时间片从Queue中提取消息
//		2.消息路由
// 		3.消息加入Consumer处理
type IBridgeQueue2Consumer interface {
	iBridge
	// 设置桥接点
	// err:
	//		ErrBridgeLinked: 		已经连接
	//		ErrBridgeQueueNil: 		pierIn=nil
	//		ErrBridgeConsumerNil: 	pierOut=nil
	SetBridgePier(pierIn queue.IMessageQueueGroup, pierOut consumer.IMessageConsumerGroup) error
	// 初始化驱动器
	// duration: 频率
	// quantity: 批量
	// err:
	// 		ErrBridgeDriverDuration: duration<0
	//		ErrBridgeDriverQuantity: quantity<1
	InitDriver(duration time.Duration, quantity int) error
}

func NewBridgeProducer2Queue() IBridgeProducer2Queue {
	return newBridgeProducer2Queue()
}

func NewBridgeQueue2Consumer() IBridgeQueue2Consumer {
	return newBridgeQueue2Consumer()
}

func newBridgeProducer2Queue() IBridgeProducer2Queue {
	return &p2qBridge{}
}

func newBridgeQueue2Consumer() IBridgeQueue2Consumer {
	return &q2cBridge{}
}

//---------------

type p2qBridge struct {
	pGroup  producer.IMessageProducerGroup
	qGroup  queue.IMessageQueueGroup
	routing routing.IRoutingStrategy

	mu      sync.RWMutex
	started bool
}

func (b *p2qBridge) SetBridgePier(pierIn producer.IMessageProducerGroup, pierOut queue.IMessageQueueGroup) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeLinked
	}
	if nil == pierIn {
		return ErrBridgeProducerNil
	}
	if nil == pierOut {
		return ErrBridgeQueueNil
	}
	b.pGroup, b.qGroup = pierIn, pierOut
	return nil
}

func (b *p2qBridge) SetRoutingMode(mode routing.RoutingMode) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeLinked
	}
	s, err := mode.NewRoutingStrategy()
	if nil != err {
		return err
	}
	b.routing = s
	return nil
}

func (b *p2qBridge) SetRoutingStrategy(strategy routing.IRoutingStrategy) error {
	if nil == strategy {
		return ErrBridgeRoutingNil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeLinked
	}
	b.routing = strategy
	return nil
}

func (b *p2qBridge) Link() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeLinked
	}
	b.routing.Config().SetRoutingTargets(b.qGroup.Config().RoutingElements())
	b.pGroup.Config().SetProducedFunc(b.onProduced, b.onMultiProduced)
	b.started = true
	return nil
}

func (b *p2qBridge) Unlink() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.started {
		return ErrBridgeNotLinked
	}
	b.pGroup, b.routing, b.qGroup = nil, nil, nil
	b.started = false
	return nil
}

func (b *p2qBridge) onProduced(msg message.IMessageContext, locateId string) {
	if nil == msg {
		return
	}
	tIds, err := b.routing.Route(msg.RoutingKey(), locateId)
	if nil != err {
		return
	}
	b.qGroup.WriteMessageToMulti(msg, tIds)
	//fmt.Println("onProduced:", msg, locateId, tIds)
	//_, _, err1, err2 := b.qGroup.WriteMessageToMulti(msg, tIds)
	//if err1 != nil {
	//	fmt.Println("err1:", err1)
	//}
	//if err2 != nil {
	//	fmt.Println("err2:", err2)
	//}
}

func (b *p2qBridge) onMultiProduced(msgArr []message.IMessageContext, locateId string) {
	if 0 == len(msgArr) {
		return
	}
	for idx, _ := range msgArr {
		tIds, err := b.routing.Route(msgArr[idx].RoutingKey(), locateId)
		if nil != err {
			return
		}
		b.qGroup.WriteMessageToMulti(msgArr[idx], tIds)
	}
}

//---------------

type q2cBridge struct {
	qGroup  queue.IMessageQueueGroup
	cGroup  consumer.IMessageConsumerGroup
	routing routing.IRoutingStrategy

	driver   ITimeSliceDriver
	msgCache [][]message.IMessageContext
	duration time.Duration
	quantity int

	mu      sync.RWMutex
	started bool
}

func (b *q2cBridge) SetBridgePier(pierIn queue.IMessageQueueGroup, pierOut consumer.IMessageConsumerGroup) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeLinked
	}
	if nil == pierIn {
		return ErrBridgeQueueNil
	}
	if nil == pierOut {
		return ErrBridgeConsumerNil
	}
	b.qGroup, b.cGroup = pierIn, pierOut
	return nil
}

func (b *q2cBridge) SetRoutingMode(mode routing.RoutingMode) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeLinked
	}
	s, err := mode.NewRoutingStrategy()
	if nil != err {
		return err
	}
	b.routing = s
	return nil
}

func (b *q2cBridge) SetRoutingStrategy(strategy routing.IRoutingStrategy) error {
	if nil == strategy {
		return ErrBridgeRoutingNil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeLinked
	}
	b.routing = strategy
	return nil
}

func (b *q2cBridge) InitDriver(duration time.Duration, quantity int) error {
	if duration < 0 {
		return ErrBridgeDriverDuration
	}
	if quantity < 1 {
		return ErrBridgeDriverQuantity
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if quantity > 1 {
		b.msgCache = make([][]message.IMessageContext, b.qGroup.Config().QueueSize(), b.qGroup.Config().QueueSize())
		for idx, _ := range b.msgCache {
			b.msgCache[idx] = make([]message.IMessageContext, quantity, quantity)
		}
	}
	b.duration, b.quantity = duration, quantity
	b.driver = NewTimeSliceDriver(duration)
	return nil

}

func (b *q2cBridge) Link() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeLinked
	}
	b.routing.Config().SetRoutingTargets(b.cGroup.Config().RoutingElements())
	b.driver.AddEventListener(EventOnTime, b.onTime)
	err := b.driver.DriverStart()
	if nil != err {
		return err
	}
	b.started = true
	return nil
}

func (b *q2cBridge) Unlink() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.started {
		return ErrBridgeNotLinked
	}
	err := b.driver.DriverStop()
	if nil != err {
		return err
	}
	b.driver.RemoveEventListener(EventOnTime, b.onTime)
	b.started = false
	return nil
}

func (b *q2cBridge) onTime(evt *eventx.EventData) {
	if b.quantity == 1 {
		b.readEachMessage()
	} else {
		b.readEachMessages()
	}
}

func (b *q2cBridge) readEachMessage() {
	b.qGroup.ForEachElement(func(index int, ele queue.IMessageContextQueue) (stop bool) {
		ctx, err := ele.ReadContext()
		if nil == err {
			b.handleMessage(ctx, ele.LocateId())
		}
		return false
	})
}

func (b *q2cBridge) readEachMessages() {
	b.qGroup.ForEachElement(func(index int, ele queue.IMessageContextQueue) (stop bool) {
		count, _ := ele.ReadContextsTo(b.msgCache[index])
		if count > 0 {
			b.handleMessages(b.msgCache[index][:count], ele.LocateId())
		}
		return false
	})
}

func (b *q2cBridge) handleMessage(msg message.IMessageContext, locateId string) {
	ids, err := b.routing.Route(msg.RoutingKey(), locateId)
	if nil != err {
		return
	}
	b.cGroup.ConsumeMessageMulti(msg, ids)
}

func (b *q2cBridge) handleMessages(msgs []message.IMessageContext, locateId string) {
	for index, _ := range msgs {
		ids, err := b.routing.Route(msgs[index].RoutingKey(), locateId)
		if nil != err {
			continue
		}
		b.cGroup.ConsumeMessageMulti(msgs[index], ids)
	}
}
