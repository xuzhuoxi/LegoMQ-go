package routing

import (
	"strings"

	"github.com/xuzhuoxi/infra-go/slicex"
)

type wordsStrategy struct {
	StrategyConfig
}

func (s *wordsStrategy) Mode() RoutingMode {
	return WordsRouting
}

func (s *wordsStrategy) Config() IRoutingStrategyConfig {
	return s
}

// Route
// 单词命中路由(大小写无关)
// 空字符串忽略
// routingKey和locateId只要命中其中一个，则判定为命中
func (s *wordsStrategy) Route(routingKey string, locateId string) (targets []string, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 {
		return nil, ErrRoutingFail
	}
	if "" != routingKey {
		targets = s.wordsMath(strings.ToLower(routingKey), targets)
	}
	if "" != locateId {
		targets = s.wordsMath(strings.ToLower(locateId), targets)
	}
	if 0 != len(targets) {
		slicex.ClearDuplicateString(targets)
		if 0 != len(targets) {
			return targets, nil
		}
	}
	return nil, ErrRoutingFail
}

func (s *wordsStrategy) match(lowerKey string, format string) bool {
	return lowerKey == strings.ToLower(format)
}

func (s *wordsStrategy) wordsMath(lowerKey string, result []string) []string {
	for idx, _ := range s.Targets {
		formats := s.Targets[idx].Formats()
		if len(formats) == 0 {
			continue
		}
		for _, routingFormat := range formats {
			if s.match(lowerKey, routingFormat) {
				result = append(result, s.Targets[idx].Id())
				break
			}
		}
	}
	return result
}

//---------------------

type caseWordsStrategy struct {
	StrategyConfig
}

func (s *caseWordsStrategy) Mode() RoutingMode {
	return CaseWordsRouting
}

func (s *caseWordsStrategy) Config() IRoutingStrategyConfig {
	return s
}

// Route
// 单词命中路由(大小写相关)
// 空字符串忽略
// routingKey和locateId只要命中其中一个，则判定为命中
func (s *caseWordsStrategy) Route(routingKey string, locateId string) (targets []string, err error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Targets) == 0 {
		return nil, ErrRoutingFail
	}
	if "" != routingKey {
		targets = s.wordsMath(routingKey, targets)
	}
	if "" != locateId {
		targets = s.wordsMath(locateId, targets)
	}
	if 0 != len(targets) {
		slicex.ClearDuplicateString(targets)
		if 0 != len(targets) {
			return targets, nil
		}
	}
	return nil, ErrRoutingFail
}

func (s *caseWordsStrategy) match(key string, format string) bool {
	return key == format
}

func (s *caseWordsStrategy) wordsMath(lowerKey string, result []string) []string {
	for idx, _ := range s.Targets {
		formats := s.Targets[idx].Formats()
		if len(formats) == 0 {
			continue
		}
		for _, routingFormat := range formats {
			if s.match(lowerKey, routingFormat) {
				result = append(result, s.Targets[idx].Id())
				break
			}
		}
	}
	return result
}

//---------------------

// NewWordsRoutingStrategy
// 创建一个Words路由策略实例
func NewWordsRoutingStrategy() IRoutingStrategy {
	return &wordsStrategy{StrategyConfig: StrategyConfig{}}
}

// NewCaseWordsRoutingStrategy
// 创建一个CaseWords路由策略实例
func NewCaseWordsRoutingStrategy() IRoutingStrategy {
	return &caseWordsStrategy{StrategyConfig: StrategyConfig{}}
}
