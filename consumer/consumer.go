package consumer

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/lang/collectionx"
	"github.com/xuzhuoxi/infra-go/logx"
)

var (
	ErrConsumerMessageNil     = errors.New("MessageConsumer: Message is nil. ")
	ErrConsumerMessagesEmpty  = errors.New("MessageConsumer: Message array is empty. ")
	ErrConsumerModeUnregister = errors.New("MessageConsumer: ConsumerMode Unregister. ")
)

type IMessageConsumer interface {
	collectionx.IOrderHashElement
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

type ConsumerSetting struct {
	Id   string
	Mode ConsumerMode
}

var (
	// 函数映射表
	consumerMap = make(map[ConsumerMode]func() IMessageConsumer)
)

func (cm ConsumerMode) NewMessageConsumer() (c IMessageConsumer, err error) {
	if v, ok := consumerMap[cm]; ok {
		return v(), nil
	} else {
		return nil, ErrConsumerModeUnregister
	}
}

// 注册
func RegisterConsumerMode(m ConsumerMode, f func() IMessageConsumer) {
	consumerMap[m] = f
}

func init() {
	RegisterConsumerMode(ClearConsumer, newClearConsumer)
	RegisterConsumerMode(PrintConsumer, newPrintConsumer)
	RegisterConsumerMode(LogConsumer, newLogConsumer)
}
