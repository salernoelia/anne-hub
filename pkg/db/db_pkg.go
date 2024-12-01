package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

// SetupDatabase initializes the database connection
func SetupDatabase() {
	dbURL := os.Getenv("DATABASE_URL")

	var dsn string
	if dbURL != "" {
		dsn = dbURL
	} else {
		dbHost := os.Getenv("DB_HOST")         
		dbPort := os.Getenv("DB_PORT") 
  		dbUsername := os.Getenv("DB_USER")      // Ensure matching
        dbPassword := os.Getenv("DB_PASS") 
		dbName := os.Getenv("DB_NAME")
		dbSSLMode := os.Getenv("DB_SSLMODE")

		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUsername, dbPassword, dbName, dbSSLMode,
		)
	}

	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Verify the connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Database connection established successfully.")
}
