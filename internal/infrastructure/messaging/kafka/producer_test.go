package kafka

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// MockKafkaMessage implements KafkaMessage interface for testing
type MockKafkaMessage struct {
	ID string
}

func (m MockKafkaMessage) GetId() string {
	return m.ID
}

func TestNewProducer(t *testing.T) {
	// Create a mock logger
	logger := logrus.New()

	// Create a mock sarama.SyncProducer
	mockProducer := &MockSyncProducer{}

	// Test creating KafkaProducer wrapper
	kafkaProducer := &KafkaProducer[MockKafkaMessage]{
		producer: mockProducer,
		logger:   logger,
		Topic:    "test-topic",
	}
	assert.NotNil(t, kafkaProducer)
	assert.Equal(t, mockProducer, kafkaProducer.producer)
	assert.Equal(t, logger, kafkaProducer.logger)
}

// MockSyncProducer implements sarama.SyncProducer for testing
type MockSyncProducer struct {
	sarama.SyncProducer
	CloseFunc func() error
}

func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return 0, 0, nil
}

func (m *MockSyncProducer) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func TestKafkaProducer_Close(t *testing.T) {
	logger := logrus.New()

	t.Run("Close nil producer", func(t *testing.T) {
		kafkaProducer := &KafkaProducer[MockKafkaMessage]{
			producer: nil,
			logger:   logger,
		}
		err := kafkaProducer.Close()
		assert.NoError(t, err)
	})

	t.Run("Close valid producer", func(t *testing.T) {
		mockProducer := &MockSyncProducer{}
		kafkaProducer := &KafkaProducer[MockKafkaMessage]{
			producer: mockProducer,
			logger:   logger,
		}
		err := kafkaProducer.Close()
		assert.NoError(t, err)
	})
}

func TestCloseKafkaProducer(t *testing.T) {
	logger := logrus.New()

	t.Run("Close nil producer", func(t *testing.T) {
		err := CloseKafkaProducer(nil, logger)
		assert.NoError(t, err)
	})

	t.Run("Close valid producer", func(t *testing.T) {
		mockProducer := &MockSyncProducer{}
		err := CloseKafkaProducer(mockProducer, logger)
		assert.NoError(t, err)
	})

	t.Run("Close producer with error", func(t *testing.T) {
		mockProducer := &MockSyncProducer{
			CloseFunc: func() error {
				return assert.AnError
			},
		}
		err := CloseKafkaProducer(mockProducer, logger)
		assert.Error(t, err)
	})
}

func TestCloseKafkaConsumerGroup(t *testing.T) {
	logger := logrus.New()

	t.Run("Close nil consumer group", func(t *testing.T) {
		err := CloseKafkaConsumerGroup(nil, logger)
		assert.NoError(t, err)
	})

	t.Run("Close valid consumer group", func(t *testing.T) {
		// Create a mock consumer group using the one from consumer_test.go
		mockConsumerGroup := &MockConsumerGroup{}
		err := CloseKafkaConsumerGroup(mockConsumerGroup, logger)
		assert.NoError(t, err)
	})

	t.Run("Close consumer group with error", func(t *testing.T) {
		// Create a mock consumer group using the one from consumer_test.go
		mockConsumerGroup := &MockConsumerGroup{
			CloseFunc: func() error {
				return assert.AnError
			},
		}
		err := CloseKafkaConsumerGroup(mockConsumerGroup, logger)
		assert.Error(t, err)
	})
}
