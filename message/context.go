package message

import "time"

type messageContext struct {
	header interface{}

	sender    string
	receiver  string
	timestamp int64
	index     uint64

	body interface{}
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

func NewMessageContext(header interface{}, sender string, receiver string, body interface{}) IMessageContext {
	rs := &messageContext{header: header, sender: sender, receiver: receiver, body: body}
	rs.timestamp = time.Now().UnixNano()
	return rs
}
