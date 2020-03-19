package routing

import "strings"

type wordsStrategy struct {
	StrategyConfig
}

func (s *wordsStrategy) Mode() RoutingMode {
	return WordsRouting
}

func (s *wordsStrategy) Config() IRoutingStrategyConfig {
	return s
}

func (s *wordsStrategy) Route(routingKey string) (targets []string, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 {
		return nil, ErrRoutingFail
	}
	lowerRoutingKey := strings.ToLower(routingKey)
	for idx, _ := range s.Targets {
		formats := s.Targets[idx].Formats()
		if len(formats) == 0 {
			continue
		}
		for _, routingFormat := range formats {
			if s.match(lowerRoutingKey, routingFormat) {
				targets = append(targets, s.Targets[idx].Id())
				break
			}
		}
	}
	if nil != targets {
		return targets, nil
	} else {
		return nil, ErrRoutingFail
	}
}

func (s *wordsStrategy) match(lowerRoutingKey string, routingFormat string) bool {
	return lowerRoutingKey == strings.ToLower(routingFormat)
}

//---------------------

type caseWordsStrategy struct {
	StrategyConfig
}

func (s *caseWordsStrategy) Mode() RoutingMode {
	return CaseWordsRouting
}

func (s *caseWordsStrategy) Route(routingKey string) (targets []IRoutingElement, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 {
		return nil, ErrRoutingFail
	}
	for idx, _ := range s.Targets {
		formats := s.Targets[idx].Formats()
		if len(formats) == 0 {
			continue
		}
		for _, routingFormat := range formats {
			if s.match(routingKey, routingFormat) {
				targets = append(targets, s.Targets[idx])
				break
			}
		}
	}
	if nil != targets {
		return targets, nil
	} else {
		return nil, ErrRoutingFail
	}
}

func (s *caseWordsStrategy) match(routingKey string, routingFormat string) bool {
	return routingKey == routingFormat
}

//---------------------

// 创建一个Words路由策略实例
func NewWordsRoutingStrategy() IRoutingStrategy {
	return &wordsStrategy{StrategyConfig: StrategyConfig{}}
}

// 创建一个CaseWords路由策略实例
func NewCaseWordsRoutingStrategy() IRoutingStrategy {
	return &caseWordsStrategy{StrategyConfig: StrategyConfig{}}
}
