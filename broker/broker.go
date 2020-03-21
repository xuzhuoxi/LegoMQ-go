package broker

import (
	"github.com/xuzhuoxi/LegoMQ-go/consumer"
	"github.com/xuzhuoxi/LegoMQ-go/producer"
	"github.com/xuzhuoxi/LegoMQ-go/queue"
	"github.com/xuzhuoxi/LegoMQ-go/routing"
	"sync"
)

type BrokerSetting struct {
	ProducerGroup   []producer.ProducerSetting
	ProducerRouting routing.RoutingMode
	QueueGroup      []queue.QueueSetting
	QueueRouting    routing.RoutingMode
	ConsumerGroup   []consumer.ConsumerSetting
}

type IMQBrokerConfig interface {
	ProducerConfig() producer.IMessageProducerGroupConfig
	QueueConfig() queue.IMessageQueueGroupConfig
	ConsumerConfig() consumer.IMessageConsumerGroupConfig
	InitBroker(setting BrokerSetting) error
}

type IMQBroker interface {
	EngineStart() error
	EngineStop() error
}

func NewMQBroker() (config IMQBrokerConfig, broker IMQBroker) {
	b := &mqBroker{}
	_, b.producerGroup = producer.NewMessageProducerGroup()
	_, b.queueGroup = queue.NewMessageQueueGroup()
	_, b.consumerGroup = consumer.NewMessageConsumerGroup()
	return b, b
}

//-------------

type mqBroker struct {
	producerGroup producer.IMessageProducerGroup
	queueGroup    queue.IMessageQueueGroup
	consumerGroup consumer.IMessageConsumerGroup
	p2qBridge     IBridgeProducer2Queue
	q2cBridge     IBridgeQueue2Consumer
	mu            sync.RWMutex
}

func (b *mqBroker) ProducerConfig() producer.IMessageProducerGroupConfig {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.producerGroup.Config()
}

func (b *mqBroker) QueueConfig() queue.IMessageQueueGroupConfig {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.queueGroup.Config()
}

func (b *mqBroker) ConsumerConfig() consumer.IMessageConsumerGroupConfig {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.consumerGroup.Config()
}

func (b *mqBroker) InitBroker(setting BrokerSetting) error {
	b.producerGroup.Config().InitProducerGroup(setting.ProducerGroup)
	b.queueGroup.Config().InitQueueGroup(setting.QueueGroup)
	b.consumerGroup.Config().InitConsumerGroup(setting.ConsumerGroup)
}
