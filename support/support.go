package support

import (
	"github.com/xuzhuoxi/infra-go/lang/collectionx"
)

type IIdSupport collectionx.IOrderHashElement

type ILocateSupport interface {
	LocateId() string
	SetLocateId(locateId string)
}

type IFormatsSupport interface {
	Formats() []string
	SetFormats(formats []string)
}

type IProducerBase interface {
	IIdSupport
	ILocateSupport
}

type IQueueBase interface {
	IIdSupport
	ILocateSupport
	IFormatsSupport
}

type IConsumerBase interface {
	IIdSupport
	IFormatsSupport
}

type IRoutingTarget interface {
	IIdSupport
	IFormatsSupport
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
