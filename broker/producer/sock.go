package producer

import (
	"github.com/xuzhuoxi/LegoMQ-go/broker"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/netx"
)

func NewSockMessageProducer(sockNetwork netx.SockNetwork) (b broker.ISockMessageProducer, err error) {
	sockServer, err := sockNetwork.NewServer()
	if nil != err {
		return nil, err
	}
	return &sockMessageProducer{sockServer: sockServer}, nil
}

type sockMessageProducer struct {
	sockServer netx.ISockServer
	queue      broker.IMessageQueue
}

func (b *sockMessageProducer) SetMessageQueue(queue broker.IMessageQueue) {
	b.queue = queue
}

func (b *sockMessageProducer) ProduceMessage(msg message.IMessageContext) error {
	if nil == b.queue {
		return broker.ErrQueueNotPrepared
	}
	b.queue.WriteMessage(msg)
	return nil
}

func (b *sockMessageProducer) SockServer() netx.ISockServer {
	return b.sockServer
}

func (b *sockMessageProducer) AppendPackHandler(handler netx.FuncPackHandler) error {
	return b.sockServer.GetPackHandlerContainer().AppendPackHandler(b.onSockPackTest)
}

func (b *sockMessageProducer) StartSockListener(params netx.SockParams) error {
	if b.sockServer.Running() {
		return netx.ErrSockServerStarted
	}
	go b.sockServer.StartServer(params)
	return nil
}

func (b *sockMessageProducer) StopSockListener() error {
	b.sockServer.GetPackHandlerContainer().ClearHandler(b.onSockPackTest)
	return b.sockServer.StopServer()
}

func (b *sockMessageProducer) onSockPackTest(data []byte, senderAddress string, other interface{}) (catch bool) {
	msg := message.NewMessageContext(senderAddress, nil, senderAddress, senderAddress, data)
	b.queue.WriteMessage(msg)
	return true
}
