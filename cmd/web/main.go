package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/plaja-app/back-end/config"
	"github.com/plaja-app/back-end/controllers"
	md "github.com/plaja-app/back-end/middleware"
	"log"
	"net/http"
)

var App config.AppConfig

func main() {
	err := setup(&App)
	if err != nil {
		log.Fatal()
	}

	bc := controllers.NewBaseController(&App)
	bm := md.NewBaseMiddleware(&App)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/signup", bc.SignUp)
	r.Post("/login", bc.Login)

	r.Route("/admin", func(r chi.Router) {
		r.Use(bm.RequireAuth)
		r.Get("/validate", bc.Validate)
	})

	fmt.Println("Running on port :8080...")
	http.ListenAndServe(":8080", r)
}
