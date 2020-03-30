package message

import (
	"fmt"
)

func NewMessageContext(routingKey string, senderHost string, header interface{}, body interface{}) IMessageContext {
	rs := &messageContext{body: body}
	rs.MessageContextHeader = NewMessageContextHeader(routingKey, senderHost, header)
	return rs
}

//-------------------------------

type messageContext struct {
	MessageContextHeader
	body interface{}
}

func (c *messageContext) String() string {
	return fmt.Sprintf("MessageContext{%s,Body=%s}", c.MessageContextHeader.String(), fmt.Sprint(c.body))
}

func (c *messageContext) Body() interface{} {
	return c.body
}
