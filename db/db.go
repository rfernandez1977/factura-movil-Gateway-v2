package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// InitDB initializes the PostgreSQL database connection
func InitDB() (*sql.DB, error) {
	connStr := "host=localhost user=postgres password=secret dbname=gateway sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Create documents table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS documents (
		id SERIAL PRIMARY KEY,
		type VARCHAR(50) NOT NULL,
		data TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`)
	if err != nil {
		log.Fatal("Failed to create documents table:", err)
	}

	return db, nil
}