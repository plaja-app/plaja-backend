package models

import "time"

// CourseExercise is the course exercise model.
type CourseExercise struct {
	ID         uint
	CourseID   uint `gorm:"not null"`
	Course     Course
	Title      string `gorm:"size:255"`
	CategoryID uint   `gorm:"not null"`
	Category   CourseExerciseCategory
	Content    string `gorm:"size:65000"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
