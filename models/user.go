package models

import "time"

type User struct {
	ID        string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
