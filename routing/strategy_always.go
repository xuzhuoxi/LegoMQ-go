package routing

type alwaysStrategy struct {
	StrategyConfig
}

func (s *alwaysStrategy) Mode() RoutingMode {
	return NeverRouting
}

func (s *alwaysStrategy) Config() IRoutingStrategyConfig {
	return s
}

func (s *alwaysStrategy) Route(routingKey string) (targetIds []string, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 {
		return nil, ErrRoutingFail
	}
	targetIds = make([]string, len(s.Targets), len(s.Targets))
	for idx, _ := range s.Targets {
		targetIds[idx] = s.Targets[idx].Id()
	}
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
