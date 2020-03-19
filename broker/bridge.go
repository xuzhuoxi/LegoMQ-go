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
)

var (
	ErrBridgeProducerNil = errors.New("")
	ErrBridgeQueueNil    = errors.New("")
	ErrBridgeConsumerNil = errors.New("")
	ErrBridgeRoutingNil  = errors.New("")
)

type IBridgeProducer2Queue interface {
	BridgeLink(in producer.IMessageProducerGroup, routing routing.IRoutingStrategy, out queue.IMessageQueueGroup) error
	BridgeUnlink() error
}

type IBridgeQueue2Consumer interface {
	BridgeLink(in queue.IMessageQueueGroup, routing routing.IRoutingStrategy, out consumer.IMessageConsumerGroup) error
	BridgeUnlink() error
}

func newBridgeProducer2Queue() IBridgeProducer2Queue {
	return &pqBridge{}
}

func newBridgeQueue2Consumer() IBridgeQueue2Consumer {
	return &qcBridge{}
}

//---------------

type pqBridge struct {
	pGroup  producer.IMessageProducerGroup
	qGroup  queue.IMessageQueueGroup
	routing routing.IRoutingStrategy
	mu      sync.RWMutex
}

func (b *pqBridge) BridgeLink(in producer.IMessageProducerGroup, routing routing.IRoutingStrategy, out queue.IMessageQueueGroup) error {
	if nil == in {
		return ErrBridgeProducerNil
	}
	if nil == out {
		return ErrBridgeQueueNil
	}
	if nil == routing {
		return ErrBridgeRoutingNil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.pGroup, b.routing, b.qGroup = in, routing, out
	b.pGroup.AddEventListener(producer.EventMessageOnProducer, b.onProduced)
	b.pGroup.AddEventListener(producer.EventMultiMessageOnProducer, b.onMultiProduced)
	return nil
}

func (b *pqBridge) BridgeUnlink() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.pGroup.RemoveEventListener(producer.EventMultiMessageOnProducer, b.onMultiProduced)
	b.pGroup.RemoveEventListener(producer.EventMessageOnProducer, b.onProduced)
	b.pGroup, b.routing, b.qGroup = nil, nil, nil
	return nil
}

func (g *pqBridge) onProduced(evt *eventx.EventData) {
	msg := evt.Data.(message.IMessageContext)
	if nil == msg {
		return
	}
	tIds, err := g.routing.Route(msg.RoutingKey())
	if nil != err {
		return
	}
	g.qGroup.WriteMessageToMulti(msg, tIds)
}

func (g *pqBridge) onMultiProduced(evt *eventx.EventData) {
	msgArr := evt.Data.([]message.IMessageContext)
	if 0 == len(msgArr) {
		return
	}
	for idx, _ := range msgArr {
		tIds, err := g.routing.Route(msgArr[idx].RoutingKey())
		if nil != err {
			return
		}
		g.qGroup.WriteMessageToMulti(msgArr[idx], tIds)
	}
}

//---------------

type qcBridge struct {
	qGroup  queue.IMessageQueueGroup
	cGroup  consumer.IMessageConsumerGroup
	routing routing.IRoutingStrategy
	mu      sync.RWMutex
}

func (b *qcBridge) BridgeLink(in queue.IMessageQueueGroup, routing routing.IRoutingStrategy, out consumer.IMessageConsumerGroup) error {
	if nil == in {
		return ErrBridgeQueueNil
	}
	if nil == out {
		return ErrBridgeConsumerNil
	}
	if nil == routing {
		return ErrBridgeRoutingNil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.qGroup, b.routing, b.cGroup = in, routing, out
	return nil
}

func (b *qcBridge) BridgeUnlink() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.qGroup, b.routing, b.cGroup = nil, nil, nil
	return nil
}

func (g *qcBridge) onProduced(evt *eventx.EventData) {
	msg := evt.Data.(message.IMessageContext)
	if nil == msg {
		return
	}
	tIds, err := g.routing.Route(msg.RoutingKey())
	if nil != err {
		return
	}
	g.qGroup.WriteMessageToMulti(msg, tIds)
}

func (g *qcBridge) onMultiProduced(evt *eventx.EventData) {
	msgArr := evt.Data.([]message.IMessageContext)
	if 0 == len(msgArr) {
		return
	}
	for idx, _ := range msgArr {
		tIds, err := g.routing.Route(msgArr[idx].RoutingKey())
		if nil != err {
			return
		}
		g.qGroup.WriteMessageToMulti(msgArr[idx], tIds)
	}
}
