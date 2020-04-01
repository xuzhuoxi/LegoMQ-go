package consumer

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
	"github.com/xuzhuoxi/infra-go/logx"
)

func NewLogConsumer() ILogMessageConsumer {
	return &logConsumer{}
}

func NewConsoleLogConsumer() ILogMessageConsumer {
	rs := &logConsumer{logger: logx.NewLogger()}
	rs.level = logx.LevelTrace
	rs.logger.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
	return rs
}

func newLogConsumer() IMessageConsumer {
	return &logConsumer{}
}

type logConsumer struct {
	support.ElementSupport
	ConsumerSettingSupport

	logger logx.ILogger
	level  logx.LogLevel
}

func (c *logConsumer) InitConsumer() error {
	if "" == c.setting.Id {
		return errors.New("Id is empty. ")
	}
	c.SetId(c.setting.Id)
	c.SetFormats(c.setting.Formats)

	c.level = c.setting.Log.Level
	c.logger = logx.NewLogger()
	for _, config := range c.setting.Log.Config {
		c.logger.SetConfig(config)
	}
	return nil
}

func (c *logConsumer) StartConsumer() error {
	return nil
}

func (c *logConsumer) StopConsumer() error {
	return nil
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
