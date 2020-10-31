package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"gobunny/commands/google"
	"gobunny/commands/registry"
	"gobunny/handlers"
)

func makeRegistry() (registry.Registry, error) {
	r := registry.New()
	r.RegisterAll(&google.Command{})

	return r, nil
}

func main() {
	logger := log.New(os.Stderr, "gobunny: ", log.Lshortfile|log.Ltime)

	commands, err := makeRegistry()
	if err != nil {
		logger.Panicf("unexpected error creating command registry: %s", err.Error())
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/q/{query}", handlers.GetQueryHandler(commands))

	http.ListenAndServe(":8080", router)
}
