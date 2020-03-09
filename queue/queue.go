package queue

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
)

type IContextReader interface {
	// 向缓存区读出一个消息
	// err: 读取异常
	ReadContext() (ctx message.IMessageContext, err error)
	// 向缓存区读出多个消息
	// ctx: 返回实际读出的消息
	// err: 读取异常
	ReadContexts(count int) (ctx []message.IMessageContext, err error)
	// 向缓存区读出多个消息
	// count: 返回实际读出的消息数量
	// err: 读取异常
	ReadContextsTo(ctx []message.IMessageContext) (count int, err error)
}

type IContextWriter interface {
	// 向缓存区写入一个消息
	// err: 写入异常
	WriteContext(ctx message.IMessageContext) error
	// 向缓存区写入多个消息
	// count: 返回实际写入数量
	// err: 写入异常
	WriteContexts(ctx []message.IMessageContext) (count int, err error)
}

type IContextReadWriter interface {
	IContextReader
	IContextWriter
}

type IContextQueue interface {
	IContextReadWriter
	// Cache最大容量
	MaxSize() int
	// Cache当前容量
	Size() int
	// 关闭
	Close()
}

type QueueMode uint8

const (
	ChannelQueue QueueMode = iota + 1
	ArrayQueueUnsafe
	ArrayQueueSafe
)

var (
	cacheRegisterErr = errors.New("QueueMode Unregister! ")
	// 函数映射表
	cacheMap = make(map[QueueMode]func(maxSize int) (c IContextQueue, err error))
)

func (cm QueueMode) NewContextQueue(maxSize int) (c IContextQueue, err error) {
	if v, ok := cacheMap[cm]; ok {
		return v(maxSize)
	} else {
		return nil, cacheRegisterErr
	}
}

// 分组策略注册入口
func RegisterQueueMode(m QueueMode, f func(maxSize int) (c IContextQueue, err error)) {
	cacheMap[m] = f
}

func init() {
	RegisterQueueMode(ChannelQueue, NewChannelQueue)
	RegisterQueueMode(ArrayQueueUnsafe, NewUnsafeArrayQueue)
	RegisterQueueMode(ArrayQueueSafe, NewSafeArrayQueue)
}
