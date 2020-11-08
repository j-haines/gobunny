package random

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"gobunny/commands"
	"gobunny/errors"
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

func (c *command) Handle(args commands.Arguments, response http.ResponseWriter, request *http.Request) error {
	min := 0
	max := math.MaxInt64

	var err error
	if len(args) == 1 {
		if max, err = strconv.Atoi(args[0]); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return nil
		}
	}

	if len(args) == 2 {
		if min, err = strconv.Atoi(args[0]); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return nil
		}

		if max, err = strconv.Atoi(args[1]); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return nil
		}
	}

	result := rand.Int63n(int64(max-min)) + int64(min)
	if _, err = response.Write([]byte(strconv.Itoa(int(result)))); err != nil {
		return errors.NewErrResponseClosed(err)
	}

	return nil
}

func (c *command) Help() string {
	return fmt.Sprintf(
		"usage: \n\tgobunny %s\n\tgobunny %s [max]\n\tgobunny %s [min] [max]",
		c.Name(),
		c.Name(),
		c.Name(),
	)
}

func (c *command) Readme() string {
	return fmt.Sprintf(
		"'gobunny %s' provides random number generation\n\n"+
			"- 'gobunny %s [max]' generates a random number between 0 and [max]\n"+
			"- 'gobunny %s [min] [max] generates a random number between [min] and [max]\n\n"+
			"aliases: %s",
		c.Name(),
		c.Name(),
		c.Name(),
		strings.Join(c.Aliases(), ", "),
	)
}
