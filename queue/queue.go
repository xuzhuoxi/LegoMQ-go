package queue

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
)

var (
	ErrQueueClosed = errors.New("ContextQueue is closed! ")
	ErrQueueFull   = errors.New("ContextQueue is Full! ")
	ErrQueueEmpty  = errors.New("ContextQueue is Empty! ")
)

var (
	ErrQueueMessageNil     = errors.New("ContextQueue: Message is nil! ")
	ErrQueueCountZero      = errors.New("ContextQueue: Count is <= 0! ")
	ErrSize                = errors.New("ContextQueue: Max size should be larger than 0. ")
	ErrQueueModeUnregister = errors.New("QueueMode Unregister! ")
)

type IMessageContextQueueReader interface {
	// 向缓存区读出一个消息
	// err: 读取异常
	//		队列关闭：ErrQueueClosed
	//		队列空：ErrQueueEmpty
	ReadContext() (ctx message.IMessageContext, err error)
	// 向缓存区读出多个消息
	// ctx: 返回实际读出的消息
	// err: 读取异常
	//		count数量异常：ErrQueueCountZero
	//		队列关闭：ErrQueueClosed
	//		队列空：ErrQueueEmpty
	ReadContexts(count int) (ctx []message.IMessageContext, err error)
	// 向缓存区读出多个消息
	// count: 返回实际读出的消息数量
	// err: 读取异常
	//		ctx数量异常：ErrQueueCountZero
	//		队列关闭：ErrQueueClosed
	//		队列空：ErrQueueEmpty
	ReadContextsTo(ctx []message.IMessageContext) (count int, err error)
}

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
	//		ctx数量异常：ErrQueueCountZero
	//		队列满：ErrQueueFull
	WriteContexts(ctx []message.IMessageContext) (count int, err error)
}

type IMessageContextQueue interface {
	support.IQueueBase
	// Cache最大容量
	MaxSize() int
	// Cache当前容量
	Size() int
	// 关闭
	Close()
	// 读操作
	IMessageContextQueueReader
	// 写操作
	IMessageContextQueueWriter
}

type QueueMode int

const (
	// 使用channel实现的队列
	ChannelQueue QueueMode = iota + 1
	// 使用数组实现的队列
	ArrayQueueUnsafe
	// 使用数组实现的队列(并发安全)
	ArrayQueueSafe
	// 自定义队列
	CustomizeQueue
)

// 函数映射表
var newQueueFuncArr = make([]func(maxSize int) (c IMessageContextQueue, err error), 16, 16)

// 创建队列实例
func (m QueueMode) NewContextQueue(maxSize int) (c IMessageContextQueue, err error) {
	if v := newQueueFuncArr[m]; nil != v {
		return v(maxSize)
	} else {
		return nil, ErrQueueModeUnregister
	}
}

// 根据创建队列实例
func NewContextQueue(id string, mode QueueMode, queueSize int, locateId string, formats []string) (q IMessageContextQueue, err error) {
	q, err = mode.NewContextQueue(queueSize)
	if nil != err {
		return nil, err
	}
	q.SetId(id)
	q.SetLocateId(locateId)
	q.SetFormats(formats)
	return
}

// 注册
func RegisterQueueMode(m QueueMode, f func(maxSize int) (c IMessageContextQueue, err error)) {
	newQueueFuncArr[m] = f
}

func init() {
	RegisterQueueMode(ChannelQueue, NewChannelQueue)
	RegisterQueueMode(ArrayQueueUnsafe, NewUnsafeArrayQueue)
	RegisterQueueMode(ArrayQueueSafe, NewSafeArrayQueue)
}
