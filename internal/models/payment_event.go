package models

// PaymentEvent defines the structure of a payment message.
type PaymentEvent struct {
	UserID string `json:"user_id"`
	Amount int    `json:"amount"`
	Status *string `json:"status,omitempty"` // The status field is optional and can be nil
}

