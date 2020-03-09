package group

type noneStrategy struct {
}

//---------------------------------

func (s *noneStrategy) Strategy() GroupMode {
	return NoneGroup
}

func (s *noneStrategy) Group(key string) (index int, err error) {
	return 0, nil
}

//---------------------------------

// 创建一个None分组的策略实例
func NewNoneStrategy(max int) (s IGroupStrategy, err error) {
	return &noneStrategy{}, nil
}

func init() {
	RegisterGroupStrategy(NoneGroup, NewNoneStrategy)
}
