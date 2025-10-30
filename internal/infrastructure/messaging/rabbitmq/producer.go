package rabbitmq

import (
	"github.com/sirupsen/logrus"
)

// Producer represents a RabbitMQ producer
type Producer struct {
	client *RabbitMQClient
	logger *logrus.Logger
}

// NewRabbitMqProducer creates a new RabbitMQ producer with a specific logger
func NewRabbitMqProducer(logger *logrus.Logger) *Producer {
	client, err := NewRabbitMQClient()
	if err != nil {
		logger.Errorf("Failed to create RabbitMQ client: %v", err)
		return nil
	}

	return &Producer{
		client: client,
		logger: logger,
	}
}

// PublishToQueue publishes a message directly to a queue
func (p *Producer) PublishToQueue(queueName, message string) error {
	if p == nil || p.client == nil {
		return nil // No producer available, silently ignore
	}
	return p.client.PublishToQueue(queueName, message)
}

// Publish publishes a message to an exchange with a routing key
func (p *Producer) Publish(exchange, routingKey, message string) error {
	if p == nil || p.client == nil {
		return nil // No producer available, silently ignore
	}
	return p.client.Publish(exchange, routingKey, message)
}

// Close closes the RabbitMQ client
func (p *Producer) Close() {
	if p != nil && p.client != nil {
		p.client.Close()
	}
}
