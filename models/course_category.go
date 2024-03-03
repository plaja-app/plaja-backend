package models

import "time"

// CourseCategory is the course category model.
type CourseCategory struct {
	ID        uint
	Title     string    `gorm:"size:255"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
