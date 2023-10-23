package consume

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

func Consumer(servers string, groupID string, topics []string, readTimeout time.Duration) {
	go consumeSingle(servers, groupID, topics, readTimeout)
}

var logger = logrus.New()

func consumeSingle(servers string, groupID string, topics []string, readTimeout time.Duration) {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)
	c, err := newKafkaConsumer(servers, groupID)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	if err := c.SubscribeTopics(topics, nil); err != nil {
		fmt.Println(err)
	}

	for {
		if readTimeout <= 0 {
			readTimeout = -1
		}

		msg, err := c.ReadMessage(readTimeout)
		if err != nil {
			kafkaErr, ok := err.(kafka.Error)
			if ok {
				if kafkaErr.Code() == kafka.ErrTimedOut {
					if readTimeout == -1 {
						continue
					}
				}
			}
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			return
		}
		// Execute Handler

		if string(msg.Value) != "" {
			var data map[string]any
			json.Unmarshal(msg.Value, &data)
			logger.WithFields(logrus.Fields{
				"topic": *msg.TopicPartition.Topic,
				"value": data,
			}).Info("consume message")
		}

	}
}

func newKafkaConsumer(servers, groupID string) (*kafka.Consumer, error) {
	// Configurations
	// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	config := &kafka.ConfigMap{
		"bootstrap.servers":        servers,
		"group.id":                 groupID,
		"auto.offset.reset":        "earliest",
		"security.protocol":        "plaintext",
		"auto.commit.interval.ms":  500,
		"enable.auto.offset.store": true,
		"socket.keepalive.enable":  true,
	}

	kc, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return kc, err
}
