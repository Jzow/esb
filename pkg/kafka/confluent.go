package kafka

import (
	"context"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

type ConfluentKafkaDriver struct {
	ctx                              context.Context
	producer                         *kafka.Producer
	hosts, username, password, topic string
}

func NewConfluentKafkaDriver(ctx context.Context, hosts, username, password, topic string) *ConfluentKafkaDriver {
	return &ConfluentKafkaDriver{
		ctx:      ctx,
		hosts:    hosts,
		username: username,
		password: password,
		topic:    topic,
	}
}

func (d *ConfluentKafkaDriver) ProduceMessage(key string, value []byte) {
	beginTime := time.Now()
	if d.producer == nil {
		conf := kafka.ConfigMap{
			// 设置接入点，请通过控制台获取对应Topic的接入点。多个逗号分隔
			"bootstrap.servers": d.hosts,
			// Kafka producer 的 ack 有 3 种机制，分别说明如下：
			// -1 或 all：Broker 在 leader 收到数据并同步给所有 ISR 中的 follower 后，才应答给 Producer 继续发送下一条（批）消息。
			// 这种配置提供了最高的数据可靠性，只要有一个已同步的副本存活就不会有消息丢失。注意：这种配置不能确保所有的副本读写入该数据才返回，
			// 可以配合 Topic 级别参数 min.insync.replicas 使用。
			// 0：生产者不等待来自 broker 同步完成的确认，继续发送下一条（批）消息。这种配置生产性能最高，但数据可靠性最低
			// （当服务器故障时可能会有数据丢失，如果 leader 已死但是 producer 不知情，则 broker 收不到消息）
			// 1： 生产者在 leader 已成功收到的数据并得到确认后再发送下一条（批）消息。
			// 这种配置是在生产吞吐和数据可靠性之间的权衡（如果leader已死但是尚未复制，则消息可能丢失）
			// 用户不显示配置时，默认值为1。用户根据自己的业务情况进行设置
			"acks": 1,
			// 请求发生错误时重试次数，建议将该值设置为大于0，失败重试最大程度保证消息不丢失
			"retries": 0,
			// 发送请求失败时到下一次重试请求之间的时间
			"retry.backoff.ms": 100,
			// producer 网络请求的超时时间。
			"socket.timeout.ms": 6000,
			// 设置客户端内部重试间隔。
			"reconnect.backoff.max.ms": 3000,
		}
		if d.username != "" && d.password != "" {
			// SASL 验证机制类型默认选用 PLAIN
			conf["sasl.mechanism"] = "PLAIN"
			// 在本地配置 ACL 策略。
			conf["security.protocol"] = "SASL_PLAINTEXT"
			// username 是实例 ID + # + 配置的用户名，password 是配置的用户密码。
			conf["sasl.username"] = d.username
			conf["sasl.password"] = d.password
		}
		p, err := kafka.NewProducer(&conf)
		if err != nil {
			util.Log("[ConfluentKafkaDriver][ProduceMessage]Topic=%v,Key=%v,Value=%v,cost %v", d.ctx, d.topic, key, value, time.Since(beginTime), err)
			return
		}
		d.producer = p
	}

	go func() {
		for e := range d.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				util.Log("[ConfluentKafkaDriver][ProduceMessage]Topic=%v,Offset=%v,Key=%v,Value=%v,sendTime=%v,cost %vms", d.ctx, ev.TopicPartition.Topic, ev.TopicPartition.Offset, ev.Key, ev.Value, time.Now().UnixMilli(), time.Since(beginTime).Milliseconds(), ev.TopicPartition.Error)
			}
		}
	}()

	err := d.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &d.topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, nil)

	if err != nil {
		util.Log("[ConfluentKafkaDriver][ProduceMessage]Topic=%v,Key=%v,Value=%v,cost %v", d.ctx, d.topic, key, value, time.Since(beginTime), err)
		return
	}

	// 等待消息传递
	d.producer.Flush(10 * 1000)
}

func (d *ConfluentKafkaDriver) Close() {
	if d.producer != nil {
		d.producer.Close()
		util.Log("[ConfluentKafkaDriver][Close]close producer,Topic=%v,close", d.ctx, d.topic)
	}
}

func (d *ConfluentKafkaDriver) ConsumeMessage(consumerGroupId string, callback func(key string, value []byte, offset int64)) {
	beginTime := time.Now()
	conf := kafka.ConfigMap{
		// 设置接入点，请通过控制台获取对应Topic的接入点。
		"bootstrap.servers": d.hosts,
		// 设置的消息消费组
		"group.id":          consumerGroupId,
		"auto.offset.reset": "earliest",
		// 使用 Kafka 消费分组机制时，消费者超时时间。当 Broker 在该时间内没有收到消费者的心跳时，认为该消费者故障失败，Broker
		// 发起重新 Rebalance 过程。目前该值的配置必须在 Broker 配置group.min.session.timeout.ms=6000和group.max.session.timeout.ms=300000 之间
		"session.timeout.ms": 10000,
	}
	if d.username != "" && d.password != "" {
		// SASL 验证机制类型默认选用 PLAIN
		conf["sasl.mechanism"] = "PLAIN"
		// 在本地配置 ACL 策略。
		conf["security.protocol"] = "SASL_PLAINTEXT"
		// username 是实例 ID + # + 配置的用户名，password 是配置的用户密码。
		conf["sasl.username"] = d.username
		conf["sasl.password"] = d.password
	}
	c, err := kafka.NewConsumer(&conf)
	if err != nil {
		util.Log("[ConfluentKafkaDriver][ConsumeMessage]Topic=%v,consumerGroupId=%v,cost %v", d.ctx, d.topic, consumerGroupId, time.Since(beginTime), err)
		return
	}
	defer c.Close()
	err = c.SubscribeTopics([]string{d.topic}, nil)
	if err != nil {
		util.Log("[ConfluentKafkaDriver][ConsumeMessage]Topic=%v,consumerGroupId=%v,cost %v", d.ctx, d.topic, consumerGroupId, time.Since(beginTime), err)
		return
	}
	util.Log("[ConfluentKafkaDriver][ConsumeMessage]host=%v, groupId=%v, topic=%v, start consuming ... !! \n", d.ctx, d.hosts, consumerGroupId, d.topic)
	fmt.Printf("[%v][ConfluentKafkaDriver][ConsumeMessage]host=%s, groupId=%s, topic=%s, start consuming ... !! \n", d.ctx.Value("stack"), d.hosts, consumerGroupId, d.topic)
	for {
		var msg *kafka.Message
		msg, err = c.ReadMessage(-1)
		if err != nil {
			util.Log("[ConfluentKafkaDriver][ConsumeMessage]msg:%v,cost %v", d.ctx, msg, err)
			continue
		}
		util.Log("[ConfluentKafkaDriver][ConsumeMessage]Topic=%v,Offset=%v,Key=%v,Value=%v,sendTime=%v,currentTime=%v,cost %vms", d.ctx, msg.TopicPartition.Topic, msg.TopicPartition.Offset, msg.Key, msg.Value, msg.Timestamp.UnixMilli(), time.Now().UnixMilli(), time.Now().Sub(msg.Timestamp).Milliseconds(), err)
		callback(string(msg.Key), msg.Value, util.ObjToInt64(msg.TopicPartition.Offset))
	}
}
