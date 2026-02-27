package queue

import (
	"sync"

	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
)

func NewChannelNonBlockingQueue(maxSize int) (c IMessageContextQueue, err error) {
	if maxSize < 1 {
		return nil, ErrQueueSize
	}
	return &channelNBlockingCache{channel: make(chan message.IMessageContext, maxSize)}, nil
}

//---------------------------------

type channelNBlockingCache struct {
	support.ElementSupport

	channel chan message.IMessageContext
	closed  bool
	mu      sync.RWMutex
}

func (c *channelNBlockingCache) MaxSize() int {
	return cap(c.channel)
}

func (c *channelNBlockingCache) Size() int {
	return len(c.channel)
}

func (c *channelNBlockingCache) WriteContext(ctx message.IMessageContext) error {
	if nil == ctx {
		return ErrQueueMessageNil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.closed {
		return ErrQueueClosed
	}
	select {
	case c.channel <- ctx:
		return nil
	default:
		return ErrQueueFull
	}
}

func (c *channelNBlockingCache) WriteContexts(ctx []message.IMessageContext) (count int, err error) {
	ctxLen := len(ctx)
	if 0 == ctxLen {
		return 0, ErrQueueMessagesEmpty
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.closed {
		return -1, ErrQueueClosed
	}
	for count = 0; count < len(ctx); count++ {
		if nil == ctx[count] {
			return count, ErrQueueMessageNil
		}
		select {
		case c.channel <- ctx[count]:
		default:
			return count, ErrQueueFull
		}
	}
	return
}

func (c *channelNBlockingCache) ReadContext() (ctx message.IMessageContext, err error) {
	select {
	case ctx = <-c.channel:
		if ctx == nil {
			return nil, ErrQueueClosed
		} else {
			return ctx, nil
		}
	default:
		return nil, ErrQueueEmpty
	}
}

func (c *channelNBlockingCache) ReadContexts(count int) (ctx []message.IMessageContext, err error) {
	if count < 1 {
		return nil, ErrQueueMessagesEmpty
	}
	ctx = make([]message.IMessageContext, count, count)
	for i := 0; i < count; i++ {
		ctx[i], err = c.ReadContext()
		if nil != err {
			return ctx[:i], err
		}
	}
	return ctx, nil
}

func (c *channelNBlockingCache) ReadContextsTo(ctx []message.IMessageContext) (count int, err error) {
	ctxLen := len(ctx)
	if 0 == ctxLen {
		return 0, ErrQueueMessagesEmpty
	}
	for count = 0; count < ctxLen; count++ {
		ctx[count], err = c.ReadContext()
		if nil != err {
			return count, err
		}
	}
	return count, nil
}

func (c *channelNBlockingCache) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	close(c.channel)
}
