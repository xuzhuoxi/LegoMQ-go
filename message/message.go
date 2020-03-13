package message

// 消息上下文
type IMessageContext interface {
	// 路由信息
	RoutingKey() string
	// 消息头
	Header() interface{}

	// 发送者信息
	Sender() string
	// 接收者信息
	Receiver() string

	// 生成时间戳
	Timestamp() int64
	// 生成序号
	Index() uint64

	// 消息体
	Body() interface{}
}
