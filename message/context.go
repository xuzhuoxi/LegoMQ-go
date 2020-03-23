package message

import (
	"fmt"
	"sync"
	"time"
)

var index = 0
var mu sync.Mutex

func NewMessageContext(routingKey string, header interface{}, sender string, receiver string, body interface{}) IMessageContext {
	rs := &messageContext{routingKey: routingKey, header: header, sender: sender, receiver: receiver, body: body}
	rs.timestamp = time.Now().UnixNano()
	mu.Lock()
	rs.index = index
	index += 1
	mu.Unlock()
	return rs
}

//-------------------------------

type messageContext struct {
	routingKey string

	header    interface{}
	sender    string
	receiver  string
	timestamp int64
	index     int

	body interface{}
}

func (c *messageContext) String() string {
	return fmt.Sprintf("[MessageContext]{routingKey=%s,header=%s,sender=%s,receiver=%s,timestamp=%d,index=%d,body=%s}",
		c.routingKey, fmt.Sprint(c.header), c.sender, c.receiver, c.timestamp, c.index, fmt.Sprint(c.body))
	panic("implement me")
}

func (c *messageContext) RoutingKey() string {
	return c.routingKey
}

func (c *messageContext) Timestamp() int64 {
	return c.timestamp
}

func (c *messageContext) Index() int {
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
