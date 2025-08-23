// ANNOTATION: DO NOT MODIFY - Core infrastructure logic shared across all services.
package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// StartExampleConsumer initializes a Kafka consumer that listens for messages on the specified topic and processes them using the provided handler function. Modify the topic and groupID as needed for your application.
func StartExampleConsumer(brokers []string, topic string, groupID string, logger *zap.Logger, handlerFunc func(data map[string]string)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	go func() {
		defer r.Close()
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				logger.Error("consumer read error", zap.Error(err))
				continue
			}
			logger.Info("received kafka message", zap.ByteString("value", m.Value))

			var payload map[string]string
			if err := json.Unmarshal(m.Value, &payload); err != nil {
				logger.Error("failed to unmarshal kafka message", zap.Error(err))
				continue
			}

			handlerFunc(payload)
		}
	}()
}
