package queue

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
)

type channelCache struct {
	channel chan message.IMessageContext
}

func (c *channelCache) MaxSize() int {
	return cap(c.channel)
}

func (c *channelCache) Size() int {
	return len(c.channel)
}

func (c *channelCache) WriteContext(ctx message.IMessageContext) error {
	c.channel <- ctx
	return nil
}

func (c *channelCache) WriteContexts(ctx []message.IMessageContext) (count int, err error) {
	if nil == ctx {
		return 0, inputCtxArrayNilError
	}
	ctxLen := len(ctx)
	if 0 == ctxLen {
		return 0, nil
	}
	for idx, _ := range ctx {
		c.channel <- ctx[idx]
	}
	return ctxLen, nil
}

func (c *channelCache) ReadContext() (ctx message.IMessageContext, err error) {
	ctx = <-c.channel
	if nil == ctx {
		return nil, cacheCloseError
	}
	return
}

func (c *channelCache) ReadContexts(count int) (ctx []message.IMessageContext, err error) {
	if count < 1 {
		return nil, errors.New("Count should >0. ")
	}
	rs := make([]message.IMessageContext, count, count)
	for i := 0; i < count; i++ {
		rs[i] = <-c.channel
		if rs[i] == nil {
			return nil, cacheCloseError
		}
	}
	return rs, nil
}

func (c *channelCache) ReadContextsTo(ctx []message.IMessageContext) (count int, err error) {
	if nil == ctx {
		return -1, outputCtxArrayNilError
	}
	ctxLen := len(ctx)
	if 0 == ctxLen {
		return 0, nil
	}
	for i := 0; i < ctxLen; i++ {
		ctx[i] = <-c.channel
		if ctx[i] == nil {
			return i + 1, cacheCloseError
		}
	}
	return
}

func (c *channelCache) Close() {
	close(c.channel)
}

//---------------------------------

func NewChannelQueue(maxSize int) (c IContextQueue, err error) {
	if maxSize <= 0 {
		return nil, sizeError
	}
	return &channelCache{channel: make(chan message.IMessageContext, 1)}, nil
}
