package main

import (
	"fmt"
	"github.com/plaja-app/back-end/config"
	"log"
	"net/http"
)

var app config.AppConfig

func main() {
	err := setup(&app)
	if err != nil {
		log.Fatal()
	}

	fmt.Println("Running on port :8080...")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
