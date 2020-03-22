package routing

import (
	"crypto/md5"
)

type hashStrategy struct {
	StrategyConfig
}

func (s *hashStrategy) Mode() RoutingMode {
	return HashAvgRouting
}

func (s *hashStrategy) Config() IRoutingStrategyConfig {
	return s
}

func (s *hashStrategy) Route(routingKey string) (targets []string, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 {
		return nil, ErrRoutingFail
	}
	val := md5.Sum([]byte(routingKey))
	rs := int(0)
	for _, v := range val {
		rs += int(v)
	}
	index := rs % len(s.Targets)
	targets = append(targets, s.Targets[index].Id())
	return
}

func (s *hashStrategy) match(routingKey string, routingFormat string) bool {
	return false
}

//---------------------------------

// 创建一个Hash分组的策略实例
func NewHashAvgRoutingStrategy() IRoutingStrategy {
	return &hashStrategy{StrategyConfig: StrategyConfig{}}
}
