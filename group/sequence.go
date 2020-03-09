package group

import (
	"errors"
	"sync"
)

var sequenceMaxErr = errors.New("SequenceStrategy Max Error! ")

type sequenceStrategy struct {
	max     int
	current int
	mu      sync.Mutex
}

func (s *sequenceStrategy) Strategy() GroupMode {
	return SequenceGroup
}

func (s *sequenceStrategy) Group(key string) (index int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	index = s.current
	s.current += 1
	if s.current >= s.max {
		s.current = 0
	}
	return
}

//---------------------------------

// 创建一个顺序分组的策略实例
func NewSequenceStrategy(max int) (s IGroupStrategy, err error) {
	if max < 1 {
		return nil, sequenceMaxErr
	} else {
		return &sequenceStrategy{max: max}, nil
	}
}

func init() {
	RegisterGroupStrategy(SequenceGroup, NewSequenceStrategy)
}
