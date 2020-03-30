package producer

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
	"net/http"
	"time"
)

func NewHttpMessageProducer() IHttpMessageProducer {
	return newHttpMessageProducer().(IHttpMessageProducer)
}

func newHttpMessageProducer() IMessageProducer {
	httpServer := netx.NewHttpServer()
	return &httpMessageProducer{httpServer: httpServer}
}

//------------------

type httpMessageProducer struct {
	eventx.EventDispatcher
	support.ElementSupport

	httpServer netx.IHttpServer
}

func (p *httpMessageProducer) NotifyMessageProduced(msg message.IMessageContext) error {
	if nil == msg {
		return ErrProducerMessageNil
	}
	p.notifyMsgProduced(msg)
	return nil
}

func (p *httpMessageProducer) NotifyMessagesProduced(msg []message.IMessageContext) error {
	if len(msg) == 0 {
		return ErrProducerMessagesEmpty
	}
	p.notifyMultiMsgProduced(msg)
	return nil
}

func (p *httpMessageProducer) HttpServer() netx.IHttpServer {
	return p.httpServer
}

func (p *httpMessageProducer) MapFunc(pattern string, f func(w http.ResponseWriter, r *http.Request)) {
	p.httpServer.MapFunc(pattern, f)
}

func (p *httpMessageProducer) StartHttpListener(addr string) error {
	if p.httpServer.Running() {
		return netx.ErrHttpServerStarted
	}
	var err error
	go func() {
		err = p.httpServer.StartServer(addr)
	}()
	time.Sleep(time.Millisecond * 20)
	return err
}

func (p *httpMessageProducer) StopHttpListener() error {
	return p.httpServer.StopServer()
}

func (p *httpMessageProducer) notifyMsgProduced(msg message.IMessageContext) {
	p.DispatchEvent(EventMessageOnProducer, p, msg)
}

func (p *httpMessageProducer) notifyMultiMsgProduced(msg []message.IMessageContext) {
	p.DispatchEvent(EventMultiMessageOnProducer, p, msg)
}

//------------------

func (p *httpMessageProducer) onRequestTest(w http.ResponseWriter, r *http.Request) error {
	msgBody := r.FormValue("msg")
	msg := message.NewMessageContext("", r.Host, nil, msgBody)
	p.notifyMsgProduced(msg)
	return nil
}
