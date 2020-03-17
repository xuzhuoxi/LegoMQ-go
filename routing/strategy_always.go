package routing

type alwaysStrategy struct {
	StrategyConfig
}

func (s *alwaysStrategy) Mode() RoutingMode {
	return NeverRouting
}

func (s *alwaysStrategy) Route(routingKey string) (targets []IRoutingElement, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 {
		return nil, ErrRoutingFail
	}
	targets = make([]IRoutingElement, len(s.Targets), len(s.Targets))
	copy(targets, s.Targets)
	return
}

func (s *alwaysStrategy) match(routingKey string, routingFormat string) bool {
	return true
}

//---------------------------------

// 创建一个Always路由策略实例
func NewAlwaysRoutingStrategy() IRoutingStrategy {
	return &alwaysStrategy{StrategyConfig: StrategyConfig{}}
}
