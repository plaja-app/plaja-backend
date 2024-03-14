package controllers

import (
	"encoding/json"
	"github.com/plaja-app/back-end/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CourseExerciseInput is the exercises input request body structure.
type CourseExerciseInput struct {
	Exercises    []ExerciseInput
	InstructorID uint
	CourseID     uint
}

// ExerciseInput is the exercise input structure.
type ExerciseInput struct {
	Title   string
	Content string
}

// CreateCourseExercises creates a new models.CourseExercise.
func (c *BaseController) CreateCourseExercises(w http.ResponseWriter, r *http.Request) {
	var body CourseExerciseInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the instructor's permission to add exercises to the course
	var course models.Course
	if err := c.App.DB.First(&course, body.CourseID).Error; err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	if course.InstructorID != body.InstructorID {
		http.Error(w, "Instructor not authorized to add exercises to this course", http.StatusUnauthorized)
		return
	}

	// Create and save new exercises
	for _, ex := range body.Exercises {
		newExercise := models.CourseExercise{
			CourseID:  body.CourseID,
			TypeID:    1,
			Length:    calculateExerciseLength(ex.Content),
			Title:     ex.Title,
			Content:   ex.Content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := c.App.DB.Create(&newExercise).Error; err != nil {
			http.Error(w, "Failed to add exercise", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// GetCourseExercises returns the queried list of models.CourseExercise.
func (c *BaseController) GetCourseExercises(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	exerciseID := query.Get("exercise_id")
	courseID := query.Get("course_id")

	w.Header().Set("Content-Type", "application/json")

	// Convert courseID to int
	courseIDInt, err := strconv.Atoi(courseID)
	if err != nil {
		http.Error(w, "Invalid course ID format", http.StatusBadRequest)
		return
	}

	var data []models.CourseExercise

	if exerciseID == "all" {
		c.App.DB.Where("course_id = ?", courseIDInt).Find(&data)
	} else {
		ids := strings.Split(exerciseID, ",")
		var intIDs []int
		for _, idStr := range ids {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid exercise ID format", http.StatusBadRequest)
				return
			}
			intIDs = append(intIDs, id)
		}
		c.App.DB.Where("id IN ? AND course_id = ?", intIDs, courseIDInt).Find(&data)
	}

	if len(data) == 0 {
		http.NotFound(w, r)
	} else {
		json.NewEncoder(w).Encode(data)
	}
}
