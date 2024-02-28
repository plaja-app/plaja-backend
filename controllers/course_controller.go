package controllers

import (
	"encoding/json"
	"github.com/plaja-app/back-end/models"
	"net/http"
	"strconv"
	"strings"
)

// GetCourse returns the queried list of models.Course.
func (c *BaseController) GetCourse(w http.ResponseWriter, r *http.Request) {
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
