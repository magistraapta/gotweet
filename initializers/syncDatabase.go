package initializers

import (
	"docker-test/model"
	"log"
)

func SyncDatabase() {
	db, err := ConnectToDatabase()

	db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
