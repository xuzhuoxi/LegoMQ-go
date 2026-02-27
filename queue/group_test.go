package queue

import (
	"fmt"
	"testing"
	"time"

	"github.com/xuzhuoxi/LegoMQ-go/message"
)

var (
	channelQueue, _ = NewChannelBlockingQueue(20)
	safeQueue, _    = NewSafeArrayQueue(20)
	ids             = []string{"C01", "S01"}
)

func init() {
	channelQueue.SetId("C01")
	safeQueue.SetId("S01")
}

func TestGroup_ReadWriteContext(t *testing.T) {
	c, g := NewMessageQueueGroup()
	c.AddQueue(channelQueue)
	c.AddQueue(safeQueue)

	go func() {
		var msg = make([]message.IMessageContext, 10, 10)
		for {
			count, err := g.ReadMessagesTo("C01", msg)
			if count > 0 {
				fmt.Println("ReadContext(C01):", count, msg[:count], err)
			}
			time.Sleep(time.Second * 2)
		}
	}()
	go func() {
		for {
			msg, err := g.ReadMessageFrom("S01")
			if nil == err {
				fmt.Println("ReadContext(S01):", msg, err)
			}
		}
	}()
	for {
		fmt.Println("To:")
		fmt.Println(g.WriteMessageTo(ctxDefault, "C01"))
		fmt.Println(g.WriteMessageTo(ctxDefault, "S01"))
		fmt.Println()

		fmt.Println("ToMulti:")
		fmt.Println(g.WriteMessageToMulti(ctxDefault, ids))
		fmt.Println()

		fmt.Println("STo:")
		fmt.Println(g.WriteMessagesTo(ctxArr, "C01"))
		fmt.Println(g.WriteMessagesTo(ctxArr, "S01"))
		fmt.Println()

		fmt.Println("SToMulti:")
		fmt.Println(g.WriteMessagesToMulti(ctxArr, ids))
		fmt.Println(g.WriteMessagesToMulti(ctxArr, ids))
		fmt.Println()
		//
		//time.Sleep(time.Second * 1)
	}
}
