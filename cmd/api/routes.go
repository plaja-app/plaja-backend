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

	r.Get("/api/v1/course_categories", c.Controller.GetCourseCategory)

	r.Post("/signup", c.Controller.SignUp)
	r.Post("/login", c.Controller.Login)

	r.Route("/admin", func(r chi.Router) {
		r.Use(m.Middleware.RequireAuth)
		r.Get("/validate", c.Controller.Validate)
	})

	return r
}
