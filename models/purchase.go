package models

import "time"

type Purchase struct {
	ID             string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Quantity       uint16
	UserId         string
	TicketOptionId string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
