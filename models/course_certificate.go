package models

import "time"

// CourseCertificate is the course certificate model.
type CourseCertificate struct {
	ID        uint
	UserID    uint   `gorm:"not null"`
	User      User   `json:"-"`
	CourseID  uint   `gorm:"not null"`
	Course    Course `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
