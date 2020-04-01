package broker

import (
	"github.com/xuzhuoxi/LegoMQ-go/consumer"
	"github.com/xuzhuoxi/LegoMQ-go/producer"
	"github.com/xuzhuoxi/LegoMQ-go/queue"
	"github.com/xuzhuoxi/LegoMQ-go/routing"
	"time"
)

type BrokerSetting struct {
	Producers []producer.ProducerSetting
	Queues    []queue.QueueSetting
	Consumers []consumer.ConsumerSetting
	Routing   BrokerRoutingSetting
}

type BrokerRoutingSetting struct {
	ProducerRouting routing.RoutingMode
	QueueRouting    routing.RoutingMode
	QueueDuration   time.Duration
}
