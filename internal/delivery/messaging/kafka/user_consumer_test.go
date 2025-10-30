package kafka

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestUserConsumer_Consume(t *testing.T) {
	// Create a logger for testing
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	// Create a user consumer
	consumer := NewUserConsumer(logger)

	// Create a sample user
	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	// Marshal the user to JSON
	userBytes, err := json.Marshal(user)
	assert.NoError(t, err)

	// Create a mock Kafka message
	message := &sarama.ConsumerMessage{
		Key:       []byte("test-key"),
		Value:     userBytes,
		Topic:     "users",
		Partition: 0,
		Offset:    0,
		Timestamp: time.Now(),
	}

	// Test consuming the message
	err = consumer.Consume(message)
	assert.NoError(t, err)
}
