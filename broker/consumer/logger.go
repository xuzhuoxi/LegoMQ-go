package consumer

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/broker"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/logx"
)

func NewLoggerConsumer() broker.ILogMessageConsumer {
	return newLoggerConsumer().(broker.ILogMessageConsumer)
}

func newLoggerConsumer() broker.IMessageConsumer {
	rs := &loggerConsumer{loopConsumer: loopConsumer{}, logger: logx.NewLogger()}
	rs.FuncHandleContext = rs.funcHandleContext
	return rs
}

type loggerConsumer struct {
	loopConsumer
	logger logx.ILogger
	level  logx.LogLevel
}

func (c *loggerConsumer) AcceptMessage(msg message.IMessageContext) error {
	return c.funcHandleContext(msg)
}

func (c *loggerConsumer) SetConsumerLevel(level logx.LogLevel) {
	c.level = level
}

func (c *loggerConsumer) GetLogger() logx.ILogger {
	return c.logger
}

func (c *loggerConsumer) funcHandleContext(ctx message.IMessageContext) error {
	switch a := ctx.(type) {
	case fmt.Stringer:
		c.logger.Logln(c.level, a.String())
	default:
		c.logger.Logln(c.level, ctx.Body())
	}
	return nil
}
