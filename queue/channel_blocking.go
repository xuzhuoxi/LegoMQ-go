package queue

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
	"sync"
)

func NewChannelBlockingQueue(maxSize int) (c IMessageContextQueue, err error) {
	if maxSize < 1 {
		return nil, ErrQueueSize
	}
	return &channelBlockingCache{channel: make(chan message.IMessageContext, maxSize)}, nil
}

//---------------------------------

type channelBlockingCache struct {
	support.ElementSupport

	channel chan message.IMessageContext
	closed  bool
	mu      sync.RWMutex
}

func (c *channelBlockingCache) MaxSize() int {
	return cap(c.channel)
}

func (c *channelBlockingCache) Size() int {
	return len(c.channel)
}

func (c *channelBlockingCache) WriteContext(ctx message.IMessageContext) error {
	if nil == ctx {
		return ErrQueueMessageNil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.closed {
		return ErrQueueClosed
	}
	c.channel <- ctx
	return nil
}

func (c *channelBlockingCache) WriteContexts(ctx []message.IMessageContext) (count int, err error) {
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
		c.channel <- ctx[count]
	}
	return ctxLen, nil
}

func (c *channelBlockingCache) ReadContext() (ctx message.IMessageContext, err error) {
	ctx = <-c.channel
	if nil == ctx {
		return nil, ErrQueueClosed
	}
	return
}

func (c *channelBlockingCache) ReadContexts(count int) (ctx []message.IMessageContext, err error) {
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
	return
}

func (c *channelBlockingCache) ReadContextsTo(ctx []message.IMessageContext) (count int, err error) {
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
	return
}

func (c *channelBlockingCache) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	close(c.channel)
}
