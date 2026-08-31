package main

import (
	"backendTemp/internal/database"
	"backendTemp/internal/di"

	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongodb, err := database.NewConnection()
	if err != nil {
		log.Fatal("Mongodb Connection can't be made: {}", err)
	}
	router, _ := di.InitializeApp(mongodb.Database)

	router.Run()
}
