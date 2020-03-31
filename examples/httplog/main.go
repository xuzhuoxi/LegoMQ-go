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
	HttpAddr         = ":10000"
	HttpPattern      = "/ToLog"
	HttpPattern_Name = "log"
)

var (
	httpProducer producer.IHttpMessageProducer
	globalBroker broker.IMQBroker
)

var (
	brokerSetting = broker.BrokerSetting{
		Producers: []producer.ProducerSetting{
			{Id: "P01", Mode: producer.HttpProducer, LocateKey: "Http0", Addr: "127.0.0.1:9000"},
		},
		Queues: []queue.QueueSetting{
			{Id: "Q01", Mode: queue.ChannelQueue, Size: 256, LocateKey: "Queue01", Formats: []string{"Http0"}},
			//{Id: "Q01_Backup", Mode: queue.ChannelQueue, Size: 256, Formats: []string{"Log02"}},
		},
		Consumers: []consumer.ConsumerSetting{
			{Id: "C01", Mode: consumer.LogConsumer, Formats: []string{"Queue01"}},
			//{Id: "C01_Backup", Mode: consumer.LogConsumer, Formats: []string{"Log02"}},
		},
		Routing: broker.BrokerRoutingSetting{
			ProducerRouting: routing.AlwaysRouting, QueueRouting: routing.CaseWordsRouting,
		},
	}
)

func main() {
	mqConfig, mqBroker := broker.NewMQBroker()
	globalBroker = mqBroker
	mqConfig.InitBroker(brokerSetting)

	p, _ := mqConfig.ProducerGroup().GetProducer("P01")
	httpProducer = p.(producer.IHttpMessageProducer)
	httpProducer.HttpServer().MapFunc(HttpPattern, onHttpRequest)
	mqConfig.ProducerConfig().AddProducer(httpProducer)

	c0, _ := mqConfig.ConsumerGroup().GetConsumer("C01")
	l0 := c0.(consumer.ILogMessageConsumer)
	l0.SetConsumerLevel(logx.LevelTrace)
	l0.GetLogger().SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})

	mqBroker.EngineStart()
	httpProducer.StartHttpListener(HttpAddr)
	fmt.Println("Start.......")
	time.Sleep(time.Hour)
}

func onHttpRequest(w http.ResponseWriter, r *http.Request) {
	logValue := r.FormValue(HttpPattern_Name)
	//fmt.Println("onHttopRequest:", logValue, r.RemoteAddr)
	msg := message.NewMessageContext("", r.RemoteAddr, nil, logValue)
	httpProducer.NotifyMessageProduced(msg)
}
