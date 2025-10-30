package rabbitmq

import (
	"encoding/json"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/messaging/rabbitmq"
	"github.com/sirupsen/logrus"
)

// UserProducer handles user-related RabbitMQ message production
type UserProducer struct {
	producer *rabbitmq.Producer
	logger   *logrus.Logger
}

// NewUserProducer creates a new UserProducer instance
func NewUserProducer(producer *rabbitmq.Producer, logger *logrus.Logger) *UserProducer {
	return &UserProducer{
		producer: producer,
		logger:   logger,
	}
}

// PublishUser sends a user entity to a RabbitMQ queue
func (p *UserProducer) PublishUser(queueName string, user *entity.User) error {
	// Check if producer is available
	if p.producer == nil {
		p.logger.Warn("RabbitMQ producer is not available, message not sent")
		return nil
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		p.logger.WithError(err).Error("error marshalling User entity")
		return err
	}

	err = p.producer.PublishToQueue(queueName, string(userJSON))
	if err != nil {
		p.logger.WithError(err).Errorf("error publishing user to queue %s", queueName)
		return err
	}

	p.logger.Infof("Published user to queue %s: %v", queueName, user.ID)
	return nil
}

// PublishUserToExchange sends a user entity to a RabbitMQ exchange with a routing key
func (p *UserProducer) PublishUserToExchange(exchange, routingKey string, user *entity.User) error {
	// Check if producer is available
	if p.producer == nil {
		p.logger.Warn("RabbitMQ producer is not available, message not sent")
		return nil
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		p.logger.WithError(err).Error("error marshalling User entity")
		return err
	}

	err = p.producer.Publish(exchange, routingKey, string(userJSON))
	if err != nil {
		p.logger.WithError(err).Errorf("error publishing user to exchange %s with routing key %s", exchange, routingKey)
		return err
	}

	p.logger.Infof("Published user to exchange %s with routing key %s: %v", exchange, routingKey, user.ID)
	return nil
}

// Close closes the underlying RabbitMQ producer
func (p *UserProducer) Close() error {
	if p.producer != nil {
		p.producer.Close()
	}
	return nil
}
