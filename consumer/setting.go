package consumer

import (
	"github.com/xuzhuoxi/infra-go/logx"
)

type ConsumerSetting struct {
	Id   string
	Mode ConsumerMode

	Formats []string

	Log ConsumerSettingLog
}

type ConsumerSettingLog struct {
	Level  logx.LogLevel
	Config []logx.LogConfig
}

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

type IConsumerSettingSupport interface {
	// 设置配置数据
	SetSetting(setting ConsumerSetting)
	// 读取配置数据
	Setting() ConsumerSetting

	// 根据配置数据初始化
	InitConsumer() error
	// 启动
	StartConsumer() error
	// 停止
	StopConsumer() error
}

type ConsumerSettingSupport struct {
	setting ConsumerSetting
}

func (s *ConsumerSettingSupport) SetSetting(setting ConsumerSetting) {
	s.setting = setting
}

func (s *ConsumerSettingSupport) Setting() ConsumerSetting {
	return s.setting
}
