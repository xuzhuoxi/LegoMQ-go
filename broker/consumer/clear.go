package consumer

import (
	"github.com/xuzhuoxi/LegoMQ-go/broker"
	"github.com/xuzhuoxi/LegoMQ-go/message"
)

func NewClearConsumer() broker.IMessageConsumer {
	return newClearConsumer()
}

func newClearConsumer() broker.IMessageConsumer {
	rs := &clearConsumer{loopConsumer{}}
	rs.FuncHandleContext = rs.funcHandleContext
	return rs
}

type clearConsumer struct {
	loopConsumer
}

func (c *clearConsumer) AcceptMessage(msg message.IMessageContext) error {
	return nil
}

func (c *clearConsumer) funcHandleContext(msg message.IMessageContext) error {
	return nil
}
