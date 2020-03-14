package mq

import "errors"

var (
	ErrIdUnknown = errors.New("IdGroup: Id Unknown. ")
	ErrIdExists  = errors.New("IdGroup: Id Exists. ")
	ErrIdNil     = errors.New("IdGroup: IdSupport is nil. ")
)

// 唯一标识支持
type IIdSupport interface {
	Id() string
	SetId(Id string)
}

type IIdGroup interface {
	// 数量
	Size() int
	// 信息
	Ids() []string

	// 加入一个Id
	// err:
	//		Id重复时返回 ErrIdExists
	Add(support IIdSupport) error
	// 加入多个Id
	// count: 成功加入的Id数量
	// err2:
	//		Id重复时返回 ErrIdExists
	Adds(supports []IIdSupport) (count int, failArr []IIdSupport, err error)
	// 移除一个Id
	// support: 返回被移除的Id
	// err:
	//		Id不存在时返回 ErrIdUnknown
	Remove(supportId string) (support IIdSupport, err error)
	// 移除多个Id
	// supports: 返回被移除的Id数组
	// err:
	//		Id不存在时返回 ErrIdUnknown
	Removes(supportIdArr []string) (supports []IIdSupport, err error)
	// 替换一个Id
	// 根据Id进行替换，如果找不到相同Id，直接加入
	Update(support IIdSupport)
	// 替换一个Id
	// 根据Id进行替换，如果找不到相同Id，直接加入
	Updates(supports []IIdSupport)
}

type IdGroup struct {
	ids   []IIdSupport
	idMap map[string]IIdSupport
}

func (g *IdGroup) Size() int {
	return len(g.ids)
}

func (g *IdGroup) Ids() []string {
	ln := len(g.ids)
	rs := make([]string, ln, ln)
	for index, q := range g.ids {
		rs[index] = q.Id()
	}
	return rs
}

func (g *IdGroup) Add(support IIdSupport) error {
	if g.isIdExists(support.Id()) {
		return ErrIdExists
	}
	g.add(support)
	return nil

}

func (g *IdGroup) Adds(supports []IIdSupport) (count int, failArr []IIdSupport, err1 error, err2 error) {
	for _, support := range supports {
		if nil == support {
			err1 = ErrIdNil
			continue
		}
		if g.isIdExists(support.Id()) {
			err2 = ErrIdExists
			failArr = append(failArr, support)
			continue
		}
		g.add(support)
		count += 1
	}
	return
}

func (g *IdGroup) Remove(supportId string) (support IIdSupport, err error) {
	if !g.isIdExists(supportId) {
		return nil, ErrIdUnknown
	}
	support = g.removeBy(supportId)
	index := g.findIndex(supportId)
	g.removeAt(index)
	return
}

func (g *IdGroup) Removes(supportIdArr []string) (supports []IIdSupport, err error) {
	for _, sId := range supportIdArr {
		if !g.isIdExists(sId) {
			err = ErrIdUnknown
			continue
		}
		queue := g.removeBy(sId)
		index := g.findIndex(sId)
		g.removeAt(index)
		supports = append(supports, queue)
	}
	return
}

func (g *IdGroup) Update(support IIdSupport) (err error) {
	sId := support.Id()
	if g.isIdExists(sId) {
		index := g.findIndex(sId)
		g.ids[index] = support
	} else {
		g.ids = append(g.ids, support)
	}
	g.idMap[sId] = support
	return nil
}

func (g *IdGroup) Updates(supports []IIdSupport) (err error) {
	for _, support := range supports {
		if nil != support {
			err = ErrIdNil
		}
		qId := support.Id()
		if g.isIdExists(qId) {
			index := g.findIndex(qId)
			g.ids[index] = support
		} else {
			g.ids = append(g.ids, support)
		}
		g.idMap[qId] = support
	}
	return
}

func (g *IdGroup) add(support IIdSupport) {
	g.ids = append(g.ids, support)
	g.idMap[support.Id()] = support
}

func (g *IdGroup) removeAt(index int) (support IIdSupport) {
	if index >= 0 && index < len(g.ids) {
		support = g.ids[index]
		g.ids = append(g.ids[:index], g.ids[index:]...)
	}
	return
}

func (g *IdGroup) removeBy(supportId string) (support IIdSupport) {
	if q, ok := g.idMap[supportId]; ok {
		delete(g.idMap, supportId)
		return q
	} else {
		return nil
	}
}

func (g *IdGroup) findIndex(supportId string) int {
	for index, q := range g.ids {
		if q.Id() == supportId {
			return index
		}
	}
	return -1
}

func (g *IdGroup) isIdExists(supportId string) bool {
	if _, ok := g.idMap[supportId]; ok {
		return true
	}
	return false
}
