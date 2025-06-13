package tests

import (
	"golang_ticket_api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	connectionString := "user=postgres password= host=localhost port=5432 dbname=ticket_api_test sslmode=disable"
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	err = db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error // reset schema
	if err != nil {
		panic("failed to reset schema")
	}

	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		panic("failed to create UUID extension: " + err.Error())
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic("failed to migrate User model: " + err.Error())
	}
	if err := db.AutoMigrate(&models.TicketOption{}); err != nil {
		panic("failed to migrate TicketOption model: " + err.Error())
	}
	if err := db.AutoMigrate(&models.Ticket{}); err != nil {
		panic("failed to migrate Ticket model: " + err.Error())
	}
	if err := db.AutoMigrate(&models.Purchase{}); err != nil {
		panic("failed to migrate Purchase model: " + err.Error())
	}

	models.DB = db

	return db
}
