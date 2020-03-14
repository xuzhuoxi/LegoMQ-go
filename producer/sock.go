package producer

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
)

func NewSockMessageProducer(sockNetwork netx.SockNetwork) (p ISockMessageProducer, err error) {
	sockServer, err := sockNetwork.NewServer()
	if nil != err {
		return nil, err
	}
	return &sockMessageProducer{sockServer: sockServer}, nil
}

type sockMessageProducer struct {
	eventx.EventDispatcher
	sockServer netx.ISockServer
}

func (p *sockMessageProducer) NotifyMessageProduced(msg message.IMessageContext) error {
	if nil == msg {
		return message.ErrMessageContextNil
	}
	p.notifyMsgProduced(msg)
	return nil
}

func (p *sockMessageProducer) SockServer() netx.ISockServer {
	return p.sockServer
}

func (p *sockMessageProducer) AppendPackHandler(handler netx.FuncPackHandler) error {
	return p.sockServer.GetPackHandlerContainer().AppendPackHandler(p.onSockPackTest)
}

func (p *sockMessageProducer) StartSockListener(params netx.SockParams) error {
	if p.sockServer.Running() {
		return netx.ErrSockServerStarted
	}
	go p.sockServer.StartServer(params)
	return nil
}

func (p *sockMessageProducer) StopSockListener() error {
	p.sockServer.GetPackHandlerContainer().ClearHandler(p.onSockPackTest)
	return p.sockServer.StopServer()
}

func (p *sockMessageProducer) notifyMsgProduced(msg message.IMessageContext) {
	p.DispatchEvent(EventOnProducer, p, msg)
}

func (p *sockMessageProducer) onSockPackTest(data []byte, senderAddress string, other interface{}) (catch bool) {
	msg := message.NewMessageContext(senderAddress, nil, senderAddress, senderAddress, data)
	p.notifyMsgProduced(msg)
	return true
}
