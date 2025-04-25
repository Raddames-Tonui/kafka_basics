// Package db provides functionality for interacting with the PostgreSQL database.
package db

import (
	"database/sql"          // Provides database interaction primitives
	"kafka/internal/models" // Imports the PaymentEvent struct definition from internal models
	"log"
)

// InsertPayment inserts a payment event into the "payment_events" table in PostgreSQL.
// Parameters:
// - db: a pointer to an initialized SQL database connection
// - event: a PaymentEvent struct containing the data to be inserted
// Returns:
// - error: any error encountered during the insert operation
func InsertPayment(db *sql.DB, event models.PaymentEvent) error {
	// Prepare a parameterized SQL INSERT statement to prevent SQL injection
	
	stmt, err := db.Prepare(`
		INSERT INTO payment_events (user_id, amount, status)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		log.Printf("[DB Error] Failed to prepare statement: %v", err)

		return err
	}
	// Ensure the prepared statement is closed after execution to free resources
	defer stmt.Close()

	// Execute the SQL statement with values from the PaymentEvent object.
	// Placeholders ($1, $2, $3) are bound to UserID, Amount, and Status respectively.
	_, err = stmt.Exec(event.UserID, event.Amount, event.Status)
	if err != nil {
		log.Printf("[DB Error] Exec failed: %v", err)

		// Return execution error if the insert fails
		return err
	}

	// Return nil if the operation was successful
	return nil
}
