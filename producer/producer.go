package producer

import (
	"errors"
	"net/http"

	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/netx/httpx"
	"github.com/xuzhuoxi/infra-go/netx/rpcx"
)

var (
	ErrProducerMessageNil     = errors.New("MessageProducer: Message is nil! ")
	ErrProducerMessagesEmpty  = errors.New("MessageProducer: Messages is empty! ")
	ErrProducerModeUnregister = errors.New("MessageProducer: ProducerMode Unregister! ")
)

// 函数映射表
var newProducerFuncArr = make([]func() IMessageProducer, 16, 16)

// ProducerMode
// 消息生产者模式
type ProducerMode int

// NewMessageProducer
// 创建消息生产者实例
// err:
//
//	ErrProducerModeUnregister:	实例化功能未注册
func (m ProducerMode) NewMessageProducer() (p IMessageProducer, err error) {
	if v := newProducerFuncArr[m]; nil != v {
		return v(), nil
	} else {
		return nil, ErrProducerModeUnregister
	}
}

const (
	HttpProducer ProducerMode = iota + 1
	SockProducer
	RPCProducer
	CustomizeProducer
)

// NewMessageProducer
// 根据创建消息生产者实例
// mode:	生产者模式
func NewMessageProducer(mode ProducerMode) (c IMessageProducer, err error) {
	return mode.NewMessageProducer()
}

// IMessageProducer
// 消息生产者接口
type IMessageProducer interface {
	eventx.IEventDispatcher // 事件支持
	support.IIdSupport      // 标识支持
	support.ILocateSupport  // 位置支持

	// NotifyMessageProduced
	// 生产消息
	// 抛出事件 EventMessageOnProducer
	// err:
	//		msg为nil时ErrProducerMessageNil
	NotifyMessageProduced(msg message.IMessageContext) error
	// NotifyMessagesProduced
	// 生产消息
	// 抛出事件 EventMultiMessageOnProducer
	// err:
	//		msg长度为0时ErrProducerMessagesEmpty
	NotifyMessagesProduced(msg []message.IMessageContext) error
}

// ISockMessageProducer
// Socket服务消息生成者
type ISockMessageProducer interface {
	IMessageProducer        // 消息生产者
	IProducerSettingSupport //消息生产者设置支持

	// InitSockServer
	// 初始化服务实例
	InitSockServer(sockNetwork netx.SockNetwork) (s netx.ISockServer, err error)
	// SockServer
	// Socket服务器
	SockServer() netx.ISockServer

	// AppendPackHandler
	// 追加Socket服务器接收的信息的处理函数
	AppendPackHandler(handler netx.FuncPackHandler) error
	// StartSockListener
	// 启动Socket服务器
	StartSockListener(params netx.SockParams) error
	// StopSockListener
	// 停止Socket服务器
	StopSockListener() error
}

// IHttpMessageProducer
// Http服务消息生成者
type IHttpMessageProducer interface {
	IMessageProducer        // 消息生产者
	IProducerSettingSupport // 消息生产者设置支持

	// InitHttpServer
	// 初始化服务实例
	InitHttpServer() (s httpx.IHttpServer, err error)
	// HttpServer
	// Http服务器
	HttpServer() httpx.IHttpServer

	// MapFunc
	// 映射Http请求对应的处理函数
	MapFunc(pattern string, f func(w http.ResponseWriter, r *http.Request))
	// StartHttpListener
	// 启动Http服务器
	// 阻塞协程
	StartHttpListener(addr string) error
	// StopHttpListener
	// 停止Http服务器
	StopHttpListener() error
}

// IRPCMessageProducer
// RPC服务消息生成者
type IRPCMessageProducer interface {
	IMessageProducer        // 消息生产者
	IProducerSettingSupport // 消息生产者设置支持

	// InitRPCServer
	// 初始化服务实例
	InitRPCServer() (s rpcx.IRPCServer, err error)
	// RPCServer
	// RPC服务器
	RPCServer() rpcx.IRPCServer

	// Register
	// 注册RPC响应对象
	Register(rcvr interface{}) error
	// StartRPCListener
	// 启动RPC服务器
	StartRPCListener(addr string) error
	// StopRPCListener
	// 停止RPC服务器
	StopRPCListener() error
}

// RegisterProducerMode
// 注册
func RegisterProducerMode(m ProducerMode, f func() IMessageProducer) {
	newProducerFuncArr[m] = f
}

func init() {
	RegisterProducerMode(HttpProducer, newHttpMessageProducer)
	RegisterProducerMode(SockProducer, newSockMessageProducer)
	RegisterProducerMode(RPCProducer, newRPCMessageProducer)
}
