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

// GetCourseCertificates returns the queried list of models.CourseCertificate.
func (c *BaseController) GetCourseCertificates(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := query.Get("id")
	userID := query.Get("user_id")
	courseID := query.Get("course_id")

	w.Header().Set("Content-Type", "application/json")

	var certificates []models.CourseCertificate
	dbQuery := c.App.DB

	if id != "" {
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

	if userID != "" {
		dbQuery = dbQuery.Where("user_id = ?", userID)
	}

	if courseID != "" {
		dbQuery = dbQuery.Where("course_id = ?", courseID)
	}

	dbQuery = dbQuery.Preload("Course")

	if err := dbQuery.Find(&certificates).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	//if len(certificates) == 0 {
	//	http.NotFound(w, r)
	//	return
	//}

	json.NewEncoder(w).Encode(certificates)
}
