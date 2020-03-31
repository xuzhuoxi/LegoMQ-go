package queue

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"testing"
	"time"
)

func TestSafeCache_ReadWriteContext(t *testing.T) {
	c, _ := NewSafeArrayQueue(20)
	go func() {
		var err error
		var ctx message.IMessageContext
		for {
			ctx, err = c.ReadContext()
			if nil == err {
				fmt.Println("ReadContext:", err, ctx)
			}
		}
	}()
	var ctx message.IMessageContext
	for {
		c.WriteContext(ctxDefault)
		ctx = message.NewMessageContext("111", "ddd", nil, "New")
		c.WriteContext(ctx)
		fmt.Println("WriteContext:", ctxDefault, ctx)
		time.Sleep(time.Second * 2)
	}
}

func TestSafeCache_ReadWriteContexts(t *testing.T) {
	c, _ := NewSafeArrayQueue(20)
	go func() {
		var err error
		var ctx []message.IMessageContext
		for {
			ctx, err = c.ReadContexts(1)
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
		time.Sleep(time.Second)
	}
}

func TestSafeCache_ReadWriteContextsTo(t *testing.T) {
	c, _ := NewSafeArrayQueue(20)
	cache := make([]message.IMessageContext, 10, 10)
	go func() {
		var err error
		var count int
		for {
			count, err = c.ReadContextsTo(cache)
			if count > 0 {
				fmt.Println("ReadContext:", count, cache[:count:count], err)
			}
		}
	}()
	var err error
	var count int
	for {
		count, err = c.WriteContexts(ctxArr)
		fmt.Println("WriteContexts:", count, err)
		time.Sleep(time.Second)
	}
}
