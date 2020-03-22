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
	rs := &logConsumer{logger: logx.NewLogger()}
	rs.level = logx.LevelTrace
	rs.logger.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
	return rs
}

type logConsumer struct {
	id      string
	formats []string
	logger  logx.ILogger
	level   logx.LogLevel
}

func (c *logConsumer) Id() string {
	return c.id
}

func (c *logConsumer) SetId(Id string) {
	c.id = Id
}

func (c *logConsumer) Formats() []string {
	return c.formats
}

func (c *logConsumer) SetFormat(formats []string) {
	c.formats = formats
}

func (c *logConsumer) ConsumeMessage(msg message.IMessageContext) error {
	return c.funcHandleContext(msg)
}

func (c *logConsumer) ConsumeMessages(msg []message.IMessageContext) (err error) {
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

func (c *logConsumer) SetConsumerLevel(level logx.LogLevel) {
	c.level = level
}

func (c *logConsumer) GetLogger() logx.ILogger {
	return c.logger
}

func (c *logConsumer) funcHandleContext(msg message.IMessageContext) error {
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
