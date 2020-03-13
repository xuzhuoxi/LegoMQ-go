package consumer

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/broker"
	"github.com/xuzhuoxi/LegoMQ-go/message"
)

func NewPrintConsumer() broker.IMessageConsumer {
	return newPrintConsumer()
}

func newPrintConsumer() broker.IMessageConsumer {
	rs := &printConsumer{loopConsumer{}}
	rs.FuncHandleContext = rs.funcHandleContext
	return rs
}

type printConsumer struct {
	loopConsumer
}

func (c *printConsumer) AcceptMessage(msg message.IMessageContext) error {
	return c.funcHandleContext(msg)
}

func (c *printConsumer) funcHandleContext(msg message.IMessageContext) error {
	switch a := msg.(type) {
	case fmt.Stringer:
		fmt.Println(a.String())
	default:
		fmt.Println(msg.Body())
	}
	return nil
}
