package random

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"

	"gobunny/commands"
	"gobunny/errors"
	"gobunny/handlers"

	"github.com/go-chi/chi"
)

type command struct {
	log *log.Logger
}

// NewCommand returns a Command which generates random numbers
func NewCommand(logger *log.Logger) commands.Command {
	return &command{
		log: logger,
	}
}

func (c *command) Aliases() []string {
	return []string{"r", "rand", "rng"}
}

func (c *command) Name() string {
	return "random"
}

func (c *command) Handle(response http.ResponseWriter, request *http.Request) error {
	var err error
	min := 0
	max := math.MaxInt64

	if param := chi.URLParam(request, "min"); len(param) > 0 {
		if min, err = strconv.Atoi(param); err != nil {
			return handlers.NewHTTPError(err.Error(), http.StatusBadRequest)
		}
	}

	if param := chi.URLParam(request, "max"); len(param) > 0 {
		if max, err = strconv.Atoi(param); err != nil {
			return handlers.NewHTTPError(err.Error(), http.StatusBadRequest)
		}
	}

	result := rand.Int63n(int64(max-min)) + int64(min)
	if _, err = response.Write([]byte(strconv.Itoa(int(result)))); err != nil {
		return errors.NewErrResponseClosed(err)
	}

	return nil
}

func (c *command) Routes() []commands.Route {
	return []commands.Route{
		{
			Method: http.MethodGet,
			Patterns: []string{
				"random",
				"random {max}",
				"random {min} {max}",
				"rand",
				"rand {max}",
				"rand {min} {max}",
				"rng",
				"rng {max}",
				"rng {min} {max}",
				"r",
				"r {max}",
				"r {min} {max}",
			},
		},
	}
}
