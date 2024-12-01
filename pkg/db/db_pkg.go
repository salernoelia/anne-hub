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


func SetupDatabase() {
    dbURL := os.Getenv("DATABASE_URL")

    var dsn string
    if dbURL != "" {
        dsn = dbURL
    } else {
        dbUsername := env.GetEnvOrFatal("DB_USERNAME")
        dbPassword := env.GetEnvOrFatal("DB_PASSWORD")
        dbName := env.GetEnvOrFatal("DB_NAME")
        dbSSLMode := env.GetEnvOrFatal("DB_SSLMODE", "require")

        dsn = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbUsername, dbPassword, dbName, dbSSLMode)
    }

    var err error
    DB, err = sqlx.Connect("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }

	fmt.Println("Database connection established.")
}