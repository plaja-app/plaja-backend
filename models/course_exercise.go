package models

import "time"

// CourseExercise is the course exercise model.
type CourseExercise struct {
	ID        uint
	Title     string `gorm:"size:255"`
	Content   string `gorm:"size:65000"`
	Length    uint
	CourseID  uint   `gorm:"not null"`
	Course    Course `json:"-"`
	TypeID    uint   `gorm:"not null"`
	Type      CourseExerciseType
	CreatedAt time.Time
	UpdatedAt time.Time
}
