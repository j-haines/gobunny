package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"gobunny/commands"
	"gobunny/commands/google"
	"gobunny/commands/marketwatch"
	"gobunny/commands/osrs"
	"gobunny/commands/random"
	"gobunny/handlers"
	"gobunny/store"
	"gobunny/store/redis"
)

type cliArgs struct {
	HTTPServerHost string
	HTTPServerPort int

	RedisAddress    string
	RedisPort       int
	RedisPassword   string
	RedisDatabaseID int
}

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

func makeStore(ctx context.Context, args *cliArgs) (store.Store, error) {
	db, err := redis.NewStore(
		ctx,
		redis.Config{
			HostAddress: fmt.Sprintf("%s:%d", args.RedisAddress, args.RedisPort),
			Password:    args.RedisPassword,
			DatabaseID:  args.RedisDatabaseID,
		},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func parseArgs(args *cliArgs) {
	flag.StringVar(&args.HTTPServerHost, "host", "localhost", "address to bind the HTTP server to")
	flag.IntVar(&args.HTTPServerPort, "port", 8080, "port to bind the HTTP server to")

	flag.StringVar(&args.RedisAddress, "redis-host", "localhost", "address of the Redis server")
	flag.IntVar(&args.RedisPort, "redis-port", 6379, "port of the Redis server")
	flag.StringVar(&args.RedisPassword, "redis-password", "", "password for the Redis server")
	flag.IntVar(&args.RedisDatabaseID, "redis-dbid", 0, "the Redis server database to use")
	flag.Parse()
}

func main() {
	args := &cliArgs{}
	parseArgs(args)

	ctx := context.Background()
	logger := log.New(os.Stderr, "gobunny: ", log.Lshortfile|log.Ltime)

	db, err := makeStore(ctx, args)
	if err != nil {
		logger.Fatalf("unexpected error creating data store: %s", err.Error())
	}

	routesFn, err := makeCommandRoutes(
		logger,
		google.NewCommand(logger),
		osrs.NewCommand(logger),
		random.NewCommand(logger),
		marketwatch.NewCommand(logger, db),
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

	bindAddr := fmt.Sprintf("%s:%d", args.HTTPServerHost, args.HTTPServerPort)
	logger.Printf("starting http server on %s", bindAddr)
	if err := http.ListenAndServe(bindAddr, router); err != nil {
		logger.Fatalf("unexpected error while running http server: '%s'", err.Error())
	}
}
