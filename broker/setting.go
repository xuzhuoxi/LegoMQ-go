package broker

import (
	"time"

	"github.com/xuzhuoxi/LegoMQ-go/consumer"
	"github.com/xuzhuoxi/LegoMQ-go/producer"
	"github.com/xuzhuoxi/LegoMQ-go/queue"
	"github.com/xuzhuoxi/LegoMQ-go/routing"
)

// BrokerSetting
// Broker配置信息
type BrokerSetting struct {
	Producers []producer.ProducerSetting // 生产者
	Queues    []queue.QueueSetting       // 队列
	Consumers []consumer.ConsumerSetting // 消息者
	Routing   BrokerRoutingSetting       // 路由
}

// BrokerRoutingSetting
// Broker路由配置信息
type BrokerRoutingSetting struct {
	ProducerRouting routing.RoutingMode `json:"PRouting"` // 队列前置路由(生产者 -> 队列)

	QueueRouting         routing.RoutingMode `json:"QRouting"`  // 队列后置路由(队列 -> 消费者)
	QueueRoutingDuration time.Duration       `json:"QDuration"` // 队列后置路由时间片
	QueueRoutingQuantity int                 `json:"QQuantity"` // 队列后置处理批量
}

func (s BrokerSetting) NewMQBroker() (b IMQBroker, err error) {
	mqConfig, mqBroker := NewMQBroker()
	err = mqConfig.InitBroker(s)
	if nil != err {
		return nil, err
	}
	return mqBroker, nil
}
