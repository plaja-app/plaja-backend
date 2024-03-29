package models

import "time"

// CourseStatus is the course status model.
type CourseStatus struct {
	ID        uint
	Title     string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
