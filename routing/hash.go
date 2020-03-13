package routing

import (
	"crypto/md5"
	"sync"
)

type hashStrategy struct {
	targets []string
	mu      sync.RWMutex
}

func (s *hashStrategy) Mode() RoutingMode {
	return HashRouting
}

func (s *hashStrategy) LenTarget() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.targets)
}

func (s *hashStrategy) SetTargets(targetIds []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.targets = targetIds
}

func (s *hashStrategy) AppendTarget(targetId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.targets = append(s.targets, targetId)
}

func (s *hashStrategy) RouteTo(routingKey string, targetIds []string) (count int, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.routeTo(routingKey, targetIds)
}

func (s *hashStrategy) Route(routingKey string) (targetIds []string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	targetIds = []string{""}
	_, err = s.routeTo(routingKey, targetIds)
	if err != nil {
		return nil, err
	}
	return targetIds, nil
}

func (s *hashStrategy) routeTo(routingKey string, targetIds []string) (count int, err error) {
	if len(s.targets) == 0 {
		return -1, ErrRoutingTargetEmpty
	}
	if len(targetIds) < 1 {
		targetIds = append(targetIds, "")
	}
	val := md5.Sum([]byte(routingKey))
	rs := int(0)
	for _, v := range val {
		rs += int(v)
	}
	index := rs % len(s.targets)
	targetIds[0] = s.targets[index]
	return 1, nil
}

//---------------------------------

// 创建一个Hash分组的策略实例
func NewHashRoutingStrategy() IRoutingStrategy {
	return &hashStrategy{}
}

func init() {
	RegisterRoutingStrategy(HashRouting, NewHashRoutingStrategy)
}
