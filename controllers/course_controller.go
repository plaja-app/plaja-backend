package controllers

import (
	"encoding/json"
	"github.com/plaja-app/back-end/models"
	"net/http"
	"strconv"
	"strings"
)

// courseCreationBody is the course creation request body structure.
type courseCreationBody struct {
	Title      string                  `json:"fullName"`
	Categories []models.CourseCategory `json:"categories"`
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

	var course models.Course

	course = models.Course{
		Title:      body.Title,
		Categories: body.Categories,
	}

	result := c.App.DB.Create(&course)
	if result.Error != nil {
		http.Error(w, "Error creating user", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
