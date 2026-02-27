package support

import (
	"github.com/xuzhuoxi/infra-go/lang/collectionx"
)

// IIdSupport
// 标识支持接口
// 提供Id读取与更改功能
// 同范围内实例的Id不允许相同
type IIdSupport collectionx.IOrderHashElement

// ILocateSupport
// 位置支持接口
// 提供LocateId读取与更改功能
// 同范围内实例的LocateId允许相同
type ILocateSupport interface {
	LocateId() string
	SetLocateId(locateId string)
}

// IFormatsSupport
// 格式匹配支持接口
// 用于匹配对照
type IFormatsSupport interface {
	Formats() []string
	SetFormats(formats []string)
}

//------------------------

type ElementSupport struct {
	id       string
	locateId string
	formats  []string
}

func (c *ElementSupport) Id() string {
	return c.id
}

func (c *ElementSupport) SetId(Id string) {
	c.id = Id
}

func (c *ElementSupport) LocateId() string {
	return c.locateId
}

func (c *ElementSupport) SetLocateId(locateId string) {
	c.locateId = locateId
}

func (c *ElementSupport) Formats() []string {
	return c.formats
}

func (c *ElementSupport) SetFormats(formats []string) {
	c.formats = formats
}
