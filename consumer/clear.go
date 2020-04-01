package consumer

import (
	"github.com/xuzhuoxi/LegoMQ-go/message"
)

func NewClearConsumer() IMessageConsumer {
	return newClearConsumer()
}

func newClearConsumer() IMessageConsumer {
	rs := &clearConsumer{}
	return rs
}

type clearConsumer struct {
	id string
}

func (c *clearConsumer) Id() string {
	return c.id
}

func (c *clearConsumer) SetId(Id string) {
	c.id = Id
}

func (c *clearConsumer) Formats() []string {
	return nil
}

func (c *clearConsumer) SetFormats(formats []string) {
	return
}

func (c *clearConsumer) ConsumeMessage(msg message.IMessageContext) error {
	if nil == msg {
		return ErrConsumerMessageNil
	}
	return nil
}

func (c *clearConsumer) ConsumeMessages(msg []message.IMessageContext) error {
	if 0 == len(msg) {
		return ErrConsumerMessagesEmpty
	}
	return nil
}
