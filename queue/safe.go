package queue

import (
	"sync"

	"github.com/xuzhuoxi/LegoMQ-go/message"
)

func NewSafeArrayQueue(maxSize int) (c IMessageContextQueue, err error) {
	us, e := newUnsafeArrayQueue(maxSize)
	if nil != e {
		return nil, e
	}
	return &safeCache{unsafeCache: *us}, nil
}

//---------------------------------

type safeCache struct {
	unsafeCache

	mu sync.RWMutex
}

func (c *safeCache) MaxSize() int {
	return c.unsafeCache.MaxSize()
}

func (c *safeCache) Size() int {
	c.mu.RUnlock()
	defer c.mu.RUnlock()
	return c.unsafeCache.Size()
}

func (c *safeCache) WriteContext(ctx message.IMessageContext) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.unsafeCache.WriteContext(ctx)
}

func (c *safeCache) WriteContexts(ctx []message.IMessageContext) (count int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.unsafeCache.WriteContexts(ctx)
}

func (c *safeCache) ReadContext() (ctx message.IMessageContext, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.unsafeCache.ReadContext()
}

func (c *safeCache) ReadContexts(count int) (ctx []message.IMessageContext, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.unsafeCache.ReadContexts(count)
}

func (c *safeCache) ReadContextsTo(ctx []message.IMessageContext) (count int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.unsafeCache.ReadContextsTo(ctx)
}

func (c *safeCache) Close() {
}
