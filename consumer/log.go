package consumer

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/logx"
)

func NewLogConsumer() ILogMessageConsumer {
	return newLogConsumer().(ILogMessageConsumer)
}

func newLogConsumer() IMessageConsumer {
	rs := &loggerConsumer{logger: logx.NewLogger()}
	return rs
}

type loggerConsumer struct {
	id     string
	logger logx.ILogger
	level  logx.LogLevel
}

func (c *loggerConsumer) Id() string {
	return c.id
}

func (c *loggerConsumer) SetId(Id string) {
	c.id = Id
}

func (c *loggerConsumer) ConsumeMessage(msg message.IMessageContext) error {
	return c.funcHandleContext(msg)
}

func (c *loggerConsumer) ConsumeMessages(msg []message.IMessageContext) (err error) {
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

func (c *loggerConsumer) SetConsumerLevel(level logx.LogLevel) {
	c.level = level
}

func (c *loggerConsumer) GetLogger() logx.ILogger {
	return c.logger
}

func (c *loggerConsumer) funcHandleContext(msg message.IMessageContext) error {
	if nil == msg {
		return ErrConsumerMessageNil
	}
	switch a := msg.(type) {
	case fmt.Stringer:
		c.logger.Logln(c.level, a.String())
	default:
		c.logger.Logln(c.level, msg.Body())
	}
	return nil
}
