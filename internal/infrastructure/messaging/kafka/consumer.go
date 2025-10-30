package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

// ConsumerHandler defines the function signature for handling Kafka messages
type ConsumerHandler func(message *sarama.ConsumerMessage) error

// ConsumerGroupHandler implements the sarama.ConsumerGroupHandler interface
type ConsumerGroupHandler struct {
	Handler ConsumerHandler
	Log     *logrus.Logger
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if message == nil {
			continue
		}

		err := h.Handler(message)
		if err != nil {
			h.Log.WithError(err).Error("Failed to process message")
			session.MarkMessage(message, "")
			continue
		}

		session.MarkMessage(message, "")
	}

	return nil
}

// ConsumeTopic starts consuming messages from a Kafka topic
func ConsumeTopic(ctx context.Context, consumerGroup sarama.ConsumerGroup, topic string, log *logrus.Logger, handler ConsumerHandler) {
	consumerHandler := &ConsumerGroupHandler{
		Handler: handler,
		Log:     log,
	}

	// Start error handling goroutine
	go func() {
		for err := range consumerGroup.Errors() {
			log.WithError(err).Error("Consumer group error")
		}
	}()

	// Consumer loop
	for {
		// Check if context is cancelled before consuming
		select {
		case <-ctx.Done():
			log.Info("Context cancelled, stopping consumer")
			if err := consumerGroup.Close(); err != nil {
				log.WithError(err).Error("Error closing consumer group")
			}
			return
		default:
		}

		if err := consumerGroup.Consume(ctx, []string{topic}, consumerHandler); err != nil {
			log.WithError(err).Error("Error from consumer")
			if err == sarama.ErrClosedConsumerGroup {
				log.Info("Consumer group closed")
				return
			}
			continue
		}

		// Check if context is cancelled after consuming
		select {
		case <-ctx.Done():
			log.Info("Context cancelled, stopping consumer")
			if err := consumerGroup.Close(); err != nil {
				log.WithError(err).Error("Error closing consumer group")
			}
			return
		default:
		}
	}
}
