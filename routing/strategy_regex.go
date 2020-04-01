package routing

import (
	"github.com/xuzhuoxi/infra-go/slicex"
	"regexp"
)

type regexStrategy struct {
	StrategyConfig
}

func (s *regexStrategy) Mode() RoutingMode {
	return RegexRouting
}

func (s *regexStrategy) Config() IRoutingStrategyConfig {
	return s
}

// 空字符串忽略
// routingKey和locateId只要命中其中一个，则判定为命中
func (s *regexStrategy) Route(routingKey string, locateId string) (targets []string, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 {
		return nil, ErrRoutingFail
	}
	if "" != routingKey {
		targets = s.regexMath(routingKey, targets)
	}
	if "" != locateId {
		targets = s.regexMath(locateId, targets)
	}
	if 0 != len(targets) {
		slicex.ClearDuplicateString(targets)
		if 0 != len(targets) {
			return targets, nil
		}
	}
	return nil, ErrRoutingFail
}

func (s *regexStrategy) match(routingKey string, routingFormat string) bool {
	matched, _ := regexp.MatchString(routingFormat, routingKey)
	return matched
}

func (s *regexStrategy) regexMath(key string, result []string) []string {
	for idx, _ := range s.Targets {
		formats := s.Targets[idx].Formats()
		if len(formats) == 0 {
			continue
		}
		for _, routingFormat := range formats {
			if s.match(key, routingFormat) {
				result = append(result, s.Targets[idx].Id())
			}
		}
	}
	return result
}

//---------------------

// 创建一个Words路由策略实例
func NewRegexRoutingStrategy() IRoutingStrategy {
	return &regexStrategy{StrategyConfig: StrategyConfig{}}
}
