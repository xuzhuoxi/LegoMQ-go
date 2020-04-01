package consumer

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
	"github.com/xuzhuoxi/infra-go/logx"
)

var (
	ErrConsumerMessageNil     = errors.New("MessageConsumer: Message is nil. ")
	ErrConsumerMessagesEmpty  = errors.New("MessageConsumer: Message array is empty. ")
	ErrConsumerModeUnregister = errors.New("MessageConsumer: ConsumerMode Unregister. ")
)

type IMessageConsumer interface {
	support.IConsumerBase
	// 消费一条消息
	// err:
	//		msg为nil时ErrConsumerMessageNil
	//		其它错误
	ConsumeMessage(msg message.IMessageContext) error
	// 消费多条消息
	// err:
	//		msg长度为0时ErrMessageContextNil
	//		其它错误
	ConsumeMessages(msg []message.IMessageContext) error
}

type ILogMessageConsumer interface {
	IMessageConsumer

	logx.ILoggerGetter
	SetConsumerLevel(level logx.LogLevel)
}

type ConsumerMode int

const (
	ClearConsumer ConsumerMode = iota + 1
	PrintConsumer
	LogConsumer
	CustomizeConsumer
)

// 函数映射表
var newConsumerFuncArr = make([]func() IMessageConsumer, 16, 16)

// 创建消费者实例
func (m ConsumerMode) NewMessageConsumer() (c IMessageConsumer, err error) {
	if v := newConsumerFuncArr[m]; nil != v {
		return v(), nil
	} else {
		return nil, ErrConsumerModeUnregister
	}
}

// 根据创建消费者实例
func NewMessageConsumer(mode ConsumerMode) (c IMessageConsumer, err error) {
	return mode.NewMessageConsumer()
}

// 注册
func RegisterConsumerMode(m ConsumerMode, f func() IMessageConsumer) {
	newConsumerFuncArr[m] = f
}

func init() {
	RegisterConsumerMode(ClearConsumer, newClearConsumer)
	RegisterConsumerMode(PrintConsumer, newPrintConsumer)
	RegisterConsumerMode(LogConsumer, newLogConsumer)
}
