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
	LevelID          uint             `gorm:"size:50;not null;"`
	Level            CourseLevel
	StatusID         uint `gorm:"not null;"`
	Status           CourseStatus
	InstructorID     uint `gorm:"not null;"`
	Price            uint
	Instructor       User
	HasCertificate   bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
