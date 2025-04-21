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

	err = db.AutoMigrate(
		&models.Purchase{},
		&models.TicketOption{},
		&models.Ticket{},
		&models.User{},
	)
	if err != nil {
		panic("failed to migrate test schema")
	}

	models.DB = db

	return db
}
