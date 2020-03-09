package queue

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
)

type unsafeCache struct {
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
	if len(c.arr) >= c.max {
		return cacheFullError
	}
	c.arr = append(c.arr, ctx)
	return nil
}

func (c *unsafeCache) WriteContexts(ctx []message.IMessageContext) (count int, err error) {
	if nil == ctx {
		return 0, inputCtxArrayNilError
	}
	ctxLen := len(ctx)
	if 0 == ctxLen {
		return 0, nil
	}
	cLen := len(c.arr)
	if cLen >= c.max {
		return -1, cacheFullError
	}
	if cLen+ctxLen <= c.max {
		count = ctxLen
	} else {
		count = c.max - cLen
	}
	c.arr = append(c.arr, ctx[:count]...)
	return count, nil
}

func (c *unsafeCache) ReadContext() (ctx message.IMessageContext, err error) {
	if len(c.arr) == 0 {
		return nil, cacheEmptyError
	}
	ctx = c.arr[0]
	c.arr = c.arr[1:]
	return
}

func (c *unsafeCache) ReadContexts(count int) (ctx []message.IMessageContext, err error) {
	cLen := len(c.arr)
	if cLen == 0 {
		return nil, cacheEmptyError
	}
	if cLen < count {
		count = cLen
	}
	ctx = make([]message.IMessageContext, count, count)
	copy(ctx, c.arr[:count:count])
	c.arr = c.arr[count:]
	return
}

func (c *unsafeCache) ReadContextsTo(ctx []message.IMessageContext) (count int, err error) {
	if nil == ctx {
		return -1, outputCtxArrayNilError
	}
	ctxLen := len(ctx)
	if 0 == ctxLen {
		return 0, nil
	}
	cLen := len(c.arr)
	if cLen < ctxLen {
		count = cLen
	} else {
		count = ctxLen
	}
	copy(ctx, c.arr[:count:count])
	c.arr = c.arr[count:]
	return
}

func (c *unsafeCache) Close() {
}

//---------------------------------

func NewUnsafeArrayQueue(maxSize int) (c IContextQueue, err error) {
	if maxSize <= 0 {
		return nil, sizeError
	}
	initialCap := 256
	if maxSize*2 > 256 {
		initialCap = maxSize * 2
	}
	return &unsafeCache{max: maxSize, arr: make([]message.IMessageContext, 0, initialCap)}, nil
}
