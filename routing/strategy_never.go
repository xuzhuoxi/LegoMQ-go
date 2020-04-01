package routing

import "github.com/xuzhuoxi/LegoMQ-go/support"

type neverStrategy struct {
}

func (s *neverStrategy) AppendRoutingTarget(target support.IRoutingTarget) error {
	return ErrRoutingUnSupport
}

func (s *neverStrategy) AppendRoutingTargets(targets []support.IRoutingTarget) error {
	return ErrRoutingUnSupport
}

func (s *neverStrategy) SetRoutingTargets(targets []support.IRoutingTarget) error {
	return ErrRoutingUnSupport
}

func (s *neverStrategy) Mode() RoutingMode {
	return NeverRouting
}

func (s *neverStrategy) Config() IRoutingStrategyConfig {
	return s
}

func (s *neverStrategy) TargetSize() int {
	return 0
}

// 忽略routingKey，locateId
// 命中失败
func (s *neverStrategy) Route(routingKey string, locateId string) (targets []string, err error) {
	return nil, ErrRoutingUnSupport
}

func (s *neverStrategy) match(key string, format string) bool {
	return false
}

//---------------------------------

// 创建一个Never路由策略实例
func NewNeverRoutingStrategy() IRoutingStrategy {
	return &neverStrategy{}
}
