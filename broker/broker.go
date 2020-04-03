package broker

import (
	"github.com/xuzhuoxi/LegoMQ-go/consumer"
	"github.com/xuzhuoxi/LegoMQ-go/producer"
	"github.com/xuzhuoxi/LegoMQ-go/queue"
	"sync"
)

// Broker操作接口
// Broker配置接口
type IMQBrokerConfig interface {
	// producer组操作接口
	ProducerGroup() (group producer.IMessageProducerGroup)
	// producer组配置接口
	ProducerConfig() (config producer.IMessageProducerGroupConfig)
	// queue组操作接口
	QueueGroup() (group queue.IMessageQueueGroup)
	// queue组配置接口
	QueueConfig() (config queue.IMessageQueueGroupConfig)
	// consumer组操作接口
	ConsumerGroup() (group consumer.IMessageConsumerGroup)
	// consumer组配置接口
	ConsumerConfig() (config consumer.IMessageConsumerGroupConfig)
	// producer、queue、consumer间的连接接口
	Bridge() (p2q IBridgeProducer2Queue, q2c IBridgeQueue2Consumer)
	// 初始化Broker
	InitBroker(setting BrokerSetting) error
}

type IMQBroker interface {
	// 启动
	EngineStart() error
	// 停止
	EngineStop() error
}

// 实例化一个MQBroker
// config:	broker配置接口
// broker:	broker操作接口
func NewMQBroker() (config IMQBrokerConfig, broker IMQBroker) {
	b := &mqBroker{}
	_, b.producerGroup = producer.NewMessageProducerGroup()
	_, b.queueGroup = queue.NewMessageQueueGroup()
	_, b.consumerGroup = consumer.NewMessageConsumerGroup()
	b.p2qBridge = NewBridgeProducer2Queue()
	b.q2cBridge = NewBridgeQueue2Consumer()
	return b, b
}

//-------------

type mqBroker struct {
	setting       BrokerSetting
	producerGroup producer.IMessageProducerGroup
	queueGroup    queue.IMessageQueueGroup
	consumerGroup consumer.IMessageConsumerGroup
	p2qBridge     IBridgeProducer2Queue
	q2cBridge     IBridgeQueue2Consumer
	mu            sync.RWMutex
}

func (b *mqBroker) ProducerGroup() (group producer.IMessageProducerGroup) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.producerGroup
}

func (b *mqBroker) ProducerConfig() (config producer.IMessageProducerGroupConfig) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.producerGroup.Config()
}

func (b *mqBroker) QueueGroup() (group queue.IMessageQueueGroup) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.queueGroup
}

func (b *mqBroker) QueueConfig() (config queue.IMessageQueueGroupConfig) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.queueGroup.Config()
}

func (b *mqBroker) ConsumerGroup() (group consumer.IMessageConsumerGroup) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.consumerGroup
}

func (b *mqBroker) ConsumerConfig() (config consumer.IMessageConsumerGroupConfig) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.consumerGroup.Config()
}

func (b *mqBroker) Bridge() (p2q IBridgeProducer2Queue, q2c IBridgeQueue2Consumer) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.p2qBridge, b.q2cBridge
}

func (b *mqBroker) InitBroker(setting BrokerSetting) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.initBroker(setting)
}

func (b *mqBroker) EngineStart() error {
	b.p2qBridge.Link()
	b.q2cBridge.Link()
	return nil
}

func (b *mqBroker) EngineStop() error {
	b.q2cBridge.Unlink()
	b.p2qBridge.Unlink()
	return nil
}

func (b *mqBroker) initBroker(setting BrokerSetting) error {
	b.setting = setting
	b.producerGroup.Config().InitProducerGroup(setting.Producers)
	b.queueGroup.Config().InitQueueGroup(setting.Queues)
	b.consumerGroup.Config().InitConsumerGroup(setting.Consumers)

	b.p2qBridge.SetBridgePier(b.producerGroup, b.queueGroup)
	b.p2qBridge.SetRoutingMode(setting.Routing.ProducerRouting)

	b.q2cBridge.SetBridgePier(b.queueGroup, b.consumerGroup)
	b.q2cBridge.SetRoutingMode(setting.Routing.QueueRouting)
	err := b.q2cBridge.InitDriver(setting.Routing.QueueBatchDuration, setting.Routing.QueueBatchQuantity)
	if nil != err {
		return err
	}
	return nil
}
