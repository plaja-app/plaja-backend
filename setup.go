package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vladyslavpavlenko/plaja/back-end/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

// envVariables holds environment variables used in the application.
type envVariables struct {
	PostgresHost   string
	PostgresUser   string
	PostgresPass   string
	PostgresDBName string
}

func setup() error {
	// Get environment variables
	env, err := loadEvnVariables()
	if err != nil {
		return err
	}

	// Connect to the database and run migrations
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		env.PostgresHost, env.PostgresUser, env.PostgresDBName, env.PostgresPass)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("could not connect: ", err)
	}

	db.AutoMigrate(&models.CourseStatus{})
	db.AutoMigrate(&models.CourseCategory{})
	db.AutoMigrate(&models.AccountType{})
	db.AutoMigrate(&models.Account{})
	db.AutoMigrate(&models.Course{})
	db.AutoMigrate(&models.CourseCertificate{})
	db.AutoMigrate(&models.EnrollmentStatus{})
	db.AutoMigrate(&models.CourseExerciseCategory{})
	db.AutoMigrate(&models.CourseExercise{})
	db.AutoMigrate(&models.Enrollment{})

	return nil
}

// loadEvnVariables loads variables from the .env file.
func loadEvnVariables() (*envVariables, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error getting environment variables: %v", err)
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPass := os.Getenv("POSTGRES_PASS")
	postgresDBName := os.Getenv("POSTGRES_DBNAME")

	return &envVariables{
		PostgresHost:   postgresHost,
		PostgresUser:   postgresUser,
		PostgresPass:   postgresPass,
		PostgresDBName: postgresDBName,
	}, nil
}
