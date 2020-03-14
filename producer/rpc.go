package producer

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
)

func NewRPCMessageProducer() IRPCMessageProducer {
	rpcServer := netx.NewRPCServer()
	return &rpcMessageProducer{rpcServer: rpcServer}
}

type rpcMessageProducer struct {
	eventx.EventDispatcher
	rpcServer netx.IRPCServer
}

func (b *rpcMessageProducer) NotifyMessageProduced(msg message.IMessageContext) error {
	if nil == msg {
		return message.ErrMessageContextNil
	}
	b.notifyMsgProduced(msg)
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

func (b *rpcMessageProducer) notifyMsgProduced(msg message.IMessageContext) {
	b.DispatchEvent(EventOnProducer, b, msg)
}
