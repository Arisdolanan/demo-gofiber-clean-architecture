package rabbitmq

import (
	"fmt"
	"log"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient represents a RabbitMQ client
type RabbitMQClient struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

// NewRabbitMQClient creates a new RabbitMQ client
func NewRabbitMQClient() (*RabbitMQClient, error) {
	rabbitConfig := configuration.GetRabbitMQConfig()

	conn, err := amqp091.Dial(rabbitConfig.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
	}, nil
}

// Publish publishes a message to a RabbitMQ exchange
func (r *RabbitMQClient) Publish(exchange, routingKey, message string) error {
	err := r.channel.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	err = r.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Message sent to exchange %s with routing key %s", exchange, routingKey)
	return nil
}

// PublishToQueue publishes a message directly to a queue (without exchange)
func (r *RabbitMQClient) PublishToQueue(queueName, message string) error {
	_, err := r.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Message sent to queue %s", queueName)
	return nil
}

// ConsumeMessagesFromQueue consumes messages from a RabbitMQ queue
func (r *RabbitMQClient) ConsumeMessagesFromQueue(queueName string, handler func(delivery amqp091.Delivery) error) error {
	_, err := r.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	deliveries, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	log.Printf("Started consuming messages from queue %s", queueName)

	for delivery := range deliveries {
		log.Printf("Received message from queue %s", queueName)
		if err := handler(delivery); err != nil {
			log.Printf("Error handling message from queue %s: %v", queueName, err)
		}
	}

	log.Printf("Stopped consuming messages from queue %s", queueName)
	return nil
}

// ConsumeMessagesFromExchange consumes messages from a RabbitMQ exchange
func (r *RabbitMQClient) ConsumeMessagesFromExchange(exchange, queueName, routingKey string, handler func(delivery amqp091.Delivery) error) error {
	err := r.channel.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	queue, err := r.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Bind queue to exchange
	err = r.channel.QueueBind(
		queue.Name, // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	deliveries, err := r.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	log.Printf("Started consuming messages from exchange %s with queue %s and routing key %s", exchange, queueName, routingKey)

	// Process messages
	for delivery := range deliveries {
		log.Printf("Received message from exchange %s with routing key %s", exchange, delivery.RoutingKey)
		if err := handler(delivery); err != nil {
			log.Printf("Error handling message from exchange %s with routing key %s: %v", exchange, delivery.RoutingKey, err)
		}
	}

	log.Printf("Stopped consuming messages from exchange %s with queue %s and routing key %s", exchange, queueName, routingKey)
	return nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQClient) Close() {
	if r.channel != nil {
		r.channel.Close()
		log.Println("RabbitMQ channel closed")
	}
	if r.conn != nil {
		r.conn.Close()
		log.Println("RabbitMQ connection closed")
	}
}
