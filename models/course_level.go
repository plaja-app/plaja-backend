package models

import "time"

// CourseLevel is the course level model.
type CourseLevel struct {
	ID        uint
	Title     string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
