package queue

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
)

func NewChannelQueue(maxSize int) (c IMessageContextQueue, err error) {
	if maxSize <= 0 {
		return nil, ErrSize
	}
	return &channelCache{channel: make(chan message.IMessageContext, 1)}, nil
}

//---------------------------------

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
		return nil, ErrQueueClosed
	}
	return
}

func (c *channelCache) ReadContexts(count int) (ctx []message.IMessageContext, err error) {
	if count < 1 {
		return nil, errors.New("Count should >0. ")
	}
	rs := make([]message.IMessageContext, count, count)
	for i := 0; i < count; i++ {
		select {
		case val := <-c.channel:
			if val == nil {
				return rs[:i], ErrQueueClosed
			} else {
				rs[i] = val
			}
		default:
			return rs[:i], ErrQueueEmpty
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
			return i + 1, ErrQueueClosed
		}
	}
	return
}

func (c *channelCache) Close() {
	close(c.channel)
}
