package main

import (
	"anne-hub/router"
	"log"

	"anne-hub/pkg/db"

	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load() // This will load the .env file
    if err != nil {
        log.Fatal("Error loading .env file")
    }


    e := router.NewRouter()

    db.SetupDatabase()



    e.Logger.Fatal(e.Start("0.0.0.0:1323"))
}