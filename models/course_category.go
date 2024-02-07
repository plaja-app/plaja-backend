package models

import "time"

// CourseCategory is the course_category model.
type CourseCategory struct {
	ID        uint
	Title     string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
