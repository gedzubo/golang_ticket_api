package models

import "time"

type TicketOption struct {
	ID         string    `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	Allocation uint16    `json:"allocation"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
