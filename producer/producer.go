package producer

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
	"net/http"
)

const (
	EventOnProducer string = "producer.EventOnProducer"
)

// 消息生产者
type IMessageProducer interface {
	eventx.IEventDispatcher
	// 生产消息
	// 抛出事件EventOnProducer
	// err:
	//		msg为nil时ErrMessageContextNil
	NotifyMessageProduced(msg message.IMessageContext) error
}

// Socket服务消息生成者
type ISockMessageProducer interface {
	IMessageProducer
	// Socket服务器
	SockServer() netx.ISockServer
	// 追加Socket服务器接收的信息的处理函数
	AppendPackHandler(handler netx.FuncPackHandler) error
	// 启动Socket服务器
	StartSockListener(params netx.SockParams) error
	// 停止Socket服务器
	StopSockListener() error
}

// Http服务消息生成者
type IHttpMessageProducer interface {
	IMessageProducer
	// Http服务器
	HttpServer() netx.IHttpServer
	// 映射Http请求对应的处理函数
	MapFunc(pattern string, f func(w http.ResponseWriter, r *http.Request))
	// 启动Http服务器
	StartHttpListener(addr string) error
	// 停止Http服务器
	StopHttpListener() error
}

// RPC服务消息生成者
type IRPCMessageProducer interface {
	IMessageProducer
	// RPC服务器
	RPCServer() netx.IRPCServer
	// 注册RPC响应对象
	Register(rcvr interface{}) error
	// 启动RPC服务器
	StartRPCListener(addr string) error
	// 停止RPC服务器
	StopRPCListener() error
}
