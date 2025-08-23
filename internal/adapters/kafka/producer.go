// ANNOTATION: DO NOT MODIFY - Core infrastructure logic shared across all services.
package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	log    *zap.Logger
}

// ANNOTATION: Kafka producer setup. Reuse as-is for message publishing.
func NewProducer(brokers []string, logger *zap.Logger) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
		log: logger,
	}
}

// ANNOTATION: Publish method to send messages to Kafka topics. DO NOT CHANGE
func (p *Producer) Publish(ctx context.Context, topic, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	}
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		p.log.Error("failed to publish kafka message", zap.Error(err))
		return err
	}
	return nil
}

// ANNOTATION: ProduceJSON method to send JSON data to Kafka topics. DO NOT CHANGE
func (p *Producer) ProduceJSON(topic string, data interface{}) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Value: value,
	})
}

// ANNOTATION: Close method to clean up resources. DO NOT CHANGE
func (p *Producer) Close() error {
	return p.writer.Close()
}
