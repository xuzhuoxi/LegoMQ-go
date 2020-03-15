package producer

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/eventx"
	"testing"
)

var settings = []ProducerSetting{
	ProducerSetting{Id: "1", Mode: HttpProducer},
	ProducerSetting{Id: "2", Mode: RPCProducer},
	ProducerSetting{Id: "3", Mode: SockProducer}}

func TestConsumerGroup(t *testing.T) {
	config, group := NewMessageProducerGroup()
	ps, _ := config.InitProducerGroup(settings)
	for idx, _ := range ps {
		ps[idx].AddEventListener(EventMessageOnProducer, func(evd *eventx.EventData) {
			fmt.Println(evd.Data)
		})
	}

	group.NotifyMessageProduced(msgNil, "1")
	var err error
	err = group.NotifyMessageProduced(msgNil, "1")
	fmt.Println("Err1:", err)
	err = group.NotifyMessageProduced(msgEmpty, "2")
	fmt.Println("Err2:", err)
	err = group.NotifyMessageProduced(msgDefault, "3")
	fmt.Println("Err3:", err)
}
