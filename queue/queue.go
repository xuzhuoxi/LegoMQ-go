package queue

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
)

var (
	ErrQueueSize           = errors.New("ContextQueue: size < 1. ")
	ErrQueueModeUnregister = errors.New("ContextQueue: QueueMode Unregister! ")
)

var (
	ErrQueueClosed = errors.New("ContextQueue: Queue is closed! ")
	ErrQueueFull   = errors.New("ContextQueue: Queue is Full! ")
	ErrQueueEmpty  = errors.New("ContextQueue: Queue is Empty! ")
)

var (
	ErrQueueMessageNil    = errors.New("ContextQueue: Message is nil! ")
	ErrQueueMessagesEmpty = errors.New("ContextQueue: Count is <= 0! ")
)

// 函数映射表
var newQueueFuncArr = make([]func(maxSize int) (c IMessageContextQueue, err error), 16, 16)

// 队列模式(队列类型)
type QueueMode int

// 创建队列实例
// err:
//		ErrQueueModeUnregister:	实例化功能未注册
func (m QueueMode) NewContextQueue(maxSize int) (c IMessageContextQueue, err error) {
	if v := newQueueFuncArr[m]; nil != v {
		return v(maxSize)
	} else {
		return nil, ErrQueueModeUnregister
	}
}

const (
	// 使用channel实现的阻塞队列
	ChannelBlockingQueue QueueMode = iota + 1
	// 使用channel实现的非阻塞队列
	ChannelNBlockingQueue
	// 使用数组实现的队列
	ArrayQueueUnsafe
	// 使用数组实现的队列(并发安全)
	ArrayQueueSafe
	// 自定义队列
	CustomizeQueue
)

// 创建队列实例
// mode:	队列模式
// maxSize:	队列的最大容量
// err:
//		ErrQueueModeUnregister:	实例化功能未注册
func NewContextQueue(mode QueueMode, maxSize int) (q IMessageContextQueue, err error) {
	return mode.NewContextQueue(maxSize)
}

// 队列读取接口
type IMessageContextQueueReader interface {
	// 向缓存区读出一个消息
	// err: 读取异常
	//		ErrQueueClosed:	队列关闭
	//		ErrQueueEmpty:	队列空
	ReadContext() (ctx message.IMessageContext, err error)
	// 向缓存区读出多个消息
	// ctx: 返回实际读出的消息
	// err: 读取异常
	//		ErrQueueMessagesEmpty:	count数量异常
	//		ErrQueueClosed:			队列关闭
	//		ErrQueueEmpty:			队列空
	ReadContexts(count int) (ctx []message.IMessageContext, err error)
	// 向缓存区读出多个消息
	// count: 返回实际读出的消息数量
	// err: 读取异常
	//		ErrQueueMessagesEmpty:	ctx数量异常
	//		ErrQueueClosed:			队列关闭
	//		ErrQueueEmpty			队列空
	ReadContextsTo(ctx []message.IMessageContext) (count int, err error)
}

// 队列写入接口
type IMessageContextQueueWriter interface {
	// 向缓存区写入一个消息
	// err: 写入异常
	//		ctx为nil：ErrQueueMessageNil
	//		队列满：ErrQueueFull
	WriteContext(ctx message.IMessageContext) error
	// 向缓存区写入多个消息
	// count: 返回实际写入数量
	// err: 写入异常
	//		ctx包含nil：ErrQueueMessageNil
	//		ctx数量异常：ErrQueueMessagesEmpty
	//		队列满：ErrQueueFull
	WriteContexts(ctx []message.IMessageContext) (count int, err error)
}

// 队列接口
type IMessageContextQueue interface {
	support.IIdSupport      // 标识支持
	support.ILocateSupport  // 位置支持
	support.IFormatsSupport // 格式匹配支持

	IMessageContextQueueReader // 读操作
	IMessageContextQueueWriter // 写操作

	// Cache最大容量
	MaxSize() int
	// Cache当前容量
	Size() int
	// 关闭
	Close()
}

// 注册
func RegisterQueueMode(m QueueMode, f func(maxSize int) (c IMessageContextQueue, err error)) {
	newQueueFuncArr[m] = f
}

func init() {
	RegisterQueueMode(ChannelBlockingQueue, NewChannelBlockingQueue)
	RegisterQueueMode(ChannelNBlockingQueue, NewChannelNonBlockingQueue)
	RegisterQueueMode(ArrayQueueUnsafe, NewUnsafeArrayQueue)
	RegisterQueueMode(ArrayQueueSafe, NewSafeArrayQueue)
}
