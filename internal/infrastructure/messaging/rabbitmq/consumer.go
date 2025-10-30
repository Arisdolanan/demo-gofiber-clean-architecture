package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

// Consumer represents a RabbitMQ consumer
type Consumer struct {
	client *RabbitMQClient
	logger *logrus.Logger
}

// NewConsumer creates a new RabbitMQ consumer
func NewConsumer() *Consumer {
	logger := logrus.New()

	client, err := NewRabbitMQClient()
	if err != nil {
		logger.Errorf("Failed to create RabbitMQ client: %v", err)
		return nil
	}

	return &Consumer{
		client: client,
		logger: logger,
	}
}

// NewConsumerWithLogger creates a new RabbitMQ consumer with a specific logger
func NewConsumerWithLogger(logger *logrus.Logger) *Consumer {
	client, err := NewRabbitMQClient()
	if err != nil {
		logger.Errorf("Failed to create RabbitMQ client: %v", err)
		return nil
	}

	return &Consumer{
		client: client,
		logger: logger,
	}
}

// ConsumeMessagesFromQueue consumes messages from a RabbitMQ queue
func (c *Consumer) ConsumeMessagesFromQueue(queueName string, handler func(delivery amqp091.Delivery) error) error {
	if c == nil || c.client == nil {
		return nil
	}

	if c.logger != nil {
		c.logger.Infof("Starting to consume messages from RabbitMQ queue: %s", queueName)
	}

	err := c.client.ConsumeMessagesFromQueue(queueName, handler)

	return err
}

// ConsumeMessagesFromExchange consumes messages from a RabbitMQ exchange
func (c *Consumer) ConsumeMessagesFromExchange(exchange, queueName, routingKey string, handler func(delivery amqp091.Delivery) error) error {
	if c == nil || c.client == nil {
		return nil
	}

	if c.logger != nil {
		c.logger.Infof("Starting to consume messages from RabbitMQ exchange: %s with queue: %s and routing key: %s", exchange, queueName, routingKey)
	}

	err := c.client.ConsumeMessagesFromExchange(exchange, queueName, routingKey, handler)

	return err
}

// Close closes the RabbitMQ consumer
func (c *Consumer) Close() {
	if c != nil && c.client != nil {
		if c.logger != nil {
			c.logger.Info("Closing RabbitMQ consumer")
		}
		c.client.Close()
	}
}
