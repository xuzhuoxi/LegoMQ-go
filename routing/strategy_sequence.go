package routing

type sequenceStrategy struct {
	StrategyConfig
	current int
}

func (s *sequenceStrategy) Mode() RoutingMode {
	return SequenceRouting
}

func (s *sequenceStrategy) Route(routingKey string) (targets []IRoutingElement, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 {
		return nil, ErrRoutingFail
	}
	if s.current >= len(s.Targets) {
		s.current = 0
	}
	targets = append(nil, s.Targets[s.current])
	s.current += 1
	return
}

func (s *sequenceStrategy) match(routingKey string, routingFormat string) bool {
	return false
}

//---------------------------------

// 创建一个顺序路由的策略实例
func NewSequenceRoutingStrategy() IRoutingStrategy {
	return &sequenceStrategy{StrategyConfig: StrategyConfig{}}
}
