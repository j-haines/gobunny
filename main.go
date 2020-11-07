package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"gobunny/commands/google"
	"gobunny/handlers"
	"gobunny/registry"
)

func makeRegistry(logger *log.Logger) (registry.Registry, error) {
	r := registry.New()
	if err := r.RegisterAll(google.NewCommand(logger)); err != nil {
		return nil, err
	}

	return r, nil
}

func main() {
	logger := log.New(os.Stderr, "gobunny: ", log.Lshortfile|log.Ltime)
	host := flag.String("host", "", "hostname to bind http server to")
	port := flag.Int("port", 8080, "port to bind http server to")

	commands, err := makeRegistry(logger)
	if err != nil {
		logger.Fatalf("unexpected error creating command registry: %s", err.Error())
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Get("/health", handlers.HealthCheckHandler())
	router.Get("/q/{query}", handlers.GetQueryHandler(commands, logger))

	bindAddr := fmt.Sprintf("%s:%d", *host, *port)
	logger.Printf("starting http server on %s", bindAddr)
	if err := http.ListenAndServe(bindAddr, router); err != nil {
		logger.Fatalf("unexpected error while running http server: %s", err.Error())
	}
}
