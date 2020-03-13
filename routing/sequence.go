package routing

import (
	"sync"
)

type sequenceStrategy struct {
	targets []string
	current int
	mu      sync.RWMutex
}

func (s *sequenceStrategy) Mode() RoutingMode {
	return SequenceRouting
}

func (s *sequenceStrategy) LenTarget() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.targets)
}

func (s *sequenceStrategy) SetTargets(targetIds []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.targets = targetIds
}

func (s *sequenceStrategy) AppendTarget(targetId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.targets = append(s.targets, targetId)
}

func (s *sequenceStrategy) RouteTo(routingKey string, targetIds []string) (count int, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.routeTo(routingKey, targetIds)
}

func (s *sequenceStrategy) Route(routingKey string) (targetIds []string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	targetIds = []string{""}
	_, err = s.routeTo(routingKey, targetIds)
	if err != nil {
		return nil, err
	}
	return targetIds, nil
}

func (s *sequenceStrategy) routeTo(routingKey string, targetIds []string) (count int, err error) {
	if len(s.targets) == 0 {
		return -1, ErrRoutingTargetEmpty
	}
	if len(targetIds) < 1 {
		targetIds = append(targetIds, "")
	}
	if s.current >= len(s.targets) {
		s.current = 0
	}
	targetIds[0] = s.targets[s.current]
	s.current += 1
	return 1, nil
}

//---------------------------------

// 创建一个顺序分组的策略实例
func NewSequenceRoutingStrategy() IRoutingStrategy {
	return &sequenceStrategy{}
}

func init() {
	RegisterRoutingStrategy(SequenceRouting, NewSequenceRoutingStrategy)
}
