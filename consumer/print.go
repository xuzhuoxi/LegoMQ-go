package consumer

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/message"
)

func NewPrintConsumer() IMessageConsumer {
	return newPrintConsumer()
}

func newPrintConsumer() IMessageConsumer {
	rs := &printConsumer{}
	return rs
}

type printConsumer struct {
	id string
}

func (c *printConsumer) Id() string {
	return c.id
}

func (c *printConsumer) SetId(Id string) {
	c.id = Id
}

func (c *printConsumer) ConsumeMessage(msg message.IMessageContext) error {
	return c.funcHandleContext(msg)
}

func (c *printConsumer) ConsumeMessages(msg []message.IMessageContext) (err error) {
	if 0 == len(msg) {
		return ErrConsumerMessagesEmpty
	}
	for idx, _ := range msg {
		e := c.funcHandleContext(msg[idx])
		if nil != e {
			err = e
		}
	}
	return
}

func (c *printConsumer) AcceptMessage(msg message.IMessageContext) error {
	return c.funcHandleContext(msg)
}

func (c *printConsumer) funcHandleContext(msg message.IMessageContext) error {
	if nil == msg {
		return ErrConsumerMessageNil
	}
	switch a := msg.(type) {
	case fmt.Stringer:
		fmt.Println(a.String())
	default:
		fmt.Println(msg.Body())
	}
	return nil
}
