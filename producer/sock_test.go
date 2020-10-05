package producer

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/netx/tcpx"
	"testing"
	"time"
)

func TestNewSockMessageProducer(t *testing.T) {
	producer := NewSockMessageProducer()
	server, err := producer.InitSockServer(netx.TcpNetwork)
	if nil != err {
		t.Fatal(err)
	}

	var sockHandler = func(data []byte, senderAddress string, other interface{}) (catch bool) {
		msg := message.NewMessageContext("", senderAddress, nil, data)
		producer.NotifyMessageProduced(msg)
		return true
	}

	server.GetPackHandlerContainer().AppendPackHandler(sockHandler)
	producer.AddEventListener(EventMessageOnProducer, onSockProduced)
	err = producer.StartSockListener(netx.SockParams{LocalAddress: "127.0.0.1:9000"})
	if nil != err {
		t.Fatal(err)
	}

	client := tcpx.NewTCPClient()
	err = client.OpenClient(netx.SockParams{RemoteAddress: "127.0.0.1:9000"})
	if nil != err {
		t.Fatal(err)
	}
	for {
		client.SendPackTo([]byte{3, 1, 3, 4})
		time.Sleep(time.Second)
	}
}

func onSockProduced(evt *eventx.EventData) {
	fmt.Println(evt.Data)
}
