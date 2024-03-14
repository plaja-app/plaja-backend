package controllers

import (
	"github.com/plaja-app/back-end/models"
	"strings"
)

// changeUserType changes the type of the specified user to the new one provided.
// Note: 1 ("Learner") < 2 ("Educator") < 3 ("Admin").
func (c *BaseController) changeUserType(userID uint, newType uint) error {
	var user models.User
	return c.App.DB.Model(&user).Where("id = ?", userID).Update("user_type_id", newType).Error
}

// calculateExerciseLength calculates ant returns the approximate length of the exercise (in minutes).
func calculateExerciseLength(content string) uint {
	words := strings.Fields(content)

	wordCount := len(words)

	length := uint(wordCount) / 200 // approximate reading speed in Ukrainian

	return length
}
