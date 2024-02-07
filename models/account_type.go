package models

import "time"

// AccountType is the account_type model.
type AccountType struct {
	ID        uint
	Title     string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
