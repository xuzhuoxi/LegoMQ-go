package broker

import (
	"github.com/xuzhuoxi/LegoMQ-go/consumer"
	"github.com/xuzhuoxi/LegoMQ-go/producer"
	"github.com/xuzhuoxi/LegoMQ-go/queue"
	"github.com/xuzhuoxi/LegoMQ-go/routing"
	"sync"
	"time"
)

type BrokerRoutingSetting struct {
	ProducerRouting routing.RoutingMode
	QueueRouting    routing.RoutingMode
	QueueSlice      time.Duration
}

type BrokerSetting struct {
	Producers []producer.ProducerSetting
	Queues    []queue.QueueSetting
	Consumers []consumer.ConsumerSetting
	Routing   BrokerRoutingSetting
}

type IMQBrokerConfig interface {
	ProducerGroup() (group producer.IMessageProducerGroup)
	ProducerConfig() (config producer.IMessageProducerGroupConfig)
	QueueGroup() (group queue.IMessageQueueGroup)
	QueueConfig() (config queue.IMessageQueueGroupConfig)
	ConsumerGroup() (group consumer.IMessageConsumerGroup)
	ConsumerConfig() (config consumer.IMessageConsumerGroupConfig)
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

func (b *mqBroker) ProducerGroup() (group producer.IMessageProducerGroup) {
	return b.producerGroup
}

func (b *mqBroker) ProducerConfig() (config producer.IMessageProducerGroupConfig) {
	return b.producerGroup.Config()
}

func (b *mqBroker) QueueGroup() (group queue.IMessageQueueGroup) {
	return b.queueGroup
}

func (b *mqBroker) QueueConfig() (config queue.IMessageQueueGroupConfig) {
	return b.queueGroup.Config()
}

func (b *mqBroker) ConsumerGroup() (group consumer.IMessageConsumerGroup) {
	return b.consumerGroup
}

func (b *mqBroker) ConsumerConfig() (config consumer.IMessageConsumerGroupConfig) {
	return b.consumerGroup.Config()
}

func (b *mqBroker) InitBroker(setting BrokerSetting) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.producerGroup.Config().InitProducerGroup(setting.Producers)
	b.queueGroup.Config().InitQueueGroup(setting.Queues)
	b.consumerGroup.Config().InitConsumerGroup(setting.Consumers)

	b.p2qBridge = NewBridgeProducer2Queue()
	b.p2qBridge.SetEntity(b.producerGroup, b.queueGroup)
	b.p2qBridge.SetRoutingMode(setting.Routing.ProducerRouting)

	b.q2cBridge = NewBridgeQueue2Consumer()
	b.q2cBridge.SetEntity(b.queueGroup, b.consumerGroup)
	b.q2cBridge.SetRoutingMode(setting.Routing.QueueRouting)
	return nil
}

func (b *mqBroker) EngineStart() error {
	b.p2qBridge.Link()
	b.q2cBridge.Link(0, 1)
	return nil
}

func (b *mqBroker) EngineStop() error {
	b.q2cBridge.Unlink()
	b.p2qBridge.Unlink()
	return nil
}
