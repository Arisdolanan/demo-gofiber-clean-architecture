package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/sirupsen/logrus"
)

// createKafkaSaramaConfig creates a standardized Kafka configuration
func createKafkaSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()

	// Producer configuration
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all replicas to commit
	config.Producer.Retry.Max = 5                    // Retry up to 5 times on failure
	config.Producer.Return.Successes = true          // Return successes to caller

	// Consumer configuration
	config.Consumer.Return.Errors = true                                        // Return errors to caller
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin // Use round-robin partition assignment
	config.Consumer.Offsets.Initial = sarama.OffsetNewest                       // Start from newest message if no offset found

	return config
}

// NewKafkaConsumerGroup creates a new Kafka consumer group for consuming messages
func NewKafkaConsumerGroup(config *configuration.Config, logger *logrus.Logger) (sarama.ConsumerGroup, error) {
	kafkaConfig := config.Messaging.Kafka

	kafkaSaramaConfig := createKafkaSaramaConfig()

	if len(kafkaConfig.Brokers) == 0 {
		return nil, fmt.Errorf("no Kafka brokers configured")
	}

	consumerGroupName := kafkaConfig.Group.ID
	if consumerGroupName == "" {
		consumerGroupName = "gofiber-worker-group"
	}

	consumerGroup, err := sarama.NewConsumerGroup(kafkaConfig.Brokers, consumerGroupName, kafkaSaramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return consumerGroup, nil
}

// NewKafkaProducer creates a new Kafka producer for sending messages
func NewKafkaProducer(config *configuration.Config, logger *logrus.Logger) sarama.SyncProducer {
	kafkaConfig := config.Messaging.Kafka

	if len(kafkaConfig.Brokers) == 0 {
		logger.Info("no Kafka brokers configured")
		return nil
	}

	if !kafkaConfig.Producer.Enabled {
		logger.Info("Kafka producer is disabled")
		return nil
	}

	saramaConfig := createKafkaSaramaConfig()

	producer, err := sarama.NewSyncProducer(kafkaConfig.Brokers, saramaConfig)
	if err != nil {
		logger.Errorf("Failed to create Kafka producer: %v", err)
		return nil
	}

	return producer
}

// CloseKafkaProducer closes a Kafka producer
func CloseKafkaProducer(producer sarama.SyncProducer, logger *logrus.Logger) error {
	if producer != nil {
		if err := producer.Close(); err != nil {
			logger.Errorf("Error closing Kafka producer: %v", err)
			return err
		}
		logger.Info("Kafka producer closed successfully")
	}
	return nil
}

// CloseKafkaConsumerGroup closes a Kafka consumer group
func CloseKafkaConsumerGroup(consumerGroup sarama.ConsumerGroup, logger *logrus.Logger) error {
	if consumerGroup != nil {
		if err := consumerGroup.Close(); err != nil {
			logger.Errorf("Error closing Kafka consumer group: %v", err)
			return err
		}
		logger.Info("Kafka consumer group closed successfully")
	}
	return nil
}

// ConsumerGroup is an interface that wraps sarama.ConsumerGroup for easier testing
type ConsumerGroup interface {
	Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error
	Errors() <-chan error
	Close() error
	Pause(partitions map[string][]int32)
	Resume(partitions map[string][]int32)
	PauseAll()
	ResumeAll()
}
