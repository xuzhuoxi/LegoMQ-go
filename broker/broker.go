package broker

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/queue"
	"github.com/xuzhuoxi/LegoMQ-go/routing"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"net/http"
)

var (
	ErrQueueNotPrepared = errors.New("Queue is not prepared. ")
	ErrQueueIndex       = errors.New("Index out of queue size. ")

	ErrConsumerStarted        = errors.New("Consumer is started. ")
	ErrConsumerStopping       = errors.New("Consumer is stopping. ")
	ErrConsumerStopped        = errors.New("Consumer is stopped. ")
	ErrConsumerQueueNil       = errors.New("ContextQueue is nil. ")
	ErrConsumerTypeUnregister = errors.New("ConsumerMode does not registered. ")
)

//-----------------

// 唯一标识支持
type IMessageIdSupport interface {
	Id() string
	SetId(Id string)
}

//-----------------

// 消息生产者
type IMessageProducer interface {
	// 设置消息队列
	SetMessageQueue(queue IMessageQueue)
	// 生产消息(把消息加入队列)
	ProduceMessage(msg message.IMessageContext) error
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

//-----------------

// 消息队列组配置接口
type IMessageQueueConfig interface {
	// 设置消息处理模式
	// beforeMode: 前置分组模式
	// queueMode: 队列模式
	// afterMode: 队列模式
	SetMode(beforeMode routing.RoutingMode, queueMode queue.QueueMode, afterMode routing.RoutingMode)
	SetMode2(groupMode, mode routing.RoutingMode, queueMode queue.QueueMode)
	// 设置消息处理缓存大小
	// queueSize: 缓存队列数量
	// contextSize: 每个队列的容量
	// afterSize: 每个afterMode
	SetSize(queueSize int, contextSize int, afterSize int)
	SetSize2(groupSize int, queueSize int)
	// 初始化
	InitQueue() error

	// 模式设置参数
	Mode() (before routing.RoutingMode, queue queue.QueueMode, after routing.RoutingMode)
	// 数量设置参数
	Size() (queueSize, contextSize, afterSize int)
}

type IMessageQueueReader interface {
	ReadMessageFrom(index int) (msg message.IMessageContext, err error)
	ReadMessagesFrom(index int, count int) (msg []message.IMessageContext, err error)
	ReadMessagesTo(index int, msg []message.IMessageContext) (count int, err error)
}

type IMessageQueueWriter interface {
	// 往消息队列写入一个消息
	// 返回错误信息，写入成功返回nil
	WriteMessage(msg message.IMessageContext) error
	// 往消息队列写入多个消息
	// 返回说明：
	// count: 成功写入消息的数量
	// err:	错误信息，没有错误则返回nil
	WriteMessages(msg []message.IMessageContext) (count int, err error)
	// 往指定消息队列写入一个消息
	// 返回错误信息，写入成功返回nil
	WriteMessageTo(msg message.IMessageContext, index int) error
	// 往指定消息队列写入多个消息
	// 返回说明：
	// count: 成功写入消息的数量
	// err:	错误信息，没有错误则返回nil
	WriteMessagesTo(msg []message.IMessageContext, index int) (count int, err error)
}

type IMessageQueue interface {
	IMessageIdSupport
	IMessageQueueReader
	IMessageQueueWriter
}

//-----------------

type IMessageConsumer interface {
	BindContextQueue(contextQueue queue.IMessageContextQueue) error
	UnbindContextQueue() error
	// 被动接收一个消息并进行处理
	AcceptMessage(msg message.IMessageContext) error
	StartConsume() error
	StopConsume() error
}

type ILogMessageConsumer interface {
	IMessageConsumer

	logx.ILoggerGetter
	SetConsumerLevel(level logx.LogLevel)
}

//---------------------------

type IMessageDriver interface {
	eventx.IEventDispatcher
	AutoOn() error
	AutoOff() error
}

//---------------------------

type IMessageBroker interface {
	UpdateProducer(producer IMessageProducer)
	Queue() (config IMessageQueueConfig, queue IMessageQueue)
	UpdateConsumer(queueIndex int, consumer IMessageConsumer)
}

type ISockMessageBroker interface {
	//IMessageBroker
	SockServer() netx.ISockServer
	StartSockListener(params netx.SockParams) error
	StopSockListener() error
}

type IHttpMessageBroker interface {
	//IMessageBroker
	HttpServer() netx.IHttpServer
	StartHttpListener(addr string) error
	StopHttpListener() error
}
