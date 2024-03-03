package models

import "time"

// UserType is the user type model.
type UserType struct {
	ID        uint
	Title     string    `gorm:"size:255"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
