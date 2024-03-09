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

	id := query.Get("id")
	statusID := query.Get("status_id")
	instructorID := query.Get("instructor_id")
	levelID := query.Get("level_id")
	hasCertificate := query.Get("has_certificate")

	w.Header().Set("Content-Type", "application/json")

	var courses []models.Course
	dbQuery := c.App.DB

	if id != "" {
		if id != "all" {
			ids := strings.Split(id, ",")
			var intIds []int
			for _, idStr := range ids {
				intId, err := strconv.Atoi(idStr)
				if err != nil {
					http.Error(w, "Invalid ID format", http.StatusBadRequest)
					return
				}
				intIds = append(intIds, intId)
			}
			dbQuery = dbQuery.Where("id IN ?", intIds)
		}
	}

	if statusID != "" {
		dbQuery = dbQuery.Where("status_id = ?", statusID)
	}

	if instructorID != "" {
		dbQuery = dbQuery.Where("instructor_id = ?", instructorID)
	}

	if levelID != "" {
		dbQuery = dbQuery.Where("level_id = ?", levelID)
	}

	if hasCertificate != "" {
		hasCertBool, err := strconv.ParseBool(hasCertificate)
		if err != nil {
			http.Error(w, "Invalid format for has_certificate", http.StatusBadRequest)
			return
		}
		dbQuery = dbQuery.Where("has_certificate = ?", hasCertBool)
	}

	dbQuery = dbQuery.Preload("Instructor").Preload("Level")

	if err := dbQuery.Find(&courses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if len(courses) == 0 {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(courses)
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
		Thumbnail:    "http://localhost:8080/api/v1/storage/service/courses/no_thumbnail.png",
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
