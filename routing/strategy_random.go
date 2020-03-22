package routing

import (
	"math/rand"
)

type randomStrategy struct {
	StrategyConfig
}

func (s *randomStrategy) Mode() RoutingMode {
	return RandomRouting
}

func (s *randomStrategy) Config() IRoutingStrategyConfig {
	return s
}

func (s *randomStrategy) Route(routingKey string) (targets []string, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 {
		return nil, ErrRoutingFail
	}
	index := rand.Intn(len(s.Targets))
	targets = append(targets, s.Targets[index].Id())
	return
}

func (s *randomStrategy) match(routingKey string, routingFormat string) bool {
	return false
}

//---------------------------------

// 创建一个随机路由的策略实例
func NewRandomRoutingStrategy() IRoutingStrategy {
	return &randomStrategy{StrategyConfig: StrategyConfig{}}
}
