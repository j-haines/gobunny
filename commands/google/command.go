package google

import (
	"fmt"
	"gobunny/commands"
	"log"
	"net/http"
	"strings"
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
	joined := strings.Join(args, " ")
	searchURL := fmt.Sprintf("https://google.com/search?q=%s", joined)

	http.Redirect(response, request, searchURL, http.StatusSeeOther)
	return nil
}

// Help implements commands.Command
func (c *command) Help(response http.ResponseWriter, request *http.Request) error {
	_, err := response.Write(
		[]byte(fmt.Sprintf(
			"usage: gobunny %s <search query>",
			c.Name(),
		)),
	)

	if err != nil {
		c.log.Printf("response closed before handler finished: %s", err.Error())
	}

	return nil
}

// Readme implements commands.Command
func (c *command) Readme(response http.ResponseWriter, request *http.Request) error {
	_, err := response.Write(
		[]byte(fmt.Sprintf(
			`"gobunny %s" provides convenient shorthand for performing Google searches\n\n`+
				`aliases: %s`,
			c.Name(),
			strings.Join(c.Aliases(), ", "),
		)),
	)

	if err != nil {
		c.log.Printf("response closed before handler finished: %s", err.Error())
	}

	return nil
}
