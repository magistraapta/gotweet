package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error

	databaseConfig := os.Getenv("DATABASE_CONFIG")

	DB, err = gorm.Open(postgres.Open(databaseConfig), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed Connect to Database ", err)
	}
}
