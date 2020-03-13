package producer

import (
	"github.com/xuzhuoxi/LegoMQ-go/broker"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/netx"
)

func NewRPCMessageProducer() broker.IRPCMessageProducer {
	rpcServer := netx.NewRPCServer()
	return &rpcMessageProducer{rpcServer: rpcServer}
}

type rpcMessageProducer struct {
	rpcServer netx.IRPCServer
	queue     broker.IMessageQueue
}

func (b *rpcMessageProducer) SetMessageQueue(queue broker.IMessageQueue) {
	b.queue = queue
}

func (b *rpcMessageProducer) ProduceMessage(msg message.IMessageContext) error {
	if nil == b.queue {
		return broker.ErrQueueNotPrepared
	}
	b.queue.WriteMessage(msg)
	return nil
}

func (b *rpcMessageProducer) RPCServer() netx.IRPCServer {
	return b.rpcServer
}

func (b *rpcMessageProducer) Register(rcvr interface{}) error {
	return b.rpcServer.Register(rcvr)
}

func (b *rpcMessageProducer) StartRPCListener(addr string) error {
	go b.rpcServer.StartServer(addr)
	return nil
}

func (b *rpcMessageProducer) StopRPCListener() error {
	b.rpcServer.StopServer()
	return nil
}
