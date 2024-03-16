package controllers

import (
	"encoding/json"
	"github.com/plaja-app/back-end/models"
	"net/http"
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

// GetCourseCategoriesStats returns a category and the number of courses associated with it.
func (c *BaseController) GetCourseCategoriesStats(w http.ResponseWriter, r *http.Request) {
	type CategoryStat struct {
		Title        string
		CoursesCount int
	}

	var stats []CategoryStat

	err := c.App.DB.Table("course_categories").Select("course_categories.title, COUNT(course_categories_junction.course_id) as courses_count").
		Joins("LEFT JOIN course_categories_junction ON course_categories.id = course_categories_junction.course_category_id").
		Joins("LEFT JOIN courses ON course_categories_junction.course_id = courses.id").
		Group("course_categories.title").
		Scan(&stats).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetCourseCategoriesAndLevelsStats returns a course level for a category and the number of courses associated with it.
func (c *BaseController) GetCourseCategoriesAndLevelsStats(w http.ResponseWriter, r *http.Request) {
	type CategoryLevelStat struct {
		CategoryTitle string
		LevelTitle    string
		CoursesCount  int
	}

	var stats []CategoryLevelStat

	err := c.App.DB.Table("course_categories").
		Select("course_categories.title as category_title, course_levels.title as level_title, COUNT(courses.id) as courses_count").
		Joins("LEFT JOIN course_categories_junction ON course_categories.id = course_categories_junction.course_category_id").
		Joins("LEFT JOIN courses ON course_categories_junction.course_id = courses.id").
		Joins("LEFT JOIN course_levels ON courses.level_id = course_levels.id"). // Assuming courses have a level_id column
		Group("course_categories.title, course_levels.title").
		Scan(&stats).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
