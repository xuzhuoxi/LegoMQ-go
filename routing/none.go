package routing

type noneStrategy struct {
}

func (s *noneStrategy) Mode() RoutingMode {
	return NoneRouting
}

func (s *noneStrategy) LenTarget() int {
	return 0
}

func (s *noneStrategy) SetTargets(targetIds []string) {
	return
}

func (s *noneStrategy) AppendTarget(targetId string) {
	return
}

func (s *noneStrategy) RouteTo(routingKey string, targetIds []string) (count int, err error) {
	return -1, ErrRoutingFail
}

func (s *noneStrategy) Route(routingKey string) (targetIds []string, err error) {
	return nil, ErrRoutingFail
}

//---------------------------------

// 创建一个None分组的策略实例
func NewNoneRoutingStrategy() IRoutingStrategy {
	return &noneStrategy{}
}

func init() {
	RegisterRoutingStrategy(NoneRouting, NewNoneRoutingStrategy)
}
