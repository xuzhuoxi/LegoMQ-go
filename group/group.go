package group

import "errors"

// 分组策略接口
type IGroupStrategy interface {
	// 策略类型
	Strategy() GroupMode
	// 分组函数
	Group(key string) (index int, err error)
}

// 消息分组策略
type GroupMode uint8

const (
	NoneGroup GroupMode = iota
	// 顺序分组
	SequenceGroup
	// 随机分组
	RandomGroup
	// Hash分组
	HashGroup
)

var (
	groupRegisterErr = errors.New("GroupMode Unregister! ")
	// 函数映射表
	groupStrategyMap = make(map[GroupMode]func(max int) (s IGroupStrategy, err error))
)

func (gm GroupMode) NewGroupStrategy(max int) (s IGroupStrategy, err error) {
	if v, ok := groupStrategyMap[gm]; ok {
		return v(max)
	} else {
		return nil, groupRegisterErr
	}
}

// 分组策略注册入口
func RegisterGroupStrategy(s GroupMode, f func(max int) (s IGroupStrategy, err error)) {
	groupStrategyMap[s] = f
}
