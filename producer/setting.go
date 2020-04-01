package producer

import (
	"github.com/xuzhuoxi/infra-go/netx"
)

type ProducerSetting struct {
	Id       string
	Mode     ProducerMode
	LocateId string
	Http     ProducerSettingHttp
	RPC      ProducerSettingRPC
	Sock     netx.SockParams
}

func (ps ProducerSetting) NewMessageProducer() (producer IMessageProducer, err error) {
	p, err := ps.Mode.NewMessageProducer()
	if nil != err {
		return nil, err
	}
	p.SetId(ps.Id)
	if sp, ok := p.(IProducerSetting); ok {
		sp.SetSetting(ps)
		err = sp.InitProducer()
		if nil != err {
			return nil, err
		}
	}
	return p, nil
}

type ProducerSettingHttp struct {
	Addr    string
	Network netx.SockNetwork
}

type ProducerSettingRPC struct {
	Addr string
}

type IProducerSetting interface {
	// 设置配置数据
	SetSetting(setting ProducerSetting)
	// 读取配置数据
	Setting() ProducerSetting

	// 根据配置数据初始化
	InitProducer() error
	// 启动
	StartProducer() error
	// 停止
	StopProducer() error
}

type ProducerSettingSupport struct {
	setting ProducerSetting
}

func (s *ProducerSettingSupport) SetSetting(setting ProducerSetting) {
	s.setting = setting
}

func (s *ProducerSettingSupport) Setting() ProducerSetting {
	return s.setting
}
