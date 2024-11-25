package models

import "time"

type TicketOption struct {
	ID         string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name       string
	Desc       string
	Allocation uint16
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
