package producer

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/netx"
	"testing"
)

var settings = []ProducerSetting{
	{Id: "1", Mode: HttpProducer, LocateId: "H1", Http: ProducerSettingHttp{Addr: "", Network: netx.TcpNetwork}},
	{Id: "2", Mode: RPCProducer, LocateId: "R2", RPC: ProducerSettingRPC{Addr: ""}},
	{Id: "3", Mode: SockProducer, LocateId: "S3", Sock: netx.SockParams{Network: netx.TcpNetwork, LocalAddress: ":9000"}}}

func TestConsumerGroup(t *testing.T) {
	var err error

	config, group := NewMessageProducerGroup()
	_, err = config.InitProducerGroup(settings)
	if err != nil {
		fmt.Println(err)
	}
	config.SetProducedFunc(func(msg message.IMessageContext, producerId string) {
		fmt.Println(producerId, msg)
	}, nil)

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
