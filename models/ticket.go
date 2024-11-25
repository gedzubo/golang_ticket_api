package models

import "time"

type Ticket struct {
	ID             string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	TicketOptionId string
	PurchaseId     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
