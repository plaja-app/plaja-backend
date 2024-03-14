package models

import "time"

// Course is the course model.
type Course struct {
	ID               uint
	Thumbnail        string
	Title            string           `gorm:"size:255"`
	ShortDescription string           `gorm:"size:255"`
	Description      string           `gorm:"size:65000"`
	Categories       []CourseCategory `gorm:"many2many:course_categories_junction;"`
	LevelID          uint             `gorm:"not null;"`
	Level            CourseLevel
	StatusID         uint         `gorm:"not null;"`
	Status           CourseStatus `json:"-"`
	InstructorID     uint         `gorm:"not null;"`
	Instructor       User
	Exercises        []CourseExercise `gorm:"foreignkey:CourseID"`
	Length           uint
	Price            uint
	HasCertificate   bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
