package google

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"gobunny/commands"
	"gobunny/errors"
)

type command struct {
	log *log.Logger
}

// NewCommand returns a Command which provides support for Google searches
func NewCommand(logger *log.Logger) commands.Command {
	return &command{
		log: logger,
	}
}

// Aliases implements commands.Command
func (c *command) Aliases() []string {
	return []string{"g"}
}

// Name implements commands.Command
func (c *command) Name() string {
	return "google"
}

// Handle implements commands.Command
func (c *command) Handle(args commands.Arguments, response http.ResponseWriter, request *http.Request) error {
	if len(args) == 0 {
		if _, err := response.Write([]byte(c.Help())); err != nil {
			return errors.NewErrResponseClosed(err)
		}

		return nil
	}

	joined := strings.Join(args, " ")
	searchURL := fmt.Sprintf("https://google.com/search?q=%s", joined)

	http.Redirect(response, request, searchURL, http.StatusSeeOther)
	return nil
}

// Help implements commands.Command
func (c *command) Help() string {
	return fmt.Sprintf(
		"usage: gobunny %s <search query>",
		c.Name(),
	)
}

// Readme implements commands.Command
func (c *command) Readme() string {
	return fmt.Sprintf(
		"'gobunny %s' provides convenient shorthand for performing Google searches\n\n"+
			"aliases: %s",
		c.Name(),
		strings.Join(c.Aliases(), ", "),
	)
}
