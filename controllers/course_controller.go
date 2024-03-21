package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/plaja-app/back-end/models"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	Title          string
	Categories     []courseCategory
	LevelID        uint
	HasCertificate bool
	// InstructorID   uint
}

// courseUpdateGeneralBody is the course general information update request body structure.
type courseUpdateGeneralBody struct {
	Title            string
	ShortDescription string
	Description      string
	Price            uint
	// InstructorID     uint
	CourseID uint
}

// GetCourses returns the queried list of models.Course.
func (c *BaseController) GetCourses(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := query.Get("id")
	userID := query.Get("user_id")
	statusID := query.Get("status_id")
	instructorID := query.Get("instructor_id")
	levelID := query.Get("level_id")
	hasCertificate := query.Get("has_certificate")
	sort := query.Get("sort")

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

	if userID != "" {
		dbQuery = dbQuery.Joins("JOIN enrollments ON enrollments.course_id = courses.id").Where("enrollments.user_id = ?", userID)
	}

	if hasCertificate != "" {
		hasCertBool, err := strconv.ParseBool(hasCertificate)
		if err != nil {
			http.Error(w, "Invalid format for has_certificate", http.StatusBadRequest)
			return
		}
		dbQuery = dbQuery.Where("has_certificate = ?", hasCertBool)
	}

	if sort != "" {
		direction := "ASC"
		if sort[0] == '-' {
			direction = "DESC"
			sort = sort[1:]
		}
		allowedSortFields := map[string]bool{
			"id":              true,
			"name":            true,
			"status_id":       true,
			"instructor_id":   true,
			"level_id":        true,
			"has_certificate": true,
			"updated_at":      true,
			"created_at":      true,
		}
		if _, ok := allowedSortFields[sort]; ok {
			dbQuery = dbQuery.Order(sort + " " + direction)
		} else {
			http.Error(w, "Invalid sort field", http.StatusBadRequest)
			return
		}
	}

	dbQuery = dbQuery.Preload("Instructor").Preload("Level").Preload("Categories")

	if err := dbQuery.Find(&courses).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(courses)
}

// CreateCourse creates a new models.Course.
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

	userCtx := r.Context().Value("user")
	if userCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, ok := userCtx.(models.User)
	if !ok {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	var course models.Course
	course = models.Course{
		Title:          body.Title,
		Thumbnail:      "http://localhost:8080/api/v1/storage/service/courses/no-thumbnail.png",
		Categories:     courseCategories,
		LevelID:        body.LevelID,
		StatusID:       1, // draft
		HasCertificate: body.HasCertificate,
		InstructorID:   user.ID,
	}

	result := c.App.DB.Create(&course)
	if result.Error != nil {
		http.Error(w, "Error creating course", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateGeneralCourse handles the update operation for general course information.
func (c *BaseController) UpdateGeneralCourse(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	var body courseUpdateGeneralBody
	body.Title = r.FormValue("Title")
	body.ShortDescription = r.FormValue("ShortDescription")
	body.Description = r.FormValue("Description")

	price, err := strconv.ParseUint(r.FormValue("Price"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}
	body.Price = uint(price)

	courseID, err := strconv.ParseUint(r.FormValue("CourseID"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid course id", http.StatusBadRequest)
		return
	}
	body.CourseID = uint(courseID)

	//instructorID, err := strconv.ParseUint(r.FormValue("InstructorID"), 10, 32)
	//if err != nil {
	//	http.Error(w, "Invalid instructor id", http.StatusBadRequest)
	//	return
	//}
	//body.InstructorID = uint(instructorID)

	userCtx := r.Context().Value("user")
	if userCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, ok := userCtx.(models.User)
	if !ok {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	var fileURL string

	file, _, err := r.FormFile("Thumbnail")
	if err == nil && file != nil {
		defer file.Close()

		storagePath := "storage/courses/thumbnails"
		os.MkdirAll(storagePath, os.ModePerm)

		filePath := filepath.Join(storagePath, fmt.Sprintf("%d-%s", body.CourseID, "thumbnail.png"))

		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to save the file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to write the file", http.StatusInternalServerError)
			return
		}

		fileURL = fmt.Sprintf("http://localhost:8080/api/v1/%s", filePath)
	}

	updateData := map[string]interface{}{
		"ShortDescription": body.ShortDescription,
		"Title":            body.Title,
		"Description":      body.Description,
		"Price":            body.Price,
		"InstructorID":     user.ID,
	}

	if fileURL != "" {
		updateData["Thumbnail"] = fileURL
	}

	var course models.Course
	result := c.App.DB.Model(&course).Where("id = ?", body.CourseID).Updates(updateData)
	if result.Error != nil {
		http.Error(w, "Error updating course", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
