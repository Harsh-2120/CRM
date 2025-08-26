package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// EnsureTopics ensures that the specified Kafka topics exist.
// It connects to the Kafka controller and creates topics if missing.
func EnsureTopics(brokers []string, topics []kafka.TopicConfig) error {
	// Connect to the first broker
	conn, err := kafka.DialContext(context.Background(), "tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("failed to dial Kafka broker %s: %w", brokers[0], err)
	}
	defer conn.Close()

	// Get controller (the broker responsible for topic management)
	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get Kafka controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to dial Kafka controller %s:%d: %w", controller.Host, controller.Port, err)
	}
	defer controllerConn.Close()

	// Create topics if not already present
	err = controllerConn.CreateTopics(topics...)
	if err != nil && err.Error() != "Topic with this name already exists" {
		return fmt.Errorf("failed to create Kafka topics: %w", err)
	}

	return nil
}
