package queue

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"sync"
)

func NewSafeArrayQueue(maxSize int) (c IMessageContextQueue, err error) {
	us, e := NewUnsafeArrayQueue(maxSize)
	if nil != e {
		return nil, e
	}
	return &safeCache{unsafe: us.(*unsafeCache)}, nil
}

//---------------------------------

type safeCache struct {
	unsafe *unsafeCache

	mu sync.RWMutex
}

func (c *safeCache) MaxSize() int {
	return c.unsafe.MaxSize()
}

func (c *safeCache) Size() int {
	c.mu.RUnlock()
	defer c.mu.RUnlock()
	return c.unsafe.Size()
}

func (c *safeCache) WriteContext(ctx message.IMessageContext) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.unsafe.WriteContext(ctx)
}

func (c *safeCache) WriteContexts(ctx []message.IMessageContext) (count int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.unsafe.WriteContexts(ctx)
}

func (c *safeCache) ReadContext() (ctx message.IMessageContext, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.unsafe.ReadContext()
}

func (c *safeCache) ReadContexts(count int) (ctx []message.IMessageContext, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.unsafe.ReadContexts(count)
}

func (c *safeCache) ReadContextsTo(ctx []message.IMessageContext) (count int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.unsafe.ReadContextsTo(ctx)
}

func (c *safeCache) Close() {
}
