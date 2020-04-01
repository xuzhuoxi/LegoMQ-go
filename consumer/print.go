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
	id      string
	formats []string
}

func (c *printConsumer) Id() string {
	return c.id
}

func (c *printConsumer) SetId(Id string) {
	c.id = Id
}

func (c *printConsumer) Formats() []string {
	return c.formats
}

func (c *printConsumer) SetFormats(formats []string) {
	c.formats = formats
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

func (c *printConsumer) funcHandleContext(msg message.IMessageContext) error {
	if nil == msg {
		return ErrConsumerMessageNil
	}
	switch a := msg.(type) {
	case fmt.Stringer:
		fmt.Println("Stringer:", a.String())
	default:
		fmt.Println("Body:", msg.Body())
	}
	return nil
}
