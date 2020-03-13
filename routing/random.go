package routing

import (
	"math/rand"
	"sync"
)

type randomStrategy struct {
	targets []string
	mu      sync.RWMutex
}

func (s *randomStrategy) Mode() RoutingMode {
	return RandomRouting
}

func (s *randomStrategy) LenTarget() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.targets)
}

func (s *randomStrategy) SetTargets(targetIds []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.targets = targetIds
}

func (s *randomStrategy) AppendTarget(targetId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.targets = append(s.targets, targetId)
}

func (s *randomStrategy) RouteTo(routingKey string, targetIds []string) (count int, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.routeTo(routingKey, targetIds)
}

func (s *randomStrategy) Route(routingKey string) (targetIds []string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	targetIds = []string{""}
	_, err = s.routeTo(routingKey, targetIds)
	if err != nil {
		return nil, err
	}
	return targetIds, nil
}

func (s *randomStrategy) routeTo(routingKey string, targetIds []string) (count int, err error) {
	if len(s.targets) == 0 {
		return -1, ErrRoutingTargetEmpty
	}
	if len(targetIds) < 1 {
		targetIds = append(targetIds, "")
	}
	index := rand.Intn(len(s.targets))
	targetIds[0] = s.targets[index]
	return 1, nil
}

//---------------------------------

// 创建一个随机分组的策略实例
func NewRandomRoutingStrategy() IRoutingStrategy {
	return &randomStrategy{}
}

func init() {
	RegisterRoutingStrategy(RandomRouting, NewRandomRoutingStrategy)
}
