package routing

type sequenceStrategy struct {
	StrategyConfig
	current int
}

func (s *sequenceStrategy) Mode() RoutingMode {
	return SequenceRouting
}

func (s *sequenceStrategy) Config() IRoutingStrategyConfig {
	return s
}

// Route
// 忽略routingKey，locateId
// 列表循环命中
func (s *sequenceStrategy) Route(routingKey string, locateId string) (targets []string, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 {
		return nil, ErrRoutingFail
	}
	if s.current >= len(s.Targets) {
		s.current = 0
	}
	targets = append(targets, s.Targets[s.current].Id())
	s.current += 1
	return
}

func (s *sequenceStrategy) match(key string, format string) bool {
	return false
}

//---------------------------------

// NewSequenceRoutingStrategy
// 创建一个顺序路由的策略实例
func NewSequenceRoutingStrategy() IRoutingStrategy {
	return &sequenceStrategy{StrategyConfig: StrategyConfig{}}
}
