package queue

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
)

func NewUnsafeArrayQueue(maxSize int) (c IMessageContextQueue, err error) {
	return newUnsafeArrayQueue(maxSize)
}

func newUnsafeArrayQueue(maxSize int) (c *unsafeCache, err error) {
	if maxSize <= 0 {
		return nil, ErrSize
	}
	initialCap := 256
	if maxSize*2 > 256 {
		initialCap = maxSize * 2
	}
	return &unsafeCache{max: maxSize, arr: make([]message.IMessageContext, 0, initialCap)}, nil
}

//---------------------------------

type unsafeCache struct {
	support.ElementSupport
	max int
	arr []message.IMessageContext
}

func (c *unsafeCache) MaxSize() int {
	return c.max
}

func (c *unsafeCache) Size() int {
	return len(c.arr)
}

func (c *unsafeCache) WriteContext(ctx message.IMessageContext) error {
	if nil == ctx {
		return ErrQueueMessageNil
	}
	if len(c.arr) >= c.max {
		return ErrQueueFull
	}
	c.arr = append(c.arr, ctx)
	return nil
}

func (c *unsafeCache) WriteContexts(ctx []message.IMessageContext) (count int, err error) {
	ctxLen := len(ctx)
	if 0 == ctxLen {
		return 0, ErrQueueCountZero
	}
	for idx, _ := range ctx {
		if len(c.arr) == c.max {
			break
		}
		if nil == ctx[idx] {
			err = ErrQueueMessageNil
			continue
		}
		c.arr = append(c.arr, ctx[idx])
		count += 1
	}
	return
}

func (c *unsafeCache) ReadContext() (ctx message.IMessageContext, err error) {
	if len(c.arr) == 0 {
		return nil, ErrQueueEmpty
	}
	ctx = c.arr[0]
	c.arr = c.arr[1:]
	return
}

func (c *unsafeCache) ReadContexts(count int) (ctx []message.IMessageContext, err error) {
	if count <= 0 {
		return nil, ErrQueueCountZero
	}
	cLen := len(c.arr)
	if cLen == 0 {
		return nil, ErrQueueEmpty
	}
	if cLen < count {
		count = cLen
	}
	ctx = make([]message.IMessageContext, count, count)
	copy(ctx, c.arr[:count])
	c.arr = c.arr[count:]
	return
}

func (c *unsafeCache) ReadContextsTo(ctx []message.IMessageContext) (count int, err error) {
	ctxLen := len(ctx)
	if 0 == ctxLen {
		return 0, ErrQueueCountZero
	}
	cLen := len(c.arr)
	if cLen < ctxLen {
		count = cLen
	} else {
		count = ctxLen
	}
	copy(ctx, c.arr[:count])
	c.arr = c.arr[count:]
	return
}

func (c *unsafeCache) Close() {
}
