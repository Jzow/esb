package kafka

import (
	"context"
	"runtime"
)

type IKafkaDriver interface {
	ProduceMessage(key string, value []byte)
	ConsumeMessage(consumerGroupId string, callback func(key string, value []byte, offset int64))
	Close()
}

type KafkaProvider struct {
	ctx                              context.Context
	driver                           IKafkaDriver
	hosts, username, password, topic string
}

func NewKafkaProvider(ctx context.Context, hosts, username, password, topic string) *KafkaProvider {
	provider := &KafkaProvider{
		ctx:      ctx,
		hosts:    hosts,
		username: username,
		password: password,
		topic:    topic,
	}
	if runtime.GOOS == "windows" {
		provider.driver = NewSegmentioKafkaDriver(ctx, hosts, username, password, topic)
	} else {
		provider.driver = NewConfluentKafkaDriver(ctx, hosts, username, password, topic)
	}
	return provider
}

func (k *KafkaProvider) ProduceMessage(key string, value []byte) {
	k.driver.ProduceMessage(key, value)
}

func (k *KafkaProvider) ConsumeMessage(consumerGroupId string, callback func(key string, value []byte, offset int64)) {
	k.driver.ConsumeMessage(consumerGroupId, callback)
}

func (k *KafkaProvider) Close() {
	k.driver.Close()
}
