package db

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	applyMigrations(dsn)
}


func applyMigrations(dsn string) {
    driver, err := postgres.WithInstance(DB.DB, &postgres.Config{})
    if err != nil {
        log.Fatalf("Could not create postgres driver: %v", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://db/migrations",
        "postgres", driver)
    if err != nil {
        log.Printf("could not create migration instance: %v", err);
    }

    // Try applying migrations normally
    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        fmt.Printf("Could not apply migrations, skipping: %v\n", err)
        // promptForForceMigration(m)
    } else {
        log.Println("Migrations applied successfully!")
    }
}
/*
func promptForForceMigration(m *migrate.Migrate) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Would you like to force a migration version? (y/n): ")
    answer, _ := reader.ReadString('\n')
    answer = strings.TrimSpace(answer)

    if strings.ToLower(answer) == "y" {
        fmt.Print("Enter the version number to force: ")
        versionStr, _ := reader.ReadString('\n')
        versionStr = strings.TrimSpace(versionStr)

        version, err := strconv.Atoi(versionStr)
        if err != nil {
            fmt.Printf("Invalid version number: %v\n", err)
            return
        }

        err = m.Force(version)
        if err != nil {
            log.Fatalf("Could not force migration to version %d: %v", version, err)
        }
        log.Printf("Forced migration to version %d", version)
    } else {
        log.Println("Migration not forced.")
    }
}

*/