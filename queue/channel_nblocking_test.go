package queue

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"testing"
)

func TestChannelNBlockingCache_ReadWriteContext(t *testing.T) {
	c, _ := NewChannelNonBlockingQueue(20)

	go func() {
		var err error
		var ctx message.IMessageContext
		for {
			ctx, err = c.ReadContext()
			fmt.Println("ReadContext:", err, ctx)
		}
	}()
	var err error
	var ctx message.IMessageContext
	for {
		err = c.WriteContext(ctxDefault)
		ctx = message.NewMessageContext("111", "ddd", nil, "New")
		err = c.WriteContext(ctx)
		fmt.Println("WriteContext:", err, ctxDefault, ctx)
	}
}

func TestChannelNBlockingCache_ReadWriteContexts(t *testing.T) {
	c, _ := NewChannelNonBlockingQueue(20)

	go func() {
		var err error
		var ctx []message.IMessageContext
		for {
			ctx, err = c.ReadContexts(3)
			if len(ctx) > 0 {
				fmt.Println("ReadContext:", len(ctx), ctx, err)
			}
		}
	}()
	var err error
	var count int
	for {
		count, err = c.WriteContexts(ctxArr)
		fmt.Println("WriteContexts:", count, err)
	}
}

func TestChannelNBlockingCache_ReadWriteContextsTo(t *testing.T) {
	c, _ := NewChannelNonBlockingQueue(2000)
	cache := make([]message.IMessageContext, 10, 10)
	go func() {
		var err error
		var count int
		for {
			count, err = c.ReadContextsTo(cache)
			if count > 0 || nil != err {
				fmt.Println("ReadContext:", count, cache[:count])
			}
			//time.Sleep(time.Second)
		}
	}()
	var err error
	var count int
	for {
		count, err = c.WriteContexts(ctxArr)
		if count > 0 || nil != err {
			fmt.Println("WriteContexts:", count, err)
		}
	}
}
