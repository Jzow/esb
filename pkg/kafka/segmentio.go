package kafka

import (
	"context"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"strings"
	"time"
)

type SegmentioKafkaDriver struct {
	ctx                              context.Context
	producer                         *kafka.Writer
	hosts, username, password, topic string
}

func NewSegmentioKafkaDriver(ctx context.Context, hosts, username, password, topic string) *SegmentioKafkaDriver {
	return &SegmentioKafkaDriver{
		ctx:      ctx,
		hosts:    hosts,
		username: username,
		password: password,
		topic:    topic,
	}
}

func (d *SegmentioKafkaDriver) ProduceMessage(key string, value []byte) {
	beginTime := time.Now()
	if d.producer == nil {
		d.producer = &kafka.Writer{
			Addr:     kafka.TCP(strings.Split(d.hosts, ",")...),
			Topic:    d.topic,
			Balancer: &kafka.Hash{},
		}
		if d.username != "" && d.password != "" {
			d.producer.Transport = &kafka.Transport{SASL: plain.Mechanism{
				Username: d.username,
				Password: d.password,
			}}
		}
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
	err := d.producer.WriteMessages(d.ctx, msg)
	timeStamp := time.Now()
	if msg.Time.Year() > 1 {
		timeStamp = msg.Time
	}
	util.Log("[SegmentioKafkaDriver][ProduceMessage]Topic=%v,Offset=%v,Key=%v,Value=%v,sendTime=%v,cost %v", d.ctx, d.topic, msg.Offset, key, value, timeStamp.UnixMilli(), time.Since(beginTime), err)
}

func (d *SegmentioKafkaDriver) Close() {
	if d.producer != nil {
		d.producer.Close()
		util.Log("[SegmentioKafkaDriver][Close]close producer,Topic=%v,close", d.ctx, d.topic)
	}
}

func (d *SegmentioKafkaDriver) ConsumeMessage(consumerGroupId string, callback func(key string, value []byte, offset int64)) {
	kafkaConfig := kafka.ReaderConfig{
		Brokers:  strings.Split(d.hosts, ","),
		GroupID:  consumerGroupId,
		Topic:    d.topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	}
	if d.username != "" && d.password != "" {
		kafkaConfig.Dialer = &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			SASLMechanism: plain.Mechanism{
				Username: d.username,
				Password: d.password,
			},
		}
	}
	kafkaReader := kafka.NewReader(kafkaConfig)
	defer kafkaReader.Close()
	util.Log("[SegmentioKafkaDriver][ConsumeMessage]host=%v, groupId=%v, topic=%v, start consuming ... !!\n", d.ctx, d.hosts, consumerGroupId, d.topic)
	fmt.Printf("[%v][SegmentioKafkaDriver][ConsumeMessage]host=%s, groupId=%s, topic=%s, start consuming ... !!\n", d.ctx.Value("stack"), d.hosts, consumerGroupId, d.topic)
	for {
		msg, err := kafkaReader.ReadMessage(d.ctx)
		util.Log("[SegmentioKafkaDriver][ConsumeMessage]Topic=%v,Offset=%v,Key=%v,Value=%v,sendTime=%v,currentTime=%v,cost %vms", d.ctx, msg.Topic, msg.Offset, msg.Key, msg.Value, msg.Time.UnixMilli(), time.Now().UnixMilli(), time.Now().Sub(msg.Time).Milliseconds(), err)
		callback(string(msg.Key), msg.Value, msg.Offset)
	}
}
