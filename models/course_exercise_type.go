package models

import "time"

// CourseExerciseType is the course_exercise_type model.
type CourseExerciseType struct {
	ID        uint
	Title     string    `gorm:"size:255"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
