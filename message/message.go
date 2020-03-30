package message

import (
	"fmt"
	"sync"
	"time"
)

var index = 0
var mu sync.Mutex

type IMessageContextHeader interface {
	// 路由信息
	RoutingKey() string
	// 发送者信息
	SenderHost() string
	// 生成时间戳
	Timestamp() int64
	// 生成序号
	Index() int
	// 业务扩展消息头
	Header() interface{}
}

// 消息上下文
type IMessageContext interface {
	IMessageContextHeader
	// 消息体
	Body() interface{}
}

type MessageContextHeader struct {
	routingKey string
	senderHost string
	timestamp  int64
	index      int

	header interface{}
}

func (h *MessageContextHeader) String() string {
	return fmt.Sprintf("Header{RoutingKey=%s,SenderHost=%s,Timestamp=%d,Index=%d,Header=%s}",
		h.routingKey, h.senderHost, h.timestamp, h.index, fmt.Sprint(h.header))
}

func (h *MessageContextHeader) Header() interface{} {
	return h.header
}

func (h *MessageContextHeader) Timestamp() int64 {
	return h.timestamp
}

func (h *MessageContextHeader) Index() int {
	return h.index
}

func (r *MessageContextHeader) RoutingKey() string {
	return r.routingKey
}

func (r *MessageContextHeader) SenderHost() string {
	return r.senderHost
}

func NewMessageContextHeader(routingKey string, senderHost string, header interface{}) MessageContextHeader {
	mu.Lock()
	defer func() {
		index += 1
		mu.Unlock()
	}()
	return MessageContextHeader{routingKey: routingKey, senderHost: senderHost,
		timestamp: time.Now().UnixNano(), index: index, header: header}
}
