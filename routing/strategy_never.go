package routing

type neverStrategy struct {
}

func (s *neverStrategy) AppendRoutingTarget(target IRoutingElement) error {
	return ErrRoutingUnSupport
}

func (s *neverStrategy) AppendRoutingTargets(targets []IRoutingElement) error {
	return ErrRoutingUnSupport
}

func (s *neverStrategy) SetRoutingTargets(targets []IRoutingElement) error {
	return ErrRoutingUnSupport
}

func (s *neverStrategy) Mode() RoutingMode {
	return NeverRouting
}

func (s *neverStrategy) TargetSize() int {
	return 0
}

func (s *neverStrategy) Route(routingKey string) (targets []IRoutingElement, err error) {
	return nil, ErrRoutingUnSupport
}

func (s *neverStrategy) match(routingKey string, routingFormat string) bool {
	return false
}

//---------------------------------

// 创建一个Never路由策略实例
func NewNeverRoutingStrategy() IRoutingStrategy {
	return &neverStrategy{}
}
