package kafka

import (
	"encoding/json"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/sirupsen/logrus"
)

// UserProducer handles user-related Kafka message production
type UserProducer struct {
	producer sarama.SyncProducer
	logger   *logrus.Logger
	topic    string
}

// NewUserProducer creates a new UserProducer instance
func NewUserProducer(producer sarama.SyncProducer, logger *logrus.Logger) *UserProducer {
	return &UserProducer{
		producer: producer,
		logger:   logger,
		topic:    "users",
	}
}

// PublishUserEvent sends a user entity as a Kafka message
func (p *UserProducer) PublishUserEvent(user *entity.User) error {
	if p.producer == nil {
		p.logger.Warn("Kafka producer is not available, message not sent")
		return nil
	}

	loginEvent := &entity.LoginEvent{
		UserID:    user.ID,
		Email:     user.Email,
		Timestamp: user.UpdatedAt,
	}

	userJSON, err := json.Marshal(loginEvent)
	if err != nil {
		p.logger.WithError(err).Error("error marshalling LoginEvent")
		return err
	}

	// Create Kafka message
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(strconv.FormatInt(user.GetID(), 10)),
		Value: sarama.ByteEncoder(userJSON),
	}

	// Send message
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.WithError(err).Errorf("error publishing user event to Kafka topic %s", p.topic)
		return err
	}

	p.logger.Infof("Published user event to Kafka topic %s, partition %d, offset %d: %v", p.topic, partition, offset, user.ID)
	return nil
}

// Close closes the underlying Kafka producer
func (p *UserProducer) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	return nil
}
