package producer

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
	"net/http"
)

const (
	EventMessageOnProducer      string = "producer.EventMessageOnProducer"
	EventMultiMessageOnProducer string = "producer.EventMultiMessageOnProducer"
)

var (
	ErrProducerMessageNil    = errors.New("MessageProducer: Message is nil! ")
	ErrProducerMessagesEmpty = errors.New("MessageProducer: Message array is empty! ")
	ErrProducerModeRegister  = errors.New("MessageProducer: ProducerMode Unregister! ")
)

// 消息生产者
type IMessageProducer interface {
	eventx.IEventDispatcher
	// 生产消息
	// 抛出事件 EventMessageOnProducer
	// err:
	//		msg为nil时ErrProducerMessageNil
	NotifyMessageProduced(msg message.IMessageContext) error
	// 生产消息
	// 抛出事件 EventMultiMessageOnProducer
	// err:
	//		msg长度为0时ErrProducerMessagesEmpty
	NotifyMessagesProduced(msg []message.IMessageContext) error
}

// Socket服务消息生成者
type ISockMessageProducer interface {
	IMessageProducer
	// 初始化服务实例
	InitSockServer(sockNetwork netx.SockNetwork) (s netx.ISockServer, err error)
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

type ProducerMode int

const (
	HttpProducer ProducerMode = iota + 1
	SockProducer
	RPCProducer
	CustomizeProducer
)

var (
	// 函数映射表
	producerMap = make(map[ProducerMode]func() IMessageProducer)
)

func (cm ProducerMode) NewMessageProducer() (p IMessageProducer, err error) {
	if v, ok := producerMap[cm]; ok {
		return v(), nil
	} else {
		return nil, ErrProducerModeRegister
	}
}

// 分组策略注册入口
func RegisterProducerMode(m ProducerMode, f func() IMessageProducer) {
	producerMap[m] = f
}

func init() {
	RegisterProducerMode(HttpProducer, newHttpMessageProducer)
	RegisterProducerMode(SockProducer, newSockMessageProducer)
	RegisterProducerMode(RPCProducer, newRPCMessageProducer)
}
