package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/config"
	kafkaDelivery "github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/messaging/kafka"
	rabbitDelivery "github.com/arisdolanan/demo-gofiber-clean-architecture/internal/delivery/messaging/rabbitmq"
	kafkaInfra "github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/messaging/kafka"
	rabbitInfra "github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/messaging/rabbitmq"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func main() {
	viperConfig, err := configuration.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	logConfig := utils.CreateLogConfigFromViper()
	logger := config.SetupLogRus(logConfig)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupGracefulShutdown(logger, cancel)

	// Run Kafka consumer in a separate goroutine
	go RunUserConsumer(viperConfig, logger, ctx)

	// Run RabbitMQ consumer in a separate goroutine
	go RunRabbitMQConsumer(logger, ctx)

	<-ctx.Done()
}

// setupGracefulShutdown handles OS signals for graceful application shutdown
func setupGracefulShutdown(logger *logrus.Logger, cancel context.CancelFunc) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChannel
		logger.Info("Shutdown signal received, initiating graceful shutdown...")

		cancel()

		time.Sleep(2 * time.Second)
		logger.Info("Graceful shutdown completed")
		os.Exit(0)
	}()
}

// RunRabbitMQConsumer demonstrates consuming messages with flexible exchange and queue names
func RunRabbitMQConsumer(logger *logrus.Logger, ctx context.Context) {
	logger.Info("Setting up RabbitMQ consumer with flexible names")

	rabbitConsumer := rabbitInfra.NewConsumerWithLogger(logger)
	if rabbitConsumer == nil {
		logger.Error("Failed to create RabbitMQ consumer")
		return
	}
	defer rabbitConsumer.Close()

	userConsumer := rabbitDelivery.NewUserConsumer(logger)

	go func() {
		err := rabbitConsumer.ConsumeMessagesFromQueue("user.notifications", func(delivery amqp091.Delivery) error {
			return userConsumer.Consume(delivery)
		})
		if err != nil {
			logger.Errorf("Failed to consume from queue: %v", err)
		}
	}()

	go func() {
		err := rabbitConsumer.ConsumeMessagesFromExchange("user.exchange", "user.notifications.queue", "notifications.email", func(delivery amqp091.Delivery) error {
			return userConsumer.ConsumeWithRouting(delivery)
		})
		if err != nil {
			logger.Errorf("Failed to consume from exchange: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Info("RabbitMQ consumer stopped")
}

func RunUserConsumer(viperConfig *configuration.Config, logger *logrus.Logger, ctx context.Context) {
	logger.Info("Setting up user consumer")

	userConsumerGroup, err := kafkaInfra.NewKafkaConsumerGroup(viperConfig, logger)
	if err != nil {
		logger.Fatalf("Failed to create Kafka consumer group: %v", err)
	}
	defer func() {
		if err := kafkaInfra.CloseKafkaConsumerGroup(userConsumerGroup, logger); err != nil {
			logger.Errorf("Failed to close consumer group: %v", err)
		}
	}()

	userHandler := kafkaDelivery.NewUserConsumer(logger)
	kafkaInfra.ConsumeTopic(ctx, userConsumerGroup, "users", logger, userHandler.Consume)
}
