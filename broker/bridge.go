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
	ErrBridgeProducerNil = errors.New("")
	ErrBridgeQueueNil    = errors.New("")
	ErrBridgeConsumerNil = errors.New("")
	ErrBridgeRoutingNil  = errors.New("")

	ErrBridgeStarted = errors.New("")
	ErrBridgeStopped = errors.New("")
)

type IBridgeProducer2Queue interface {
	SetEntity(in producer.IMessageProducerGroup, out queue.IMessageQueueGroup) error
	SetRoutingMode(mode routing.RoutingMode) error
	SetRoutingStrategy(strategy routing.IRoutingStrategy) error
	Link() error
	Unlink() error
}

type IBridgeQueue2Consumer interface {
	SetEntity(in queue.IMessageQueueGroup, out consumer.IMessageConsumerGroup) error
	SetRoutingMode(mode routing.RoutingMode) error
	SetRoutingStrategy(strategy routing.IRoutingStrategy) error
	Link(duration time.Duration, maxMessage int) error
	Unlink() error
}

func NewBridgeProducer2Queue() IBridgeProducer2Queue {
	return newBridgeProducer2Queue()
}

func NewBridgeQueue2Consumer() IBridgeQueue2Consumer {
	return newBridgeQueue2Consumer()
}

//---------------

func newBridgeProducer2Queue() IBridgeProducer2Queue {
	return &p2qBridge{}
}

func newBridgeQueue2Consumer() IBridgeQueue2Consumer {
	return &q2cBridge{}
}

type p2qBridge struct {
	pGroup  producer.IMessageProducerGroup
	qGroup  queue.IMessageQueueGroup
	routing routing.IRoutingStrategy
	mu      sync.RWMutex
	started bool
}

func (b *p2qBridge) SetEntity(in producer.IMessageProducerGroup, out queue.IMessageQueueGroup) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeStarted
	}
	if nil == in {
		return ErrBridgeProducerNil
	}
	if nil == out {
		return ErrBridgeQueueNil
	}
	b.pGroup, b.qGroup = in, out
	return nil
}

func (b *p2qBridge) SetRoutingMode(mode routing.RoutingMode) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeStarted
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
		return ErrBridgeStarted
	}
	b.routing = strategy
	return nil
}

func (b *p2qBridge) Link() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeStarted
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
		return ErrBridgeStopped
	}
	b.pGroup, b.routing, b.qGroup = nil, nil, nil
	b.started = false
	return nil
}

func (b *p2qBridge) onProduced(msg message.IMessageContext, locateKey string) {
	if nil == msg {
		return
	}
	tIds, err := b.routing.Route(msg.RoutingKey(), locateKey)
	if nil != err {
		return
	}
	b.qGroup.WriteMessageToMulti(msg, tIds)
	//fmt.Println("onProduced:", msg, locateKey, tIds)
	//_, _, err1, err2 := b.qGroup.WriteMessageToMulti(msg, tIds)
	//if err1 != nil {
	//	fmt.Println("err1:", err1)
	//}
	//if err2 != nil {
	//	fmt.Println("err2:", err2)
	//}
}

func (b *p2qBridge) onMultiProduced(msgArr []message.IMessageContext, locateKey string) {
	if 0 == len(msgArr) {
		return
	}
	for idx, _ := range msgArr {
		tIds, err := b.routing.Route(msgArr[idx].RoutingKey(), locateKey)
		if nil != err {
			return
		}
		b.qGroup.WriteMessageToMulti(msgArr[idx], tIds)
	}
}

//---------------

type q2cBridge struct {
	qGroup   queue.IMessageQueueGroup
	cGroup   consumer.IMessageConsumerGroup
	driver   ITimeSliceDriver
	routing  routing.IRoutingStrategy
	mu       sync.RWMutex
	started  bool
	msgCache [][]message.IMessageContext
}

func (b *q2cBridge) SetEntity(in queue.IMessageQueueGroup, out consumer.IMessageConsumerGroup) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeStarted
	}
	if nil == in {
		return ErrBridgeQueueNil
	}
	if nil == out {
		return ErrBridgeConsumerNil
	}
	b.qGroup, b.cGroup = in, out
	return nil
}

func (b *q2cBridge) SetRoutingMode(mode routing.RoutingMode) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeStarted
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
		return ErrBridgeStarted
	}
	b.routing = strategy
	return nil
}

func (b *q2cBridge) Link(duration time.Duration, maxMessage int) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBridgeStarted
	}
	b.initCache(maxMessage)
	b.routing.Config().SetRoutingTargets(b.cGroup.Config().RoutingElements())
	b.driver = NewTimeSliceDriver(duration)
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
		return ErrBridgeStopped
	}
	err := b.driver.DriverStop()
	if nil != err {
		return err
	}
	b.driver.RemoveEventListener(EventOnTime, b.onTime)
	b.started = false
	return nil
}

func (b *q2cBridge) initCache(maxMessage int) {
	b.msgCache = make([][]message.IMessageContext, b.qGroup.Config().QueueSize(), b.qGroup.Config().QueueSize())
	for idx, _ := range b.msgCache {
		b.msgCache[idx] = make([]message.IMessageContext, maxMessage, maxMessage)
	}
}

func (b *q2cBridge) onTime(evt *eventx.EventData) {
	//fmt.Println(22222222)
	b.qGroup.ForEachElement(func(index int, ele queue.IMessageContextQueue) (stop bool) {
		//ctx, err := ele.ReadContext()
		//if nil == err {
		//	b.handleMessage(ctx, ele.LocateKey())
		//}
		//return false
		count, _ := ele.ReadContextsTo(b.msgCache[index])
		if count > 0 {
			b.handleMessages(b.msgCache[index][:count], ele.LocateKey())
		}
		return false
	})
}

func (b *q2cBridge) handleMessage(msg message.IMessageContext, locateKey string) {
	ids, err := b.routing.Route(msg.RoutingKey(), locateKey)
	if nil != err {
		return
	}
	b.cGroup.ConsumeMessageMulti(msg, ids)
}

func (b *q2cBridge) handleMessages(msgs []message.IMessageContext, locateKey string) {
	for index, _ := range msgs {
		ids, err := b.routing.Route(msgs[index].RoutingKey(), locateKey)
		if nil != err {
			continue
		}
		b.cGroup.ConsumeMessageMulti(msgs[index], ids)
	}
}
