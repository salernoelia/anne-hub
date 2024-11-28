package main

import (
	"anne-hub/router"
	"log"

	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load() // This will load the .env file
    if err != nil {
        log.Fatal("Error loading .env file")
    }


    // db.SetDatabase()
    e := router.NewRouter()


    e.Logger.Fatal(e.Start("0.0.0.0:1323"))
}