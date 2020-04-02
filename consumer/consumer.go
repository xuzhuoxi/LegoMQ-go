package consumer

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/support"
	"github.com/xuzhuoxi/infra-go/logx"
)

var (
	ErrConsumerMessageNil     = errors.New("MessageConsumer: Message is nil. ")
	ErrConsumerMessagesEmpty  = errors.New("MessageConsumer: Messages is empty. ")
	ErrConsumerModeUnregister = errors.New("MessageConsumer: ConsumerMode Unregister. ")
)

// 消息消费者接口
type IMessageConsumer interface {
	support.IIdSupport      // 标识支持
	support.IFormatsSupport // 格式匹配支持

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

// 日志消息消费者接口
type ILogMessageConsumer interface {
	IMessageConsumer        // 消息消费者
	IConsumerSettingSupport // 消息消费者设置支持
	logx.ILoggerGetter      // 日志获取

	// 设置默认日志等级
	SetConsumerLevel(level logx.LogLevel)
}

// 消息消费者模式
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
// err:
//		ErrConsumerModeUnregister: 实例化功能未注册
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
