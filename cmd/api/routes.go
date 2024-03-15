package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/plaja-app/back-end/config"
	c "github.com/plaja-app/back-end/controllers"
	m "github.com/plaja-app/back-end/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/api/v1/course-categories", c.Controller.GetCourseCategories)
	r.Get("/api/v1/course-levels", c.Controller.GetCourseLevels)
	r.Get("/api/v1/courses", c.Controller.GetCourses)

	r.Get("/api/v1/course-certificates", c.Controller.GetCourseCertificates)

	r.Post("/api/v1/users/signup", c.Controller.SignUp)
	r.Post("/api/v1/users/login", c.Controller.Login)
	r.Post("/api/v1/users/logout", c.Controller.Logout)

	r.Get("/api/v1/users", c.Controller.GetUsers)

	r.Get("/api/v1/course-exercises", c.Controller.GetCourseExercises)

	r.Group(func(r chi.Router) {
		r.Use(m.Middleware.RequireAuth)
		r.Get("/api/v1/users/getme", c.Controller.GetMe)
		r.Post("/api/v1/users/update-user", c.Controller.UpdateUser)

		r.Post("/api/v1/courses/create", c.Controller.CreateCourse)
		r.Post("/api/v1/courses/update-general", c.Controller.UpdateGeneralCourse)

		r.Post("/api/v1/teaching-applications/create", c.Controller.CreateTeachingApplication)

		r.Post("/api/v1/course-certificates/create", c.Controller.CreateCourseCertificate)
		r.Post("/api/v1/course-exercises/create-update", c.Controller.CreateOrUpdateCourseExercises)
	})

	r.Get("/api/v1/storage/*", c.Controller.GetImage)

	return r
}
