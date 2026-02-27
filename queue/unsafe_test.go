package queue

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/xuzhuoxi/LegoMQ-go/message"
)

func TestUnsafeCache_ReadWriteContext(t *testing.T) {
	c, _ := NewUnsafeArrayQueue(20)
	var mu sync.Mutex
	go func() {
		var err error
		var ctx message.IMessageContext
		for {
			mu.Lock()
			ctx, err = c.ReadContext()
			if nil == err {
				fmt.Println("ReadContext:", err, ctx)
			}
			mu.Unlock()
		}
	}()
	var ctx message.IMessageContext
	for {
		mu.Lock()
		c.WriteContext(ctxDefault)
		ctx = message.NewMessageContext("111", "ddd", nil, "New")
		c.WriteContext(ctx)
		fmt.Println("WriteContext:", ctxDefault, ctx)
		mu.Unlock()
		time.Sleep(time.Second * 2)
	}
}

func TestUnsafeCache_ReadWriteContexts(t *testing.T) {
	c, _ := NewUnsafeArrayQueue(20)
	var mu sync.Mutex
	go func() {
		var err error
		var ctx []message.IMessageContext
		for {
			mu.Lock()
			ctx, err = c.ReadContexts(3)
			if len(ctx) > 0 {
				fmt.Println("ReadContext:", len(ctx), ctx, err)
			}
			mu.Unlock()
		}
	}()
	var err error
	var count int
	for {
		mu.Lock()
		count, err = c.WriteContexts(ctxArr)
		fmt.Println("WriteContexts:", count, err)
		mu.Unlock()
		time.Sleep(time.Second)
	}
}

func TestUnsafeCache_ReadWriteContextsTo(t *testing.T) {
	c, _ := NewUnsafeArrayQueue(20)
	cache := make([]message.IMessageContext, 10, 10)
	var mu sync.Mutex
	go func() {
		var err error
		var count int
		for {
			mu.Lock()
			count, err = c.ReadContextsTo(cache)
			if count > 0 {
				fmt.Println("ReadContext:", count, cache[:count:count], err)
			}
			mu.Unlock()
		}
	}()
	var err error
	var count int
	for {
		mu.Lock()
		count, err = c.WriteContexts(ctxArr)
		fmt.Println("WriteContexts:", count, err)
		mu.Unlock()
		time.Sleep(time.Second)
	}
}
