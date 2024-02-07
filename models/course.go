package models

import "time"

// Course is the course model.
type Course struct {
	ID             uint
	InstructorID   uint
	Instructor     Account
	Title          string           `gorm:"size:255"`
	Descriptions   string           `gorm:"size:65000"`
	Categories     []CourseCategory `gorm:"many2many:course_categories_junction;"`
	StatusID       uint             `gorm:"not null"`
	Status         CourseStatus
	HasCertificate bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
