package kafka

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewUserProducer(t *testing.T) {
	// Create a mock logger
	logger := logrus.New()

	// Create a mock producer
	var producer sarama.SyncProducer

	// Create a user producer
	userProducer := NewUserProducer(producer, logger)

	// Test that the user producer is created correctly
	assert.NotNil(t, userProducer)
	assert.Equal(t, logger, userProducer.logger)
	assert.Equal(t, producer, userProducer.producer)
	assert.Equal(t, "users", userProducer.topic)
}

// MockSyncProducer implements sarama.SyncProducer for testing
type MockSyncProducer struct {
	sarama.SyncProducer
	CloseFunc       func() error
	SendMessageFunc func(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(msg)
	}
	return 0, 0, nil
}

func (m *MockSyncProducer) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func TestUserProducer_PublishUserEvent(t *testing.T) {
	logger := logrus.New()

	t.Run("PublishUserEvent with nil producer", func(t *testing.T) {
		userProducer := &UserProducer{
			producer: nil,
			logger:   logger,
			topic:    "users",
		}

		user := &entity.User{
			ID:        1,
			Email:     "test@example.com",
			UpdatedAt: time.Now(),
		}

		err := userProducer.PublishUserEvent(user)
		assert.NoError(t, err)
	})

	t.Run("PublishUserEvent with valid producer", func(t *testing.T) {
		mockProducer := &MockSyncProducer{
			SendMessageFunc: func(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
				return 0, 0, nil
			},
		}

		userProducer := &UserProducer{
			producer: mockProducer,
			logger:   logger,
			topic:    "users",
		}

		user := &entity.User{
			ID:        1,
			Email:     "test@example.com",
			UpdatedAt: time.Now(),
		}

		err := userProducer.PublishUserEvent(user)
		assert.NoError(t, err)
	})

	t.Run("PublishUserEvent with error", func(t *testing.T) {
		mockProducer := &MockSyncProducer{
			SendMessageFunc: func(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
				return 0, 0, assert.AnError
			},
		}

		userProducer := &UserProducer{
			producer: mockProducer,
			logger:   logger,
			topic:    "users",
		}

		user := &entity.User{
			ID:        1,
			Email:     "test@example.com",
			UpdatedAt: time.Now(),
		}

		err := userProducer.PublishUserEvent(user)
		assert.Error(t, err)
	})
}

func TestUserProducer_Close(t *testing.T) {
	logger := logrus.New()

	t.Run("Close with nil producer", func(t *testing.T) {
		userProducer := &UserProducer{
			producer: nil,
			logger:   logger,
			topic:    "users",
		}
		err := userProducer.Close()
		assert.NoError(t, err)
	})

	t.Run("Close with valid producer", func(t *testing.T) {
		mockProducer := &MockSyncProducer{}
		userProducer := &UserProducer{
			producer: mockProducer,
			logger:   logger,
			topic:    "users",
		}
		err := userProducer.Close()
		assert.NoError(t, err)
	})

	t.Run("Close with error", func(t *testing.T) {
		mockProducer := &MockSyncProducer{
			CloseFunc: func() error {
				return assert.AnError
			},
		}
		userProducer := &UserProducer{
			producer: mockProducer,
			logger:   logger,
			topic:    "users",
		}
		err := userProducer.Close()
		assert.Error(t, err)
	})
}
