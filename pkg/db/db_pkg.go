package db

import (
	"fmt"
	"log"
	"os"

	"anne-hub/pkg/env"

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
		dbHost := env.GetEnvOrFatal("DB_HOST")         
		dbPort := env.GetEnvOrFatal("DB_PORT", "5432") 
		dbUsername := env.GetEnvOrFatal("DB_USERNAME")
		dbPassword := env.GetEnvOrFatal("DB_PASSWORD")
		dbName := env.GetEnvOrFatal("DB_NAME")
		dbSSLMode := env.GetEnvOrFatal("DB_SSLMODE", "require")

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
