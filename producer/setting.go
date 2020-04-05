package producer

import (
	"github.com/xuzhuoxi/infra-go/netx"
)

// 消息生产者设置
type ProducerSetting struct {
	Id       string              // 标识
	Mode     ProducerMode        // 消息生产者模式
	LocateId string              // 位置信息
	Http     ProducerSettingHttp // Http消息接收服务器设置
	RPC      ProducerSettingRPC  // RPC消息接收服务器设置
	Sock     netx.SockParams     // Socket消息接收服务器设置
}

func (ps ProducerSetting) NewMessageProducer() (producer IMessageProducer, err error) {
	p, err := ps.Mode.NewMessageProducer()
	if nil != err {
		return nil, err
	}
	p.SetId(ps.Id)
	if sp, ok := p.(IProducerSettingSupport); ok {
		sp.SetSetting(ps)
		err = sp.InitProducer()
		if nil != err {
			return nil, err
		}
	}
	return p, nil
}

// 消息生产者设置
// Http消息接收服务器
type ProducerSettingHttp struct {
	Addr string
}

// 消息生产者设置
// RPC消息接收服务器
type ProducerSettingRPC struct {
	Addr string
}

// 消息生产者设置支持接口
type IProducerSettingSupport interface {
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
