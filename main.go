package main

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Purchase struct {
	ID             string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Quantity       uint16
	UserId         string
	TicketOptionId string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type TicketOption struct {
	ID         string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name       string
	Desc       string
	Allocation uint16
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Ticket struct {
	ID             string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	TicketOptionId string
	PurchaseId     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func main() {
	connection_string := "user=postgres password= host=127.0.0.1 port=5432 dbname=ticket_api_development sslmode=disable"
	db, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
}
