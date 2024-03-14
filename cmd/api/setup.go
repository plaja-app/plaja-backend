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

	err = db.AutoMigrate(&models.CourseExerciseType{})
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

	err = db.AutoMigrate(&models.TeachingApplication{})
	if err != nil {
		return err
	}

	// populate tables with initial data
	err = createInitialUserTypes(db)
	if err != nil {
		return errors.New(fmt.Sprint("error creating initial user types:", err))
	}

	err = createInitialUsers(db)
	if err != nil {
		return errors.New(fmt.Sprint("error creating initial users:", err))
	}

	err = createInitialCourseLevels(db)
	if err != nil {
		return errors.New(fmt.Sprint("error creating initial course levels:", err))
	}

	err = createInitialCourseStatuses(db)
	if err != nil {
		return errors.New(fmt.Sprint("error creating initial course statuses:", err))
	}

	err = createInitialCourseCategories(db)
	if err != nil {
		return errors.New(fmt.Sprint("error creating initial course categories:", err))
	}

	err = createInitialCourses(db)
	if err != nil {
		return errors.New(fmt.Sprint("error creating initial courses:", err))
	}

	err = createInitialCourseExerciseTypes(db)
	if err != nil {
		return errors.New(fmt.Sprint("error creating initial course exercise types:", err))
	}

	return nil
}

// createInitialUserTypes creates initial user types in user_types table.
func createInitialUserTypes(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.UserType{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	initialData := []models.UserType{
		{Title: "Learner", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Educator", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Admin", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	if err := db.Create(&initialData).Error; err != nil {
		return err
	}

	return nil
}

// createInitialUsers creates initial users types in users table.
func createInitialUsers(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	initialData := []models.User{
		{
			FirstName:  "Plaja",
			LastName:   "Team",
			Email:      "mail@plaja.io",
			UserTypeID: 3,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&initialData).Error; err != nil {
		return err
	}

	return nil
}

// createInitialCourseLevels creates initial course levels in course_levels table.
func createInitialCourseLevels(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.CourseLevel{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	initialData := []models.CourseLevel{
		{Title: "Початковий", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Середній", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Високий", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	if err := db.Create(&initialData).Error; err != nil {
		return err
	}

	return nil
}

// createInitialCourseStatuses creates initial course statuses types in course_statuses table.
func createInitialCourseStatuses(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.CourseStatus{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	initialData := []models.CourseStatus{
		{Title: "draft", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "being validated", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "revisions required", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "published", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "suspended", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "archived", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	if err := db.Create(&initialData).Error; err != nil {
		return err
	}

	return nil
}

// createInitialCourseExerciseTypes creates initial course exercise types in course_exercise_types table.
func createInitialCourseExerciseTypes(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.CourseExerciseType{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	initialData := []models.CourseExerciseType{
		{Title: "article", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "video", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	if err := db.Create(&initialData).Error; err != nil {
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

	initialData := []models.CourseCategory{
		{Title: "Go", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "C++", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "C#", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Rust", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Ruby", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Python", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	if err := db.Create(&initialData).Error; err != nil {
		return err
	}

	return nil
}

// createInitialCourses creates initial courses in courses table.
func createInitialCourses(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.Course{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	initialData := []models.Course{
		{
			Title:            "Розробка сучасних веб-застосунків із Go",
			Thumbnail:        "https://img-c.udemycdn.com/course/480x270/3579383_3c67_4.jpg",
			ShortDescription: "Навчіться створювати сучасні вед-застосунки з Go, HTML, CSS та JavaScript. Курс від професора та знавця своєї справи.",
			InstructorID:     1,
			LevelID:          1,
			Price:            399,
			HasCertificate:   true,
			StatusID:         4,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},

		{
			Title:            "Використання мікросервісів у Go",
			Thumbnail:        "https://img-b.udemycdn.com/course/480x270/4606320_764e_2.jpg",
			ShortDescription: "Створюйте високодоступні, масштабовані, відмовостійкі розподілені додатки з Go.",
			InstructorID:     1,
			Price:            599,
			LevelID:          2,
			StatusID:         4,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},

		{
			Title:            "Svelte та SvelteKit: повний курс",
			Thumbnail:        "https://img-c.udemycdn.com/course/480x270/5557070_a5f3_3.jpg",
			ShortDescription: "Створюйте та розгортайте високопродуктивні, доступні, рендерингові веб-застосунки зі Svelte та SvelteKit.",
			InstructorID:     1,
			LevelID:          3,
			Price:            199,
			StatusID:         4,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},

		{
			Title:            "Шаблони проєктування в C++/C#",
			Thumbnail:        "https://bs-uploads.toptal.io/blackfish-uploads/components/blog_post_page/content/cover_image_file/cover_image/1285782/retina_500x200_COVER-dcbcd112f1d502d97d7f2467c1ce21da.png",
			ShortDescription: "Дізнайтеся про шаблони проєктування та їх застосування при розробці застосунків на C++.",
			InstructorID:     1,
			LevelID:          1,
			Price:            1199,
			HasCertificate:   true,
			StatusID:         4,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},

		{
			Title:            "Створення курсів на Plaja",
			Thumbnail:        "http://localhost:8080/api/v1/storage/service/courses/1-thumbnail.png",
			ShortDescription: "Курс для тих, хто хоче навчитися створювати власні відкриті або платні курси на платформі Plaja.",
			InstructorID:     1,
			LevelID:          2,
			StatusID:         4,
			Price:            0,
			HasCertificate:   true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	}

	if err := db.Create(&initialData).Error; err != nil {
		return err
	}

	return nil
}
