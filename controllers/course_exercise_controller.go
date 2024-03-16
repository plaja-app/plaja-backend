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
	Exercises         []ExerciseInput
	InstructorID      uint
	CourseID          uint
	ExercisesToDelete []uint
}

// ExerciseInput is the exercise input structure.
type ExerciseInput struct {
	ID      uint
	Title   string
	Content string
}

// CreateOrUpdateCourseExercises creates new records of type models.CourseExercise or
// updates the exising ones if ID is provided.
func (c *BaseController) CreateOrUpdateCourseExercises(w http.ResponseWriter, r *http.Request) {
	var body CourseExerciseInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var course models.Course
	if err := c.App.DB.First(&course, body.CourseID).Error; err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	if course.InstructorID != body.InstructorID {
		http.Error(w, "Instructor not authorized to add or update exercises in this course", http.StatusUnauthorized)
		return
	}

	for _, ex := range body.Exercises {
		if ex.ID != 0 {
			var existingExercise models.CourseExercise
			if err := c.App.DB.First(&existingExercise, "id = ?", ex.ID).Error; err != nil {
				http.Error(w, "Exercise not found", http.StatusNotFound)
				return
			}

			existingExercise.Title = ex.Title
			existingExercise.Content = ex.Content
			existingExercise.Length = calculateExerciseLength(ex.Content)
			existingExercise.UpdatedAt = time.Now()

			if err := c.App.DB.Save(&existingExercise).Error; err != nil {
				http.Error(w, "Failed to update exercise", http.StatusInternalServerError)
				return
			}
		} else {
			newExercise := models.CourseExercise{
				CourseID:  body.CourseID,
				TypeID:    1, // Assuming TypeID is a fixed value as per your original logic
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
	}

	if len(body.ExercisesToDelete) > 0 {
		c.App.DB.Where("id IN ?", body.ExercisesToDelete).Delete(&models.CourseExercise{})
	}

	var totalCourseLength uint
	if err := c.App.DB.Model(&models.CourseExercise{}).Where("course_id = ?", course.ID).Select("sum(length)").Row().Scan(&totalCourseLength); err != nil {
		http.Error(w, "Failed to calculate total course length", http.StatusInternalServerError)
		return
	}

	course.Length = totalCourseLength
	if err := c.App.DB.Save(&course).Error; err != nil {
		http.Error(w, "Failed to update course length", http.StatusInternalServerError)
		return
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
		c.App.DB.Order("id").Where("course_id = ?", courseIDInt).Find(&data)
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
		c.App.DB.Order("id").Where("id IN ? AND course_id = ?", intIDs, courseIDInt).Find(&data)
	}

	if len(data) == 0 {
		data = make([]models.CourseExercise, 0)
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(data)
}
