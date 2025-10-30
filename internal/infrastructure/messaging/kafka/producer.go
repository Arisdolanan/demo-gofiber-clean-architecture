package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

// KafkaMessage represents a message that can be sent to Kafka
type KafkaMessage interface {
	GetId() string
}

// KafkaProducer wraps a sarama.SyncProducer for easy message sending
type KafkaProducer[T KafkaMessage] struct {
	producer sarama.SyncProducer
	Topic    string
	logger   *logrus.Logger
}

func (kp *KafkaProducer[T]) GetTopic() *string {
	return &kp.Topic
}

// SendMessage sends a message to a Kafka topic
func (kp *KafkaProducer[T]) SendMessage(event T) error {
	value, err := json.Marshal(event)
	if err != nil {
		kp.logger.WithError(err).Error("failed to marshal event")
		return err
	}

	// Create message with topic and content
	msg := &sarama.ProducerMessage{
		Topic: kp.Topic,
		Key:   sarama.StringEncoder(event.GetId()),
		Value: sarama.ByteEncoder(value),
	}

	// Send message and get result
	partition, offset, err := kp.producer.SendMessage(msg)
	if err != nil {
		kp.logger.Errorf("Failed to send message to Kafka topic %s: %v", kp.Topic, err)
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	kp.logger.Infof("Message sent to topic %s, partition %d, offset %d", kp.Topic, partition, offset)
	return nil
}

// Close closes the Kafka producer using the centralized close function
func (kp *KafkaProducer[T]) Close() error {
	return CloseKafkaProducer(kp.producer, kp.logger)
}
