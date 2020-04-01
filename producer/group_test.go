package producer

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"testing"
)

var settings = []ProducerSetting{
	{Id: "1", Mode: HttpProducer, LocateId: "H1"},
	{Id: "2", Mode: RPCProducer, LocateId: "R2"},
	{Id: "3", Mode: SockProducer, LocateId: "S3"}}

func TestConsumerGroup(t *testing.T) {
	config, group := NewMessageProducerGroup()
	config.InitProducerGroup(settings)
	config.SetProducedFunc(func(msg message.IMessageContext, producerId string) {
		fmt.Println(producerId, msg)
	}, nil)

	var err error
	err = group.NotifyMessageProduced(msgNil, "1")
	if nil != err {
		fmt.Println("Err1:", err)
	}
	err = group.NotifyMessageProduced(msgEmpty, "2")
	if nil != err {
		fmt.Println("Err2:", err)
	}
	err = group.NotifyMessageProduced(msgDefault, "3")
	if nil != err {
		fmt.Println("Err3:", err)
	}
	err = group.NotifyMessageProduced(msgDefault, "999")
	if nil != err {
		fmt.Println("Err4:", err)
	}
}
