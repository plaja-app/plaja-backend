package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plaja-app/back-end/models"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// signupBody is the signup request body structure.
type signupBody struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// loginBody is the login request body structure.
type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// userUpdateGeneralBody is the user general information update request body structure.
type userUpdateGeneralBody struct {
	FirstName string
	LastName  string
}

// GetMe returns the model of the current models.User.
func (c *BaseController) GetMe(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// SignUp handles the signup request.
func (c *BaseController) SignUp(w http.ResponseWriter, r *http.Request) {
	var body signupBody

	// decode the request body
	err := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	if body.FirstName == "" ||
		body.LastName == "" ||
		!validateEmail(body.Email) ||
		len(body.Password) < 8 {
		http.Error(w, "Bad credentials provided", http.StatusBadRequest)
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusBadRequest)
		return
	}

	// create a new user model
	user := models.User{
		FirstName:  body.FirstName,
		LastName:   body.LastName,
		Email:      body.Email,
		Password:   string(hashedPassword),
		UserTypeID: 1,
	}

	// add user to the database
	result := c.App.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Error creating user", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login handles the login request.
func (c *BaseController) Login(w http.ResponseWriter, r *http.Request) {
	var body loginBody

	// get the email and password off request body
	err := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// look up the requested user
	var user models.User

	c.App.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		http.Error(w, "Invalid email", http.StatusUnauthorized)
		return
	}

	// compare sent in password with saved user password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// generate a JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 14).Unix(),
	})

	// sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(c.App.Env.JWTSecret))
	if err != nil {
		http.Error(w, "Failed to create JWT token", http.StatusBadRequest)
		return
	}

	// create and set a cookie
	cookie := http.Cookie{
		Name:     "pja_user_jwt",
		Path:     "/",
		Value:    tokenString,
		MaxAge:   3600 * 24 * 30,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

// Logout handles the logout request by invalidating the user's session cookie.
func (c *BaseController) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "pja_user_jwt",
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

// GetUsers returns the queried list of models.User.
func (c *BaseController) GetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("id")

	w.Header().Set("Content-Type", "application/json")

	var data []models.User

	if idParam == "all" {
		err := c.App.DB.Model(&models.User{}).Preload("UserType").Find(&data).Error
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
		err := c.App.DB.Where("id IN ?", intIds).Preload("UserType").Find(&data).Error
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

// UpdateUser handles the update request of the user's information.
func (c *BaseController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	var body userUpdateGeneralBody
	body.FirstName = r.FormValue("FirstName")
	body.LastName = r.FormValue("LastName")

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

	file, _, err := r.FormFile("ProfilePic")
	if err == nil && file != nil {
		defer file.Close()

		storagePath := "storage/users/profile-pictures"
		os.MkdirAll(storagePath, os.ModePerm)

		filePath := filepath.Join(storagePath, fmt.Sprintf("%d-%s", user.ID, "pp.png"))

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
		"FirstName": body.FirstName,
		"LastName":  body.LastName,
	}

	if fileURL != "" {
		updateData["ProfilePic"] = fileURL
	}

	result := c.App.DB.Model(&user).Where("id = ?", user.ID).Updates(updateData)
	if result.Error != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
