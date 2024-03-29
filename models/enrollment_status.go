package models

import "time"

// EnrollmentStatus is the enrollment status model.
type EnrollmentStatus struct {
	ID        uint
	Title     string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
