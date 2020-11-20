package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"gobunny/commands"
	"gobunny/commands/google"
	"gobunny/commands/osrs"
	"gobunny/commands/random"
	"gobunny/handlers"
)

func makeCommandHandler(l *log.Logger, c commands.Command) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if err := c.Handle(resp, req); err != nil {
			if herr, ok := err.(*handlers.HTTPError); ok {
				resp.WriteHeader(herr.Status)
				if _, err := resp.Write([]byte(herr.Error())); err != nil {
					l.Printf("unexpected error writing to response: '%s'", err.Error())
				}

				return
			}

			resp.WriteHeader(http.StatusInternalServerError)
			l.Printf("unrecognized error in command handler: '%s'", err.Error())
		}
	}
}

func makeCommandRoutes(l *log.Logger, c ...commands.Command) (func(chi.Router), error) {
	return func(router chi.Router) {
		for _, command := range c {
			for _, route := range command.Routes() {
				var fn func(string, http.HandlerFunc)
				switch route.Method {
				case http.MethodDelete:
					fn = router.Delete
				case http.MethodGet:
					fn = router.Get
				case http.MethodPatch:
					fn = router.Patch
				case http.MethodPost:
					fn = router.Post
				case http.MethodPut:
					fn = router.Put
				}

				for _, pattern := range route.Patterns {
					fn(fmt.Sprintf("/%s", pattern), makeCommandHandler(l, command))
				}
			}
		}
	}, nil
}

func main() {
	logger := log.New(os.Stderr, "gobunny: ", log.Lshortfile|log.Ltime)
	host := flag.String("host", "", "hostname to bind http server to")
	port := flag.Int("port", 8080, "port to bind http server to")

	routesFn, err := makeCommandRoutes(
		logger,
		google.NewCommand(logger),
		osrs.NewCommand(logger),
		random.NewCommand(logger),
	)
	if err != nil {
		logger.Fatalf("unexpected error while creating registry routes: '%s'", err.Error())
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Get("/health", handlers.HealthCheckHandler())
	router.Route("/q", routesFn)

	bindAddr := fmt.Sprintf("%s:%d", *host, *port)
	logger.Printf("starting http server on %s", bindAddr)
	if err := http.ListenAndServe(bindAddr, router); err != nil {
		logger.Fatalf("unexpected error while running http server: '%s'", err.Error())
	}
}
