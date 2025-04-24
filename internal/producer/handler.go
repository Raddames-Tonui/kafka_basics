package producer

import (
	"encoding/json"
	"kafka/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PaymentHandler(producer *KafkaProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.PaymentEvent

		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		value, err := json.Marshal(event)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal event"})
			return
		}

		err = producer.Publish([]byte(event.UserID), value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send to Kafka"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "event published"})
	}
}
