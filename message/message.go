package message

import (
	"fmt"
	"sync"
	"time"
)

var index = 0
var mu sync.Mutex

// IMessageContextHeader
// 消息头接口
type IMessageContextHeader interface {
	// RoutingKey
	// 路由信息
	RoutingKey() string
	// SenderAddr
	// 发送者信息
	SenderAddr() string
	// Timestamp
	// 生成时间戳
	Timestamp() int64
	// Index
	// 生成序号
	Index() int
	// Header
	// 业务扩展消息头
	Header() interface{}
}

// IMessageContext
// 消息上下文
type IMessageContext interface {
	// IMessageContextHeader
	// 消息头
	IMessageContextHeader
	// Body
	// 消息体
	Body() interface{}
}

type MessageContextHeader struct {
	routingKey string
	senderAddr string
	timestamp  int64
	index      int

	header interface{}
}

func (h *MessageContextHeader) String() string {
	return fmt.Sprintf("Header{RoutingKey=%s,SenderAddr=%s,Timestamp=%d,Index=%d,Header=%s}",
		h.routingKey, h.senderAddr, h.timestamp, h.index, fmt.Sprint(h.header))
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

func (r *MessageContextHeader) SenderAddr() string {
	return r.senderAddr
}

func NewMessageContextHeader(routingKey string, senderAddr string, header interface{}) MessageContextHeader {
	mu.Lock()
	defer func() {
		index += 1
		mu.Unlock()
	}()
	return MessageContextHeader{routingKey: routingKey, senderAddr: senderAddr,
		timestamp: time.Now().UnixNano(), index: index, header: header}
}
