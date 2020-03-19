package broker

import (
	"github.com/xuzhuoxi/LegoMQ-go/consumer"
	"github.com/xuzhuoxi/LegoMQ-go/producer"
	"github.com/xuzhuoxi/LegoMQ-go/queue"
)

type IMQBrokerConfig interface {
	ProducerConfig() producer.IMessageProducerGroupConfig
	QueueConfig() queue.IMessageQueueGroupConfig
	ConsumerConfig() consumer.IMessageConsumerGroup
}

type IMQBroker interface {
	EngineStart() error
	EngineStop() error
}
