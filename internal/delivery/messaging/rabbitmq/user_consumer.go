package rabbitmq

import (
	"encoding/json"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

// UserConsumer handles user-related RabbitMQ messages
type UserConsumer struct {
	Log *logrus.Logger
}

// NewUserConsumer creates a new UserConsumer instance
func NewUserConsumer(log *logrus.Logger) *UserConsumer {
	return &UserConsumer{
		Log: log,
	}
}

// Consume handles incoming user-related RabbitMQ messages from a queue
func (c *UserConsumer) Consume(delivery amqp091.Delivery) error {
	// Check if logger is available
	if c.Log == nil {
		return nil
	}

	user := new(entity.User)
	if err := json.Unmarshal(delivery.Body, user); err != nil {
		c.Log.WithError(err).Error("error unmarshalling User event")
		return err
	}

	// TODO process event
	c.Log.Infof("Received user event: %v", user)
	return nil
}

// ConsumeWithRouting handles incoming user-related RabbitMQ messages with routing key information
func (c *UserConsumer) ConsumeWithRouting(delivery amqp091.Delivery) error {
	// Check if logger is available
	if c.Log == nil {
		return nil
	}

	user := new(entity.User)
	if err := json.Unmarshal(delivery.Body, user); err != nil {
		c.Log.WithError(err).Error("error unmarshalling User event")
		return err
	}

	// TODO process event based on routing key
	c.Log.Infof("Received user event with routing key %s: %v", delivery.RoutingKey, user)
	return nil
}
