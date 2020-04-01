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

// 忽略locateId
func (s *hashStrategy) Route(routingKey string, locateId string) (targets []string, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 || "" == routingKey {
		return nil, ErrRoutingFail
	}
	index := s.md5Sum(routingKey, len(s.Targets))
	targets = append(targets, s.Targets[index].Id())
	return
}

func (s *hashStrategy) match(key string, format string) bool {
	return false
}

func (s *hashStrategy) md5Sum(key string, max int) int {
	val := md5.Sum([]byte(key))
	rs := int(0)
	for _, v := range val {
		rs += int(v)
	}
	rs = rs % max
	return rs
}

//---------------------------------

// 创建一个Hash分组的策略实例
func NewHashAvgRoutingStrategy() IRoutingStrategy {
	return &hashStrategy{StrategyConfig: StrategyConfig{}}
}
