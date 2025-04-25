package consumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"kafka/internal/db"
	"kafka/internal/models"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "github.com/lib/pq"


	"github.com/segmentio/kafka-go"
)

// StartConsumer listens to a Kafka topic and stores received payment events into PostgreSQL.
func StartConsumer(broker, topic, groupID string, dbConn *sql.DB) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("[Shutdown] Termination signal received. Shutting down consumer...")
		cancel()
	}()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1e3,  // 1KB
		MaxBytes: 10e6, // 10MB
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("[Kafka Error] Failed to close reader: %v", err)
		}
	}()

	log.Printf("[Kafka] Consumer started | Broker: %s | Topic: %s | GroupID: %s", broker, topic, groupID)

	for {
		select {
		case <-ctx.Done():
			log.Println("[Kafka] Context cancelled. Exiting consumer loop.")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("[Kafka Error] Failed to read message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			var event models.PaymentEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("[JSON Error] Invalid message format: %v | Raw: %s", err, string(msg.Value))
				continue
			}

			log.Printf("[Received] Payment Event Fetched: %+v", event)

			if err := db.InsertPayment(dbConn, event); err != nil {
				log.Printf("[DB Error] Could not insert for UserID %s: %v", event.UserID, err)
				continue
			}

			log.Printf("[Success] Stored payment event for UserID: %s", event.UserID)
		}
	}
}
