package entity

import (
	"encoding/json"
	"strconv"
	"time"
)

// LoginEvent represents a user login event
type LoginEvent struct {
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

// GetId returns the unique identifier for the event (required by KafkaMessage interface)
func (le *LoginEvent) GetId() string {
	return strconv.FormatInt(le.UserID, 10)
}

// MarshalBinary implements encoding.BinaryMarshaler for Kafka message serialization
func (le *LoginEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(le)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for Kafka message deserialization
func (le *LoginEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, le)
}
