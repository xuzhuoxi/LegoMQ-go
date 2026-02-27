package producer

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx/rpcx"
)

func NewRPCMessageProducer() IRPCMessageProducer {
	return &rpcMessageProducer{}
}

func newRPCMessageProducer() IMessageProducer {
	return &rpcMessageProducer{}
}

//------------------

type rpcMessageProducer struct {
	eventx.EventDispatcher
	support.ElementSupport
	ProducerSettingSupport

	rpcServer rpcx.IRPCServer
}

func (p *rpcMessageProducer) InitProducer() error {
	if "" == p.setting.Id {
		return support.ErrIdEmpty
	}
	p.SetId(p.setting.Id)
	p.SetLocateId(p.setting.LocateId)
	p.rpcServer = rpcx.NewRPCServer()
	return nil
}

func (p *rpcMessageProducer) StartProducer() error {
	return p.start(p.setting.RPC.Addr)
}

func (p *rpcMessageProducer) StopProducer() error {
	return p.stop()
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

func (p *rpcMessageProducer) InitRPCServer() (s rpcx.IRPCServer, err error) {
	p.rpcServer = rpcx.NewRPCServer()
	return p.rpcServer, nil
}

func (p *rpcMessageProducer) RPCServer() rpcx.IRPCServer {
	return p.rpcServer
}

func (p *rpcMessageProducer) Register(rcvr interface{}) error {
	return p.rpcServer.Register(rcvr)
}

func (p *rpcMessageProducer) StartRPCListener(addr string) error {
	return p.start(addr)
}

func (p *rpcMessageProducer) StopRPCListener() error {
	return p.stop()
}

func (p *rpcMessageProducer) start(addr string) error {
	p.rpcServer.StartServer(addr)
	return nil
}

func (p *rpcMessageProducer) stop() error {
	p.rpcServer.StopServer()
	return nil
}

func (p *rpcMessageProducer) notifyMsgProduced(msg message.IMessageContext) {
	p.DispatchEvent(EventMessageOnProducer, p, msg)
}

func (p *rpcMessageProducer) notifyMultiMsgProduced(msg []message.IMessageContext) {
	p.DispatchEvent(EventMultiMessageOnProducer, p, msg)
}
