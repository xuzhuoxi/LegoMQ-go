package routing

// 消息路由策略
type RoutingMode uint8

func (gm RoutingMode) NewGroupStrategy() (s IRoutingStrategy, err error) {
	if v, ok := groupStrategyMap[gm]; ok {
		return v(), nil
	} else {
		return nil, ErrRougingRegister
	}
}

const (
	// 无路由
	NoneRouting RoutingMode = iota
	// 顺序路由
	SequenceRouting
	// 随机路由
	RandomRouting
	// Hash路由
	HashRouting
	// 自定义路由
	CustomizeRouting
)

var (
	// 函数映射表
	groupStrategyMap = make(map[RoutingMode]func() IRoutingStrategy)
)

// 路由策略接口
type IRoutingStrategy interface {
	// 策略类型
	Mode() RoutingMode
	// 路由目标数量
	LenTarget() int

	// 设置路由目标
	SetTargets(targetIds []string)
	// 追加路由目标
	AppendTarget(targetId string)

	// 路由函数
	RouteTo(routingKey string, targetIds []string) (count int, err error)
	// 路由函数
	Route(routingKey string) (targetIds []string, err error)
}

// 路由策略注册入口
func RegisterRoutingStrategy(s RoutingMode, f func() IRoutingStrategy) {
	groupStrategyMap[s] = f
}
