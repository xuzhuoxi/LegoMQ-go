package consumer

import (
	"github.com/xuzhuoxi/LegoMQ-go/broker"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/queue"
	"time"
)

type loopConsumer struct {
	ContextQueue queue.IMessageContextQueue
	Started      bool
	Stopping     bool

	count    int
	duration time.Duration

	FuncHandleContext func(msg message.IMessageContext) error
}

func (c *loopConsumer) BindContextQueue(contextQueue queue.IMessageContextQueue) error {
	if c.Started || c.Stopping {
		return broker.ErrConsumerStarted
	}
	if nil == contextQueue {
		return broker.ErrConsumerQueueNil
	}
	c.ContextQueue = contextQueue
	return nil
}

func (c *loopConsumer) UnbindContextQueue() error {
	if c.Started || c.Stopping {
		return broker.ErrConsumerStarted
	}
	if nil == c.ContextQueue {
		return broker.ErrConsumerQueueNil
	}
	c.ContextQueue = nil
	return nil
}

func (c *loopConsumer) StartConsume() error {
	if c.Started {
		return broker.ErrConsumerStarted
	}
	if c.Stopping {
		return broker.ErrConsumerStopping
	}
	c.startLoop()
	return nil
}

func (c *loopConsumer) StopConsume() error {
	if !c.Started {
		return broker.ErrConsumerStopped
	}
	if c.Stopping {
		return broker.ErrConsumerStopping
	}
	c.stopLoop()
	return nil
}

// 总觉得缺点什么
func (c *loopConsumer) startLoop() {
	c.Started = true
	if c.count == 1 {
		c.uniLoop()
	} else {
		c.multiLoop()
	}
	c.Started = false
	c.Stopping = false
}

func (c *loopConsumer) uniLoop() {
	var err error
	var msg message.IMessageContext
	for c.Started && !c.Stopping {
		msg, err = c.ContextQueue.ReadContext()
		if nil != err {
			if err == queue.ErrQueueClosed {
				break
			}
		}
		c.handleContext(msg)
	}
}

func (c *loopConsumer) multiLoop() {
	var err error
	var count int
	var ctx = make([]message.IMessageContext, c.count, c.count)
	for c.Started && !c.Stopping {
		count, err = c.ContextQueue.ReadContextsTo(ctx)
		if nil != err {
			if err == queue.ErrQueueClosed {
				break
			}
		}
		c.handleContexts(ctx[:count:count])
	}
}

// 非即时
func (c *loopConsumer) stopLoop() {
	c.Stopping = true
}

func (c *loopConsumer) handleContext(ctx message.IMessageContext) error {
	return c.FuncHandleContext(ctx)
}

func (c *loopConsumer) handleContexts(ctxArr []message.IMessageContext) (count int, err error) {
	for index, ctx := range ctxArr {
		err = c.FuncHandleContext(ctx)
		if nil != err {
			return index + 1, err
		}
	}
	return len(ctxArr), nil
}
