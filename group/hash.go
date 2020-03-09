package group

import (
	"crypto/md5"
	"errors"
)

var hashMaxErr = errors.New("HashStrategy Max Error! ")

type hashStrategy struct {
	max int
}

//---------------------------------

func (s *hashStrategy) Strategy() GroupMode {
	return HashGroup
}

func (s *hashStrategy) Group(key string) (index int, err error) {
	val := md5.Sum([]byte(key))
	rs := int(0)
	for _, v := range val {
		rs += int(v)
	}
	index = rs % s.max
	return
}

//---------------------------------

// 创建一个Hash分组的策略实例
func NewHashStrategy(max int) (s IGroupStrategy, err error) {
	if max < 1 {
		return nil, hashMaxErr
	} else {
		return &hashStrategy{max: max}, nil
	}
}

func init() {
	RegisterGroupStrategy(HashGroup, NewHashStrategy)
}
