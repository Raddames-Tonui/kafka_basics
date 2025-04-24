
// 🔹 payment.go (Kafka Producer logic)
//  Infrastructure  Manages Kafka producer setup, message publishing, and closing connection
// Responsibility: Handles the lower-level Kafka publishing logic.

// Key Features:
// Defines a KafkaProducer struct that wraps the Kafka writer.

// Encapsulates producer setup (NewKafkaProducer).

// Provides a Publish method to send messages to Kafka.

// Manages Kafka connection lifecycle (Close method).

// Purpose:
// This file acts as a Kafka utility or infrastructure layer. It abstracts Kafka details so other parts of your app don’t need to deal with it directly.


package producer

import (
	"context"
	"github.com/segmentio/kafka-go" // Kafka client library for Go
	"log"
)

// KafkaProducer wraps the kafka.Writer for producing messages to Kafka topics.
type KafkaProducer struct {
	writer *kafka.Writer // Kafka writer instance used to send messages
}

// NewKafkaProducer initializes and returns a KafkaProducer instance.
// - brokerURL: address of the Kafka broker (e.g., "localhost:9092")
// - topic: the Kafka topic to which messages will be published
func NewKafkaProducer(brokerURL, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerURL},       // List of Kafka brokers
		Topic:    topic,                     // Target topic for messages
		Balancer: &kafka.LeastBytes{},       // Load balancing strategy: sends messages to the partition with the fewest bytes
	})
	return &KafkaProducer{writer: writer}
}

// Publish sends a single message to the Kafka topic.
// - key: used for partitioning; same key always goes to the same partition
// - value: the actual message payload
func (p *KafkaProducer) Publish(key, value []byte) error {
	msg := kafka.Message{
		Key:   key,   // Useful for keyed partitioning and ordering
		Value: value, // Message body
	}
	// Send the message with context for cancellation/timeouts
	return p.writer.WriteMessages(context.Background(), msg)
}

// Close gracefully shuts down the Kafka producer to release network resources.
func (p *KafkaProducer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Printf("Error closing Kafka producer: %v", err)
	}
}
