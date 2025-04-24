package producer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer initializes and returns a new Kafka producer.
func NewKafkaProducer(brokerURL, topic string) *KafkaProducer{
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerURL},
		Topic: topic,
		Balancer: &kafka.LeastBytes{},
	})
	return &KafkaProducer{writer: writer}
}

// Publish sends the message to Kafka.
func (p *KafkaProducer) Publish(key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}
	return p.writer.WriteMessages(context.Background(), msg)
}

// Close the Kafka producer connection.
func (p *KafkaProducer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Printf("Error closing Kafka producer: %v", err)
	}
}
