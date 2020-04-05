package main

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/broker"
	"github.com/xuzhuoxi/LegoMQ-go/consumer"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/LegoMQ-go/producer"
	"github.com/xuzhuoxi/LegoMQ-go/queue"
	"github.com/xuzhuoxi/LegoMQ-go/routing"
	"github.com/xuzhuoxi/infra-go/logx"
	"net/http"
	"time"
)

const (
	HttpPattern      = "/ToLog"
	HttpPattern_Name = "log"
)

var (
	httpProducer producer.IHttpMessageProducer
)

var (
	brokerSetting = broker.BrokerSetting{
		Producers: []producer.ProducerSetting{
			{Id: "P01", Mode: producer.HttpProducer, LocateId: "Http0",
				Http: producer.ProducerSettingHttp{Addr: ":10000"}},
		},
		Queues: []queue.QueueSetting{
			{Id: "Q01", Mode: queue.ChannelBlockingQueue, Size: 256, LocateId: "Queue01", Formats: []string{"Http0"}},
			//{Id: "Q01_Backup", Mode: queue.ChannelBlockingQueue, Size: 256, Formats: []string{"Log02"}},
		},
		Consumers: []consumer.ConsumerSetting{
			{Id: "C01", Mode: consumer.LogConsumer, Formats: []string{"Queue01"},
				Log: consumer.ConsumerSettingLog{Level: logx.LevelTrace,
					Config: []logx.LogConfig{{Type: logx.TypeConsole, Level: logx.LevelAll}}}},
			//{Id: "C01_Backup", Mode: consumer.LogConsumer, Formats: []string{"Log02"}},
		},
		Routing: broker.BrokerRoutingSetting{
			ProducerRouting: routing.AlwaysRouting,
			QueueRouting:    routing.CaseWordsRouting, QueueRoutingDuration: 0, QueueRoutingQuantity: 1,
		},
	}
)

func main() {
	mqConfig, mqBroker := broker.NewMQBroker()
	err := mqConfig.InitBroker(brokerSetting)
	if nil != err {
		fmt.Println(err)
		return
	}

	p, _ := mqConfig.ProducerGroup().GetProducer("P01")
	httpProducer = p.(producer.IHttpMessageProducer)
	httpProducer.HttpServer().MapFunc(HttpPattern, onHttpRequest)

	mqBroker.EngineStart()
	httpProducer.StartProducer()

	fmt.Println("Start.......")
	time.Sleep(time.Hour)
}

func onHttpRequest(w http.ResponseWriter, r *http.Request) {
	logValue := r.FormValue(HttpPattern_Name)
	//fmt.Println("onHttopRequest:", logValue, r.RemoteAddr)
	msg := message.NewMessageContext("", r.RemoteAddr, nil, logValue)
	httpProducer.NotifyMessageProduced(msg)
}
