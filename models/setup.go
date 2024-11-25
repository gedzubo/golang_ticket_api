package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	connection_string := "user=postgres password= host=127.0.0.1 port=5432 dbname=ticket_api_development sslmode=disable"
	db, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	DB = db
}
