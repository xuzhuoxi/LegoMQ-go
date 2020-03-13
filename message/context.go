package message

import "time"

func NewMessageContext(routingKey string, header interface{}, sender string, receiver string, body interface{}) IMessageContext {
	rs := &messageContext{routingKey: routingKey, header: header, sender: sender, receiver: receiver, body: body}
	rs.timestamp = time.Now().UnixNano()
	return rs
}

//-------------------------------

type messageContext struct {
	routingKey string

	header    interface{}
	sender    string
	receiver  string
	timestamp int64
	index     uint64

	body interface{}
}

func (c *messageContext) RoutingKey() string {
	return c.routingKey
}

func (c *messageContext) Timestamp() int64 {
	return c.timestamp
}

func (c *messageContext) Index() uint64 {
	return c.index
}

func (c *messageContext) Header() interface{} {
	return c.header
}

func (c *messageContext) Sender() string {
	return c.sender
}

func (c *messageContext) Receiver() string {
	return c.receiver
}

func (c *messageContext) Body() interface{} {
	return c.body
}
