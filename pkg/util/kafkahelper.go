package util

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

func KafkaProduce(pool *sync.Pool, key string, value []byte) error {
	beginTime := time.Now()
	kafkaWriter := pool.Get().(*kafka.Writer)
	defer pool.Put(kafkaWriter)
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
	err := kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		logging.Errorf("[Kafka]failed to produce message,cost %v,topic=%s,key=%s, error: %s", time.Since(beginTime), kafkaWriter.Topic, key, err.Error())
		return err
	}
	logging.Infof("[Kafka]produce message successfully,cost %v,topic=%s,key=%s", time.Since(beginTime), kafkaWriter.Topic, key)
	return nil
}

func KafkaBatchProduce(pool *sync.Pool, keyVals map[string][]byte) error {
	length := len(keyVals)
	if length == 0 {
		return errors.New("keyVals is empty")
	}

	beginTime := time.Now()
	kafkaWriter := pool.Get().(*kafka.Writer)
	defer pool.Put(kafkaWriter)

	keyValDict := make(map[string]any)
	var msgs []kafka.Message
	for key, value := range keyVals {
		keyValDict[key] = string(value)
		msgs = append(msgs, kafka.Message{
			Key:   []byte(key),
			Value: value,
		})
	}
	logBytes, _ := json.Marshal(keyValDict)

	var err error
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = kafkaWriter.WriteMessages(ctx, msgs...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			logging.Errorf("[Kafka]failed to produce message,i=%d, topic=%s error: %s, msg:%s", i, kafkaWriter.Topic, err.Error(), logBytes)
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			logging.Errorf("[Kafka]failed to produce message,i=%d,cost %v,topic=%s, error: %s, msg:%s", i, time.Since(beginTime), kafkaWriter.Topic, err.Error(), logBytes)
			return err
		} else {
			break
		}
	}

	logging.Infof("[Kafka]produce message successfully, topic=%s, cost %v, msg:%s", kafkaWriter.Topic, time.Since(beginTime), logBytes)

	return nil
}
