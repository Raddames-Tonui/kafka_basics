// 🔹 handler.go (Gin HTTP handler)
// API / Delivery Layer — Handles HTTP requests, binds JSON, and triggers Kafka publishing logic
//
// Responsibility:
// Acts as the HTTP interface layer for incoming API requests.
//
// Key Features:
// - Exposes a PaymentHandler that:
//   - Receives HTTP POST requests
//   - Parses incoming JSON into a PaymentEvent
//   - Calls KafkaProducer to publish to Kafka
//   - Handles HTTP response codes and error messages
//
// Purpose:
// This file belongs to the API layer. It connects the web request to backend logic.
// It depends on KafkaProducer but doesn't care how Kafka works—just that it can publish an event.

package producer

import (
	"context"
	"encoding/json"
	"kafka/internal/models"          // PaymentEvent struct is defined here
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"       // Gin framework for HTTP routing
	"github.com/segmentio/kafka-go"  // Kafka client for Go
)

// PaymentHandler returns a Gin handler that processes payment events and publishes them to Kafka.
// It includes structured logging, timeout context, and clean error handling.
func PaymentHandler(producer *KafkaProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.PaymentEvent // Struct to hold incoming JSON data

		// Attempt to bind incoming JSON to the PaymentEvent struct.
		// If the client sends invalid or malformed JSON, respond with a 400 Bad Request.
		if err := c.ShouldBindJSON(&event); err != nil {
			log.Printf("[ERROR] Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		// Marshal the PaymentEvent struct to JSON bytes for sending to Kafka.
		// This transforms the structured data into a format Kafka can store.
		value, err := json.Marshal(event)
		if err != nil {
			log.Printf("[ERROR] Failed to marshal event: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal event"})
			return
		}

		// Create a timeout context to ensure that Kafka write operation does not hang indefinitely.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Publish the message to Kafka.
		// The message key is the UserID — Kafka will use this to route events to the same partition.
		// This helps maintain ordering of events for a specific user.
		err = producer.writer.WriteMessages(ctx, kafka.Message{
			Key:   []byte(event.UserID),
			Value: value,
		})
		if err != nil {
			log.Printf("[ERROR] Failed to publish event to Kafka: %v", err)
			// Optional: you could also record this failure using Prometheus metrics
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send to Kafka"})
			return
		}

		// Log success to console or logging system
		log.Printf("[INFO] Event published to Kafka: userID=%s", event.UserID)
		// Optional: increment Prometheus success counter here

		// Respond to the client that the event was published successfully
		c.JSON(http.StatusOK, gin.H{"status": "event published"})
	}
}
