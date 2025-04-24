
// 🔹 handler.go (Gin HTTP handler)
// handler.go	API / Delivery Layer	Handles HTTP requests, binds JSON, and triggers Kafka publishing logic
// Responsibility: Acts as the HTTP interface layer for incoming API requests.

// Key Features:
// Exposes a PaymentHandler that:

// Receives HTTP POST requests.

// Parses incoming JSON into a PaymentEvent.

// Calls KafkaProducer.Publish to publish to Kafka.

// Handles HTTP response codes and error messages.

// Purpose:
// This file belongs to the API layer. It connects the web request to backend logic. It depends on KafkaProducer but doesn't care how Kafka works—just that it can publish an event.

package producer

import (
	"context"
	"encoding/json"
	"kafka/internal/models"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// PaymentHandler returns a Gin handler that processes payment events and publishes them to Kafka.
// It includes structured logging, timeout context, and clean error handling.
func PaymentHandler(producer *KafkaProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.PaymentEvent // Struct to hold incoming JSON data

		// Attempt to bind incoming JSON to the PaymentEvent struct
		if err := c.ShouldBindJSON(&event); err != nil {
			log.Printf("[ERROR] Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		// Marshal the PaymentEvent struct to JSON bytes for Kafka
		value, err := json.Marshal(event)
		if err != nil {
			log.Printf("[ERROR] Failed to marshal event: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal event"})
			return
		}

		// Define a timeout context to avoid blocking indefinitely when publishing
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt to publish the event to Kafka using the UserID as the key (ensures ordering per user)
		err = producer.writer.WriteMessages(ctx, kafka.Message{
			Key:   []byte(event.UserID),
			Value: value,
		})
		if err != nil {
			log.Printf("[ERROR] Failed to publish event to Kafka: %v", err)
			// Optional: increment Prometheus failure counter here
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send to Kafka"})
			return
		}

		log.Printf("[INFO] Event published to Kafka: userID=%s", event.UserID)
		// Optional: increment Prometheus success counter here

		// Respond with a success message
		c.JSON(http.StatusOK, gin.H{"status": "event published"})
	}
}
