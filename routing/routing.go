package routing

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/support"
)

var (
	ErrRougingRegister = errors.New("RoutingMode Unregister! ")

	ErrRoutingUnSupport    = errors.New("RoutingStrategy route failed! ")
	ErrRoutingTargetNil    = errors.New("RoutingStrategy TargetIds is empty! ")
	ErrRoutingTargetsEmpty = errors.New("RoutingStrategy TargetIds is empty! ")
	ErrRoutingFail         = errors.New("RoutingStrategy route failed! ")
)

// 消息路由策略
type RoutingMode uint8

// 函数映射表
var newRoutingFuncArr = make([]func() IRoutingStrategy, 32, 32)

func (m RoutingMode) NewRoutingStrategy() (s IRoutingStrategy, err error) {
	if v := newRoutingFuncArr[m]; nil != v {
		return v(), nil
	} else {
		return nil, ErrRougingRegister
	}
}

const (
	// 无路由
	NeverRouting RoutingMode = iota
	// 全部路由
	AlwaysRouting
	// 顺序路由
	SequenceRouting
	// 随机路由
	RandomRouting
	// Hash路由
	HashAvgRouting
	// 单词匹配路由
	WordsRouting
	// 单词匹配路由
	CaseWordsRouting
	// 正则路由
	RegexRouting
	// 自定义路由
	CustomizeRouting
)

type IRoutingStrategyConfig interface {
	// 路由目标数量
	TargetSize() int
	// 追加路由目标
	AppendRoutingTarget(target support.IRoutingTarget) error
	// 追加路由目标
	AppendRoutingTargets(targets []support.IRoutingTarget) error
	// 设置路由目标
	SetRoutingTargets(targets []support.IRoutingTarget) error
}

// 路由策略接口
type IRoutingStrategy interface {
	// 策略类型
	Mode() RoutingMode
	// 配置入口
	Config() IRoutingStrategyConfig

	// 路由函数
	Route(routingKey string, locateKey string) (targetIds []string, err error)
	// 匹配检查
	match(key string, format string) bool
}

// 路由策略注册入口
func RegisterRoutingStrategy(s RoutingMode, f func() IRoutingStrategy) {
	newRoutingFuncArr[s] = f
}

func init() {
	RegisterRoutingStrategy(NeverRouting, NewNeverRoutingStrategy)
	RegisterRoutingStrategy(AlwaysRouting, NewAlwaysRoutingStrategy)
	RegisterRoutingStrategy(SequenceRouting, NewSequenceRoutingStrategy)
	RegisterRoutingStrategy(RandomRouting, NewRandomRoutingStrategy)
	RegisterRoutingStrategy(HashAvgRouting, NewHashAvgRoutingStrategy)
	RegisterRoutingStrategy(WordsRouting, NewWordsRoutingStrategy)
	RegisterRoutingStrategy(CaseWordsRouting, NewCaseWordsRoutingStrategy)
	RegisterRoutingStrategy(RegexRouting, NewRegexRoutingStrategy)
}
