package consumer

import (
	"fmt"
	"testing"
)

var settings = []ConsumerSetting{
	ConsumerSetting{Id: "1", Mode: ClearConsumer},
	ConsumerSetting{Id: "2", Mode: PrintConsumer},
	ConsumerSetting{Id: "3", Mode: LogConsumer}}

func TestConsumerGroup(t *testing.T) {
	config, group := NewMessageConsumerGroup()
	config.InitConsumerGroup(settings)
	group.ConsumeMessage(msgNil, "1")
	var err error
	err = group.ConsumeMessage(msgNil, "1")
	fmt.Println("Err1:", err)
	err = group.ConsumeMessage(msgEmpty, "2")
	fmt.Println("Err2:", err)
	err = group.ConsumeMessage(msgDefault, "3")
	fmt.Println("Err3:", err)
}
