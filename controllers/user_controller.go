package controllers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plaja-app/back-end/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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

// GetMe returns the model of the current user (see models.User).
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
	//var body signupBody
	//
	//// decode the request body
	//err := json.NewDecoder(r.Body).Decode(&body)
	//defer r.Body.Close()
	//if err != nil {
	//	http.Error(w, "Failed to read body", http.StatusBadRequest)
	//	return
	//}
	//
	//// hash password
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	//if err != nil {
	//	http.Error(w, "Error hashing password", http.StatusBadRequest)
	//	return
	//}
	//
	//// parse name
	//fullNameSlice := strings.Split(body.FullName, " ")
	//firstName, lastName := fullNameSlice[0], fullNameSlice[1]
	//
	//// create a new user model
	//user := models.User{
	//	FirstName:  firstName,
	//	LastName:   lastName,
	//	Email:      body.Email,
	//	Password:   string(hashedPassword),
	//	UserTypeID: 1,
	//}
	//
	//// add user to the database
	//result := c.App.DB.Create(&user)
	//if result.Error != nil {
	//	http.Error(w, "Error creating user", http.StatusConflict)
	//	return
	//}
	//
	//w.WriteHeader(http.StatusCreated)
}
