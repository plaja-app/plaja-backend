package controllers

import (
	"encoding/json"
	"errors"
	"github.com/plaja-app/back-end/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

// enrollmentBody is the course enrollment request body structure.
type enrollmentBody struct {
	UserID   uint
	CourseID uint
}

// GetEnrollments returns the queried list of models.Enrollment.
func (c *BaseController) GetEnrollments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userID := query.Get("user_id")
	courseID := query.Get("course_id")
	statusID := query.Get("status_id")

	w.Header().Set("Content-Type", "application/json")

	var enrollments []models.Enrollment
	dbQuery := c.App.DB

	if userID != "" {
		ids := strings.Split(userID, ",")
		var intIds []int
		for _, idStr := range ids {
			intId, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}
			intIds = append(intIds, intId)
		}
		dbQuery = dbQuery.Where("user_id IN ?", intIds)
	}

	if statusID != "" {
		ids := strings.Split(statusID, ",")
		var intIds []int
		for _, idStr := range ids {
			intId, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}
			intIds = append(intIds, intId)
		}
		dbQuery = dbQuery.Where("status_id IN ?", intIds)
	}

	if courseID != "" {
		dbQuery = dbQuery.Where("course_id = ?", courseID)
	}

	if err := dbQuery.Find(&enrollments).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(enrollments)
}

// CreateEnrollment creates a new models.Enrollment.
func (c *BaseController) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var body enrollmentBody

	err := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var enrollment models.Enrollment
	enrollment = models.Enrollment{
		UserID:         body.UserID,
		CourseID:       body.CourseID,
		StatusID:       1, // enrolled
		Progress:       0,
		LastExerciseID: 1,
	}

	result := c.App.DB.Create(&enrollment)
	if result.Error != nil {
		http.Error(w, "Error creating enrollment", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
