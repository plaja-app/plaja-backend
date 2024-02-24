package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/plaja-app/back-end/config"
	c "github.com/plaja-app/back-end/controllers"
	m "github.com/plaja-app/back-end/middleware"
	"github.com/plaja-app/back-end/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func setup(app *config.AppConfig) error {
	// Get environment variables
	env, err := loadEvnVariables()
	if err != nil {
		return err
	}

	app.Env = env

	// Connect to the database and run migrations
	db, err := connectToPostgresAndMigrate(env)
	if err != nil {
		return err
	}

	app.DB = db

	// Run database migrations
	err = runDatabaseMigrations(db)
	if err != nil {
		return err
	}

	// Create controllers
	bc := c.NewBaseController(app)
	c.NewControllers(bc)

	// Create middleware
	bm := m.NewBaseMiddleware(app)
	m.NewMiddleware(bm)

	return nil
}

// loadEvnVariables loads variables from the .env file.
func loadEvnVariables() (*config.EnvVariables, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error getting environment variables: %v", err)
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPass := os.Getenv("POSTGRES_PASS")
	postgresDBName := os.Getenv("POSTGRES_DBNAME")
	jwtSecret := os.Getenv("JWT_SECRET")

	return &config.EnvVariables{
		PostgresHost:   postgresHost,
		PostgresUser:   postgresUser,
		PostgresPass:   postgresPass,
		PostgresDBName: postgresDBName,
		JWTSecret:      jwtSecret,
	}, nil
}

// connectToPostgresAndMigrate initializes a PostgreSQL db session and runs GORM migrations.
func connectToPostgresAndMigrate(env *config.EnvVariables) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		env.PostgresHost, env.PostgresUser, env.PostgresDBName, env.PostgresPass)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("could not connect: ", err)
	}

	return db, nil
}

func runDatabaseMigrations(db *gorm.DB) error {
	// create tables
	err := db.AutoMigrate(&models.CourseStatus{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.CourseCategory{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.UserType{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.Course{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.CourseCertificate{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.EnrollmentStatus{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.CourseExerciseCategory{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.CourseExercise{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.Enrollment{})
	if err != nil {
		return err
	}

	// populate tables with initial data
	err = createInitialUserTypes(db)
	if err != nil {
		return errors.New(fmt.Sprint("error creating initial user types:", err))
	}

	err = createInitialCourseCategories(db)
	if err != nil {
		return errors.New(fmt.Sprint("error creating initial user types:", err))
	}

	return nil
}

// createInitialUserTypes creates initial account types in account_types table.
func createInitialUserTypes(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.UserType{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	initialTypes := []models.UserType{
		{Title: "Learner", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Educator", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Admin", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	if err := db.Create(&initialTypes).Error; err != nil {
		return err
	}

	return nil
}

// createInitialCourseCategories creates initial course categories in course_categories table.
func createInitialCourseCategories(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.CourseCategory{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	initialTypes := []models.CourseCategory{
		{Title: "Go", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "C++", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "C#", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Rust", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Ruby", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Python", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	if err := db.Create(&initialTypes).Error; err != nil {
		return err
	}

	return nil
}
