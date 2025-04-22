package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	_ "github.com/lib/pq"
)

func main() {
	// Kafka writer setup ( Producer Setup) 
	// Creates a Kafka writer that connects to Kafka and sends messages to the payment_events topic.
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "payment_events",
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	// PostgreSQL setup
	// Opens a connection to the PostgreSQL database (payments_events).
	db, err := sql.Open("postgres", "postgres://postgres:123@localhost:5432/payments_events?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS payments (
			id SERIAL PRIMARY KEY,
			user_id TEXT NOT NULL,
			amount INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	// Gin Framework (API Setup)
	r := gin.Default()

	// POST /payment
	r.POST("/payment", func(c *gin.Context) {
		//  Bind the incoming JSON payload to the payload struct.
		var payload struct {
			UserID string `json:"user_id"`
			Amount int    `json:"amount"`
		}

		if err := c.BindJSON(&payload); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// Create message
		message := fmt.Sprintf(`{"user_id":"%s","amount":%d}`, payload.UserID, payload.Amount)

		// Send to Kafka
		// Once the data is received, the backend sends it as a message to Kafka using the Kafka producer.
		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(payload.UserID),
			Value: []byte(message),
		})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to send Kafka message"})
			return
		}

		// Save to database
		// The same payment data is inserted into the PostgreSQL database.
		_, err = db.Exec(`INSERT INTO payments(user_id, amount, created_at) VALUES($1, $2, $3)`,
			payload.UserID, payload.Amount, time.Now())
		if err != nil {
			c.JSON(500, gin.H{"error": "Database insert failed"})
			return
		}

		c.JSON(200, gin.H{"message": "Payment processed"})
	})

	r.Run(":8081")
 // Server running on port 8080
}
