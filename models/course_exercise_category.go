package models

import "time"

// CourseExerciseCategory is the course exercise_category model.
type CourseExerciseCategory struct {
	ID        uint
	Title     string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
