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
	host := flag.String("host", "", "hostname to bind http server to")
	port := flag.Int("port", 8080, "port to bind http server to")

	commands, err := makeRegistry()
	if err != nil {
		logger.Panicf("unexpected error creating command registry: %s", err.Error())
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Get("/q/{query}", handlers.GetQueryHandler(commands))

	bindAddr := fmt.Sprintf("%s:%d", *host, *port)
	logger.Printf("starting http server on %s", bindAddr)
	http.ListenAndServe(bindAddr, router)
}
