package models

import "time"

// User is the user model.
type User struct {
	ID         uint `gorm:"type:int;"`
	ProfilePic string
	FirstName  string `gorm:"size:255;"`
	LastName   string `gorm:"size:255;"`
	Email      string `gorm:"size:255;unique;"`
	Password   string `gorm:"size:255" json:"-"`
	UserTypeID uint   `gorm:"not null"`
	UserType   UserType
	CreatedAt  time.Time
	UpdatedAt  time.Time `json:"-"`
}
