package producer

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
)

func NewRPCMessageProducer() IRPCMessageProducer {
	return newRPCMessageProducer().(IRPCMessageProducer)
}

func newRPCMessageProducer() IMessageProducer {
	rpcServer := netx.NewRPCServer()
	return &rpcMessageProducer{rpcServer: rpcServer}
}

//------------------

type rpcMessageProducer struct {
	eventx.EventDispatcher
	id        string
	rpcServer netx.IRPCServer
}

func (p *rpcMessageProducer) Id() string {
	return p.id
}

func (p *rpcMessageProducer) SetId(Id string) {
	p.id = Id
}

func (p *rpcMessageProducer) NotifyMessageProduced(msg message.IMessageContext) error {
	if nil == msg {
		return ErrProducerMessageNil
	}
	p.notifyMsgProduced(msg)
	return nil
}

func (p *rpcMessageProducer) NotifyMessagesProduced(msg []message.IMessageContext) error {
	if len(msg) == 0 {
		return ErrProducerMessagesEmpty
	}
	p.notifyMultiMsgProduced(msg)
	return nil
}

func (p *rpcMessageProducer) RPCServer() netx.IRPCServer {
	return p.rpcServer
}

func (p *rpcMessageProducer) Register(rcvr interface{}) error {
	return p.rpcServer.Register(rcvr)
}

func (p *rpcMessageProducer) StartRPCListener(addr string) error {
	go p.rpcServer.StartServer(addr)
	return nil
}

func (p *rpcMessageProducer) StopRPCListener() error {
	p.rpcServer.StopServer()
	return nil
}

func (p *rpcMessageProducer) notifyMsgProduced(msg message.IMessageContext) {
	p.DispatchEvent(EventMessageOnProducer, p, msg)
}

func (p *rpcMessageProducer) notifyMultiMsgProduced(msg []message.IMessageContext) {
	p.DispatchEvent(EventMultiMessageOnProducer, p, msg)
}
