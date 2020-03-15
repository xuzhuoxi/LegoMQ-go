package consumer

import "github.com/xuzhuoxi/LegoMQ-go/message"

var (
	msgNil     message.IMessageContext = nil
	msgEmpty                           = message.NewMessageContext("", nil, "", "", nil)
	msgDefault                         = message.NewMessageContext("1", nil, "192.168.1.1", "192.168.1.1", "default")
)
