package group

import (
	"errors"
	"math/rand"
)

var randomMaxErr = errors.New("RandomStrategy Max Error! ")

type randomStrategy struct {
	max int
}

func (s *randomStrategy) Strategy() GroupMode {
	return RandomGroup
}

func (s *randomStrategy) Group(key string) (index int, err error) {
	index = rand.Intn(s.max)
	return
}

//---------------------------------

// 创建一个随机分组的策略实例
func NewRandomStrategy(max int) (s IGroupStrategy, err error) {
	if max < 1 {
		return nil, randomMaxErr
	} else {
		return &randomStrategy{max: max}, nil
	}
}

func init() {
	RegisterGroupStrategy(RandomGroup, NewRandomStrategy)
}
