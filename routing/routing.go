package routing

import (
	"errors"
	"github.com/xuzhuoxi/LegoMQ-go/support"
)

var (
	ErrRougingRegister = errors.New("RoutingStrategy: RoutingMode Unregister! ")

	ErrRoutingUnSupport    = errors.New("RoutingStrategy: Routing is not supported! ")
	ErrRoutingTargetNil    = errors.New("RoutingStrategy: Target is nil! ")
	ErrRoutingTargetsEmpty = errors.New("RoutingStrategy: Targets is empty! ")
	ErrRoutingFail         = errors.New("RoutingStrategy: Route fail! ")
)

// 函数映射表
var newRoutingFuncArr = make([]func() IRoutingStrategy, 32, 32)

// 路由策略模式
type RoutingMode uint8

// 创建策略实例
// err:
//		ErrRougingRegister:	实例化功能未注册
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
	// Hash平均值路由
	HashAvgRouting
	// 单词匹配路由(大小写无关)
	WordsRouting
	// 单词匹配路由(大小写相关)
	CaseWordsRouting
	// 正则路由
	RegexRouting
	// 自定义路由
	CustomizeRouting
)

// 创建策略实例
// mode:	路由策略模式
// err:
//		ErrRougingRegister:	实例化功能未注册
func NewRoutingStrategy(mode RoutingMode) (s IRoutingStrategy, err error) {
	return mode.NewRoutingStrategy()
}

// 路由目标支持接口
// 实现本接口后允许作为目标加入到路由策略中
type IRoutingTarget interface {
	support.IIdSupport
	support.IFormatsSupport
}

// 路由配置支持接口
// 用于处理目标的加入、移除与清空功能
type IRoutingStrategyConfig interface {
	// 路由目标数量
	TargetSize() int
	// 追加路由目标
	AppendRoutingTarget(target IRoutingTarget) error
	// 追加路由目标
	AppendRoutingTargets(targets []IRoutingTarget) error
	// 设置路由目标
	SetRoutingTargets(targets []IRoutingTarget) error
}

// 路由策略接口
type IRoutingStrategy interface {
	// 策略类型
	Mode() RoutingMode
	// 配置入口
	Config() IRoutingStrategyConfig

	// 路由函数
	Route(routingKey string, locateId string) (targetIds []string, err error)
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
