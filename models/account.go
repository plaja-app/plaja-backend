package models

import "time"

// Account is the account (i.e. user) model.
type Account struct {
	ID            uint   `gorm:"type:int;"`
	FirstName     string `gorm:"size:255;"`
	LastName      string `gorm:"size:255"`
	UserName      string `gorm:"size:255"`
	Email         string `gorm:"size:255"`
	Password      string `gorm:"size:255"`
	AccountTypeID uint   `gorm:"not null"`
	AccountType   AccountType
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
