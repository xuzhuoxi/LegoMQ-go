package routing

import "regexp"

type regexStrategy struct {
	StrategyConfig
}

func (s *regexStrategy) Mode() RoutingMode {
	return RegexRouting
}

func (s *regexStrategy) Config() IRoutingStrategyConfig {
	return s
}

func (s *regexStrategy) Route(routingKey string) (targets []string, err error) {
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
				targets = append(targets, s.Targets[idx].Id())
			}
		}
	}
	if nil != targets {
		return targets, nil
	} else {
		return nil, ErrRoutingFail
	}
}

func (s *regexStrategy) match(routingKey string, routingFormat string) bool {
	matched, _ := regexp.MatchString(routingFormat, routingKey)
	return matched
}

//---------------------

// 创建一个Words路由策略实例
func NewRegexRoutingStrategy() IRoutingStrategy {
	return &regexStrategy{StrategyConfig: StrategyConfig{}}
}
