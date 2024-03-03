package controllers

import (
	"encoding/json"
	"github.com/plaja-app/back-end/models"
	"net/http"
	"strconv"
	"strings"
)

// courseCategory is the models.CourseCategory DTO.
type courseCategory struct {
	ID    uint
	Title string
}

// courseCreationBody is the course creation request body structure.
type courseCreationBody struct {
	Title        string
	Categories   []courseCategory
	LevelID      uint
	InstructorID uint
}

// GetCourses returns the queried list of models.Course.
func (c *BaseController) GetCourses(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("id")

	w.Header().Set("Content-Type", "application/json")

	var data []models.Course

	if idParam == "all" {
		err := c.App.DB.Model(&models.Course{}).Preload("Instructor").Preload("Level").Find(&data).Error
		if err != nil {
			return
		}
	} else {
		ids := strings.Split(idParam, ",")
		var intIds []int
		for _, idStr := range ids {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}
			intIds = append(intIds, id)
		}
		err := c.App.DB.Where("id IN ?", intIds).Preload("Instructor").Preload("Level").Find(&data).Error
		if err != nil {
			return
		}
	}

	if len(data) == 0 {
		http.NotFound(w, r)
	} else {
		json.NewEncoder(w).Encode(data)
	}
}

// CreateCourse creates a new course in the courses table.
func (c *BaseController) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var body courseCreationBody

	err := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var courseCategories []models.CourseCategory
	for _, c := range body.Categories {
		category := models.CourseCategory{
			ID:    c.ID,
			Title: c.Title,
		}

		courseCategories = append(courseCategories, category)
	}

	var course models.Course
	course = models.Course{
		Title:        body.Title,
		Categories:   courseCategories,
		LevelID:      body.LevelID,
		StatusID:     1, // draft
		InstructorID: body.InstructorID,
	}

	result := c.App.DB.Create(&course)
	if result.Error != nil {
		http.Error(w, "Error creating course", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
