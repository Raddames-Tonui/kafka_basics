// package main defines the entry point of the application.
// This package must contain the `main()` function for the Go runtime to execute it.
package main

import (
	// Importing the internal Kafka producer logic from the internal module.
	// This includes the Kafka producer setup and HTTP handler for publishing events.
	"kafka/internal/producer"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize a new Kafka producer that connects to the broker at localhost:9092
	// and targets the topic "payment_events".
	// This encapsulates Kafka writer creation using segmentio/kafka-go.
	prod := producer.NewKafkaProducer("localhost:9092", "payment_events")

	// Ensure the Kafka producer is gracefully closed when the application shuts down.
	defer prod.Close()

	// Create a new Gin HTTP router instance.
	// This will be used to register API routes and start the web server.
	router := gin.Default()

	// Register an HTTP POST endpoint at "/payment".
	// Incoming POST requests to this route will be handled by the PaymentHandler function,
	// which will parse the JSON body and send the event to Kafka.
	router.POST("/payment", producer.PaymentHandler(prod))

	router.Run(":8080")
}
