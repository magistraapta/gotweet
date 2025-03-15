package main

import (
	"docker-test/initializers"
	"docker-test/internal/api"
	"log"
)

func init() {
	initializers.LoadEnvFile()
	initializers.ConnectToDatabase()
	initializers.SyncDatabase()
}

func main() {
	//connect to db
	db, err := initializers.ConnectToDatabase()

	if err != nil {
		log.Fatal("Failed to connect database")
	}

	router := api.ApiRouter(db)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
}
