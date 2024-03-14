package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/plaja-app/back-end/models"
	"gorm.io/gorm"
	"image"
	"log"
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

// CreateCourseCertificate creates a new models.CourseCertificate.
func (c *BaseController) CreateCourseCertificate(w http.ResponseWriter, r *http.Request) {
	_, err := generateCertificate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

// loadImage loads an image from the specified path.
func loadImage(path string) (image.Image, error) {
	img, err := gg.LoadPNG(path)
	if err != nil {
		return nil, fmt.Errorf("error loading image from %s: %v", path, err)
	}
	return img, nil
}

// drawString draws the string with specified parameters.
func drawString(dc *gg.Context, text string, x, y float64, width float64, alignment gg.Align, fontSize float64, font string) error {
	if err := dc.LoadFontFace(fmt.Sprintf("./storage/service/fonts/%s.ttf", font), fontSize); err != nil {
		return fmt.Errorf("error loading font %s: %v", font, err)
	}
	dc.DrawStringWrapped(text, x, y, 0, 0, width, 1.5, alignment)

	return nil
}

// generateCertificate generates a new certificate (.jpg) and saves it to ./storage/certificates.
// Returns the path to the generated certificate and an error.
func generateCertificate() (string, error) {
	dc := gg.NewContext(1200, 800)

	// Add background
	img, err := loadImage("./storage/service/certificates/background.png")
	if err != nil {
		return "", err
	}
	dc.DrawImage(img, 0, 0)

	// Add logo and signature
	img, err = loadImage("./storage/service/logo/logo-dark.png")
	if err != nil {
		return "", err
	}
	dc.DrawImage(img, 63, 590)

	img, err = loadImage("./storage/service/other/signature.png")
	if err != nil {
		return "", err
	}
	dc.DrawImage(img, 875, 615)

	// Add base text
	dc.SetRGB(0, 0, 0)

	err = drawString(dc, "цей сертифікат засвідчує, що", 110, 220, 500, gg.AlignLeft, 24, "Onest-Regular")
	if err != nil {
		return "", err
	}

	err = drawString(dc, "успішно завершив (-ла) курс", 110, 380, 500, gg.AlignLeft, 24, "Onest-Regular")
	if err != nil {
		return "", err
	}

	err = drawString(dc, "інструктор:", 110, 540, 200, gg.AlignLeft, 24, "Onest-Regular")
	if err != nil {
		return "", err
	}

	err = drawString(dc, "тривалість:", 110, 572, 200, gg.AlignLeft, 24, "Onest-Regular")
	if err != nil {
		return "", err
	}

	err = drawString(dc, "засновник, Plaja", 910, 695, 200, gg.AlignRight, 24, "Onest-Regular")
	if err != nil {
		return "", err
	}

	// Add semi-transparent text
	dc.SetRGBA(0, 0, 0, 0.3)
	err = drawString(dc, fmt.Sprintf("ідентифікатор: %d", 0), 615, 85, 500, gg.AlignRight, 14, "Onest-Regular")
	if err != nil {
		return "", err
	}

	err = drawString(dc, fmt.Sprintf("видано %s", "11 березня 2023"), 615, 105, 500, gg.AlignRight, 14, "Onest-Regular")
	if err != nil {
		return "", err
	}

	// Add actual information with different font sizes
	dc.SetRGB(0, 0, 0)
	err = drawString(dc, "Ім'я Прізвище", 110, 257, 980, gg.AlignLeft, 56, "Onest-Medium")
	if err != nil {
		return "", err
	}

	err = drawString(dc, "Назва курсу", 110, 415, 980, gg.AlignLeft, 36, "Onest-Medium")
	if err != nil {
		return "", err
	}

	err = drawString(dc, "N/A", 246, 540, 500, gg.AlignLeft, 24, "Onest-Medium")
	if err != nil {
		return "", err
	}

	err = drawString(dc, "N/A", 246, 572, 500, gg.AlignLeft, 24, "Onest-Medium")
	if err != nil {
		return "", err
	}

	// Save the final image
	path := fmt.Sprintf("/storage/certificates/%d-certificate.png", 0)
	if err := dc.SavePNG(fmt.Sprintf(".%s", path)); err != nil {
		log.Fatalf("error saving image to %s: %v", path, err)
	}

	return path, nil
}
