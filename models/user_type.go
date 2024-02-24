package models

import "time"

// UserType is the user_type model.
type UserType struct {
	ID        uint
	Title     string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}