package producer

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
	"time"
)

func NewSockMessageProducer() ISockMessageProducer {
	return &sockMessageProducer{}
}

func newSockMessageProducer() IMessageProducer {
	return &sockMessageProducer{}
}

//------------------

type sockMessageProducer struct {
	eventx.EventDispatcher
	support.ElementSupport
	ProducerSettingSupport

	sockServer netx.ISockServer
}

func (p *sockMessageProducer) InitProducer() error {
	if "" == p.setting.Id {
		return errors.New("Id is empty. ")
	}
	p.SetId(p.setting.Id)
	p.SetLocateId(p.setting.LocateId)
	return nil
}

func (p *sockMessageProducer) StartProducer() error {
	return p.start(p.setting.Sock)
}

func (p *sockMessageProducer) StopProducer() error {
	return p.stop()
}

func (p *sockMessageProducer) NotifyMessageProduced(msg message.IMessageContext) error {
	if nil == msg {
		return ErrProducerMessagesEmpty
	}
	p.notifyMsgProduced(msg)
	return nil
}

func (p *sockMessageProducer) NotifyMessagesProduced(msg []message.IMessageContext) error {
	if len(msg) == 0 {
		return ErrProducerMessagesEmpty
	}
	p.notifyMultiMsgProduced(msg)
	return nil
}

func (p *sockMessageProducer) InitSockServer(sockNetwork netx.SockNetwork) (s netx.ISockServer, err error) {
	s, err = sockNetwork.NewServer()
	if nil != err {
		return nil, err
	}
	p.sockServer = s
	return
}

func (p *sockMessageProducer) SockServer() netx.ISockServer {
	return p.sockServer
}

func (p *sockMessageProducer) AppendPackHandler(handler netx.FuncPackHandler) error {
	return p.sockServer.GetPackHandlerContainer().AppendPackHandler(p.onSockPackTest)
}

func (p *sockMessageProducer) StartSockListener(params netx.SockParams) error {
	return p.start(params)
}

func (p *sockMessageProducer) StopSockListener() error {
	return p.stop()
}

func (p *sockMessageProducer) start(params netx.SockParams) error {
	if p.sockServer.Running() {
		return netx.ErrSockServerStarted
	}
	var err error
	go func() {
		err = p.sockServer.StartServer(params)
	}()
	time.Sleep(time.Millisecond * 20)
	return err
}

func (p *sockMessageProducer) stop() error {
	p.sockServer.GetPackHandlerContainer().ClearHandler(p.onSockPackTest)
	return p.sockServer.StopServer()
}

func (p *sockMessageProducer) notifyMsgProduced(msg message.IMessageContext) {
	p.DispatchEvent(EventMessageOnProducer, p, msg)
}

func (p *sockMessageProducer) notifyMultiMsgProduced(msg []message.IMessageContext) {
	p.DispatchEvent(EventMultiMessageOnProducer, p, msg)
}

func (p *sockMessageProducer) onSockPackTest(data []byte, senderAddress string, other interface{}) (catch bool) {
	msg := message.NewMessageContext("", senderAddress, nil, data)
	p.notifyMsgProduced(msg)
	return true
}
