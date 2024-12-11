package main

import (
	"anne-hub/router"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"anne-hub/pkg/db"

	"github.com/joho/godotenv"
)

func main() {
   err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found. Proceeding with environment variables.")
    }




    e := router.NewRouter()

    e.Static("/files", "./static")

    db.SetupDatabase()

    // In main.go
    go func() {
        if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
            e.Logger.Fatal("Shutting down the server")
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the server with
    // a timeout of 10 seconds.
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := e.Shutdown(ctx); err != nil {
        e.Logger.Fatal(err)
    }
}