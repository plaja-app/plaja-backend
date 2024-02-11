package controllers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plaja-app/back-end/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// SignupBody is the signup request body structure.
type SignupBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginBody is the login request body structure.
type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *BaseController) Validate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logged in"))
	w.WriteHeader(http.StatusOK)
}

// SignUp handles the signup request.
func (c *BaseController) SignUp(w http.ResponseWriter, r *http.Request) {
	var body SignupBody

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
		Email:      body.Email,
		Password:   string(hashedPassword),
		UserTypeID: 1,
	}

	// add user to the database
	result := c.App.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Error creating user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Login handles the login request.
func (c *BaseController) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginBody

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
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	// compare sent in password with saved user password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusBadRequest)
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
		Value:    tokenString,
		MaxAge:   3600 * 24 * 30,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}
