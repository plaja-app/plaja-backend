package controllers

import (
	"encoding/json"
	"github.com/plaja-app/back-end/models"
	"net/http"
)

// teachingApplicationBody is the teaching application request body structure.
type teachingApplicationBody struct {
	UserID         uint
	Experience     string
	Motivation     string
	PlatformChoice string
}

// CreateTeachingApplication creates a new models.TeachingApplication.
func (c *BaseController) CreateTeachingApplication(w http.ResponseWriter, r *http.Request) {
	var body teachingApplicationBody

	err := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var application models.TeachingApplication
	application = models.TeachingApplication{
		UserID:         body.UserID,
		Experience:     body.Experience,
		Motivation:     body.Motivation,
		PlatformChoice: body.PlatformChoice,
	}

	result := c.App.DB.Create(&application)
	if result.Error != nil {
		http.Error(w, "Error creating teaching application", http.StatusBadRequest)
		return
	}

	err = c.changeUserType(application.UserID, 2)
	if err != nil {
		http.Error(w, "Error promoting user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
