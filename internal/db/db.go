// Package db provides database initialization logic for the application.
package db

import (
	"database/sql" // Standard library for SQL database interactions
	"fmt"          // For formatting connection string
	"log"          // For logging fatal errors
	"os"           // For accessing environment variables
)

// InitDB initializes a connection to a PostgreSQL database and returns the connection pool.
// It reads connection credentials from environment variables with fallback defaults.
func InitDB() *sql.DB {
	// Format the PostgreSQL connection string using environment variables or fallbacks.
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("POSTGRES_USER", "postgres"),        // Username
		getEnv("POSTGRES_PASSWORD", "password"),    // Password
		getEnv("POSTGRES_HOST", "localhost"),       // Host (default: localhost)
		getEnv("POSTGRES_PORT", "5432"),            // Port (default: 5432)
		getEnv("POSTGRES_DB", "payments_events"),   // Database name
	)

	// Open a new database connection using the formatted connection string.
	// The driver "postgres" should be registered via a PostgreSQL driver like lib/pq.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// If connection initialization fails, terminate the application with a fatal log.
		log.Fatal("Error connecting to the database:", err)
	}

	// Return the database connection object.
	return db
}

// getEnv returns the value of the environment variable `key` if it exists.
// If the variable is not set, it returns the provided `fallback` value instead.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
