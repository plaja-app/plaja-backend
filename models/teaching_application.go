package models

import "time"

// TeachingApplication is the teaching application model.
type TeachingApplication struct {
	UserID         uint   `gorm:"primaryKey;autoIncrement:false;not null"`
	User           User   `json:"-"`
	Experience     string `gorm:"size:255;"`
	Motivation     string `gorm:"size:255;"`
	PlatformChoice string `gorm:"size:255;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
