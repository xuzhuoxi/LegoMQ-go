package queue

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
)

func NewChannelQueue(maxSize int) (c IMessageContextQueue, err error) {
	if maxSize <= 0 {
		return nil, ErrSize
	}
	return &channelCache{channel: make(chan message.IMessageContext, maxSize)}, nil
}

//---------------------------------

type channelCache struct {
	support.ElementSupport

	channel chan message.IMessageContext
}

func (c *channelCache) MaxSize() int {
	return cap(c.channel)
}

func (c *channelCache) Size() int {
	return len(c.channel)
}

func (c *channelCache) WriteContext(ctx message.IMessageContext) error {
	if nil == ctx {
		return ErrQueueMessageNil
	}
	c.channel <- ctx
	return nil
}

func (c *channelCache) WriteContexts(ctx []message.IMessageContext) (count int, err error) {
	ctxLen := len(ctx)
	if 0 == ctxLen {
		return 0, ErrQueueCountZero
	}
	for count = 0; count < len(ctx); count++ {
		if nil == ctx[count] {
			err = ErrQueueMessageNil
			continue
		}
		c.channel <- ctx[count]
	}
	return
}

// 阻塞
func (c *channelCache) ReadContext() (ctx message.IMessageContext, err error) {
	ctx = <-c.channel
	if nil == ctx {
		return nil, ErrQueueClosed
	}
	return
}

// 非阻塞
func (c *channelCache) ReadContexts(count int) (ctx []message.IMessageContext, err error) {
	if count < 1 {
		return nil, ErrQueueCountZero
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

// 非阻塞
func (c *channelCache) ReadContextsTo(ctx []message.IMessageContext) (count int, err error) {
	ctxLen := len(ctx)
	if 0 == ctxLen {
		return 0, ErrQueueCountZero
	}
	for count = 0; count < ctxLen; count++ {
		select {
		case val := <-c.channel:
			if val == nil {
				return count, ErrQueueClosed
			} else {
				ctx[count] = val
			}
		default:
			return count, ErrQueueEmpty
		}
	}
	return
}

func (c *channelCache) Close() {
	close(c.channel)
}
