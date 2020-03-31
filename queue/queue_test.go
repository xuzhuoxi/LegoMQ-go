package queue

import "github.com/xuzhuoxi/LegoMQ-go/message"

var (
	ctxNil     message.IMessageContext
	ctxDefault = message.NewMessageContext("", "", nil, "Default")
	ctxArr     = []message.IMessageContext{message.NewMessageContext("", "", nil, "Default0"),
		message.NewMessageContext("", "", nil, "Default1")}
)
