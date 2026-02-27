package producer

import (
	"net/http"

	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/netx/httpx"
)

func NewHttpMessageProducer() IHttpMessageProducer {
	return &httpMessageProducer{}
}

func newHttpMessageProducer() IMessageProducer {
	return &httpMessageProducer{}
}

//------------------

type httpMessageProducer struct {
	eventx.EventDispatcher
	support.ElementSupport
	ProducerSettingSupport

	httpServer httpx.IHttpServer
}

func (p *httpMessageProducer) InitProducer() error {
	if "" == p.setting.Id {
		return support.ErrIdEmpty
	}
	p.SetId(p.setting.Id)
	p.SetLocateId(p.setting.LocateId)
	p.httpServer = httpx.NewHttpServer()
	return nil
}

func (p *httpMessageProducer) StartProducer() error {
	return p.start(p.setting.Http.Addr)
}

func (p *httpMessageProducer) StopProducer() error {
	return p.stop()
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

func (p *httpMessageProducer) InitHttpServer() (s httpx.IHttpServer, err error) {
	p.httpServer = httpx.NewHttpServer()
	return p.httpServer, nil
}

func (p *httpMessageProducer) HttpServer() httpx.IHttpServer {
	return p.httpServer
}

func (p *httpMessageProducer) MapFunc(pattern string, f func(w http.ResponseWriter, r *http.Request)) {
	p.httpServer.MapFunc(pattern, f)
}

func (p *httpMessageProducer) StartHttpListener(addr string) error {
	return p.start(addr)
}

func (p *httpMessageProducer) StopHttpListener() error {
	return p.stop()
}

//--------------------

func (p *httpMessageProducer) start(addr string) error {
	if p.httpServer.Running() {
		return netx.ErrHttpServerStarted
	}
	return p.httpServer.StartServer(addr)
}

func (p *httpMessageProducer) stop() error {
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
