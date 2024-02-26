package controllers

import (
	"encoding/json"
	"github.com/plaja-app/back-end/models"
	"net/http"
	"strconv"
	"strings"
)

// GetCourseCategory returns the queried list of models.CourseCategory.
func (c *BaseController) GetCourseCategory(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("id")
	sortParam := query.Get("sort")

	w.Header().Set("Content-Type", "application/json")

	var data []models.CourseCategory

	if idParam == "all" {
		if sortParam == "title" {
			c.App.DB.Order("title").Find(&data)
		} else {
			c.App.DB.Find(&data)
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
		c.App.DB.Where("id IN ?", intIds).Find(&data)
	}

	if len(data) == 0 {
		http.NotFound(w, r)
	} else {
		json.NewEncoder(w).Encode(data)
	}
}
