package producer

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
	"testing"
	"time"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith struct {
	Producer IRPCMessageProducer
}

func (t *Arith) Add(args Args, reply *Reply) error {
	reply.C = args.A + args.B
	msg := message.NewMessageContext("", "", nil, reply.C)
	t.Producer.NotifyMessageProduced(msg)
	return nil
}

func TestNewRPCMessageProducer(t *testing.T) {
	producer := NewRPCMessageProducer()
	producer.InitRPCServer()
	RObject := &Arith{Producer: producer}
	producer.Register(RObject)
	producer.AddEventListener(EventMessageOnProducer, onRPCProduced)
	producer.StartRPCListener("127.0.0.1:9000")

	client := netx.NewRPCClient(netx.RpcNetworkTCP)
	client.Dial("127.0.0.1:9000")
	args := &Args{}
	reply := new(Reply)
	var A, B int
	for {
		A, B = A+1, B+2
		args.A, args.B = A, B
		client.Call("Arith.Add", args, reply)
		fmt.Println(reply.C)
		time.Sleep(time.Second)
	}
}
func onRPCProduced(evt *eventx.EventData) {
	fmt.Println(evt.Data)
}
