package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/sirupsen/logrus"
)

// UserConsumer handles user-related Kafka messages
type UserConsumer struct {
	Log *logrus.Logger
}

// NewUserConsumer creates a new UserConsumer instance
func NewUserConsumer(log *logrus.Logger) *UserConsumer {
	return &UserConsumer{
		Log: log,
	}
}

// Consume handles incoming user-related Kafka messages
func (c UserConsumer) Consume(message *sarama.ConsumerMessage) error {
	user := new(entity.User)
	if err := json.Unmarshal(message.Value, user); err != nil {
		c.Log.WithError(err).Error("error unmarshalling User event")
		return err
	}

	c.Log.Infof("Received topic users with event: %v from partition %d", user, message.Partition)
	return nil
}
