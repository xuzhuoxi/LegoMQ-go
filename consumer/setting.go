package consumer

import (
	"github.com/xuzhuoxi/infra-go/logx"
)

// ConsumerSetting
// 消息消费者设置
type ConsumerSetting struct {
	Id      string       // 标识
	Mode    ConsumerMode // 消息生产者模式
	Formats []string     // 格式匹配信息

	Log ConsumerSettingLog //日志记录配置
}

// ConsumerSettingLog
// 日志消息消费者设置
type ConsumerSettingLog struct {
	Level  logx.LogLevel    // 默认日志等级
	Config []logx.LogConfig // 日志配置
}

// NewMessageConsumer
// 创建日志消息消费者实例
// err:
//
//	ErrConsumerModeUnregister: 实例化功能未注册
func (cs ConsumerSetting) NewMessageConsumer() (consumer IMessageConsumer, err error) {
	p, err := cs.Mode.NewMessageConsumer()
	if nil != err {
		return nil, err
	}
	p.SetId(cs.Id)
	if sp, ok := p.(IConsumerSettingSupport); ok {
		sp.SetSetting(cs)
		err = sp.InitConsumer()
		if nil != err {
			return nil, err
		}
	}
	return p, nil
}

// IConsumerSettingSupport
// 消息消费者设置支持接口
type IConsumerSettingSupport interface {
	// SetSetting
	// 设置配置数据
	SetSetting(setting interface{})
	// Setting
	// 读取配置数据
	Setting() interface{}

	// InitConsumer
	// 根据配置数据初始化
	InitConsumer() error
	// StartConsumer
	// 启动
	StartConsumer() error
	// StopConsumer
	// 停止
	StopConsumer() error
}

type ConsumerSettingSupport struct {
	setting ConsumerSetting
}

func (s *ConsumerSettingSupport) SetSetting(setting interface{}) {
	s.setting = setting.(ConsumerSetting)
}

func (s *ConsumerSettingSupport) Setting() interface{} {
	return s.setting
}
