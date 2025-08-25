package kafka

// import (
// 	"context"
// 	"fmt"

// 	"github.com/segmentio/kafka-go"
// )

// // EnsureTopics ensures that the specified Kafka topics exist.
// func EnsureTopics(brokers []string, topics []kafka.TopicConfig) error {
// 	conn, err := kafka.DialContext(context.Background(), "tcp", brokers[0])
// 	if err != nil {
// 		return fmt.Errorf("failed to dial Kafka broker %s: %w", brokers[0], err)
// 	}
// 	defer conn.Close()

// 	controller, err := conn.Controller()
// 	if err != nil {
// 		return fmt.Errorf("failed to get Kafka controller: %w", err)
// 	}

// 	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
// 	if err != nil {
// 		return fmt.Errorf("failed to dial Kafka controller %s:%d: %w", controller.Host, controller.Port, err)
// 	}
// 	defer controllerConn.Close()

// 	// No context and no variadic error — pass topic configs directly
// 	err = controllerConn.CreateTopics(topics...)
// 	if err != nil {
// 		if err.Error() == "Topic with this name already exists" {
// 			appLogger.Info("Kafka topic already exists, skipping creation.")
// 		} else {
// 			return fmt.Errorf("failed to create Kafka topics: %w", err)
// 		}
// 	}

// 	appLogger.Info("Kafka topics ensured successfully.")
// 	return nil
// }
