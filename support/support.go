package support

import (
	"github.com/xuzhuoxi/infra-go/lang/collectionx"
)

type IIdSupport collectionx.IOrderHashElement

type ILocateSupport interface {
	LocateKey() string
	SetLocateKey(locateKey string)
}

type IFormatsSupport interface {
	Formats() []string
	SetFormat(formats []string)
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
	id        string
	locateKey string
	formats   []string
}

func (c *ElementSupport) Id() string {
	return c.id
}

func (c *ElementSupport) SetId(Id string) {
	c.id = Id
}

func (c *ElementSupport) LocateKey() string {
	return c.locateKey
}

func (c *ElementSupport) SetLocateKey(locateKey string) {
	c.locateKey = locateKey
}

func (c *ElementSupport) Formats() []string {
	return c.formats
}

func (c *ElementSupport) SetFormat(formats []string) {
	c.formats = formats
}
