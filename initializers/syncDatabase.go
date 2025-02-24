package initializers

import (
	"docker-test/model"
	"log"
)

func SyncDatabase() {
	if err := DB.AutoMigrate(&model.User{}, &model.Post{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
