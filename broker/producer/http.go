package producer

import (
	"github.com/xuzhuoxi/LegoMQ-go/broker"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/netx"
	"net/http"
)

func NewHttpMessageProducer() broker.IHttpMessageProducer {
	httpServer := netx.NewHttpServer()
	return &httpMessageProducer{httpServer: httpServer}
}

type httpMessageProducer struct {
	httpServer netx.IHttpServer
	queue      broker.IMessageQueue
}

func (b *httpMessageProducer) SetMessageQueue(queue broker.IMessageQueue) {
	b.queue = queue
}

func (b *httpMessageProducer) ProduceMessage(msg message.IMessageContext) error {
	if nil == b.queue {
		return broker.ErrQueueNotPrepared
	}
	b.queue.WriteMessage(msg)
	return nil
}

func (b *httpMessageProducer) HttpServer() netx.IHttpServer {
	return b.httpServer
}

func (b *httpMessageProducer) MapFunc(pattern string, f func(w http.ResponseWriter, r *http.Request)) {
	b.httpServer.MapFunc(pattern, f)
}

func (b *httpMessageProducer) StartHttpListener(addr string) error {
	if b.httpServer.Running() {
		return netx.ErrHttpServerStarted
	}
	go b.httpServer.StartServer(addr)
	return nil
}

func (b *httpMessageProducer) StopHttpListener() error {
	return b.httpServer.StopServer()
}

func (b *httpMessageProducer) onRequestTest(w http.ResponseWriter, r *http.Request) error {
	host := r.Host
	msgBody := r.FormValue("msg")
	msg := message.NewMessageContext(host, nil, host, host, msgBody)
	b.queue.WriteMessage(msg)
	return nil
}
