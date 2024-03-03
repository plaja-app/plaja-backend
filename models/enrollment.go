package models

import "time"

// Enrollment is the enrollment model.
type Enrollment struct {
	UserID         uint   `gorm:"primaryKey;autoIncrement:false;not null"`
	User           User   `json:"-"`
	CourseID       uint   `gorm:"primaryKey;autoIncrement:false;not null"`
	Course         Course `json:"-"`
	Progress       uint
	StatusID       uint             `gorm:"not null"`
	Status         EnrollmentStatus `json:"-"`
	LastExerciseID uint             `gorm:"not null"`
	LastExercise   CourseExercise   `json:"-"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
