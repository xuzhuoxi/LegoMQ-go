package producer

import (
	"errors"
	"github.com/xuzhuoxi/infra-go/netx"
)

var ErrSettingNotSupport = errors.New("IProducerSetting not support. ")

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
	if sp, ok := p.(IProducerSetting); ok {
		sp.SetSetting(ps)
		sp.InitProducer()
		return p, nil
	}
	return nil, ErrSettingNotSupport
}

type ProducerSettingHttp struct {
	Addr    string
	Network netx.SockNetwork
}

type ProducerSettingRPC struct {
	Addr string
}

type IProducerSetting interface {
	SetSetting(setting ProducerSetting)
	Setting() ProducerSetting

	InitProducer() error
	StartProducer() error
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
