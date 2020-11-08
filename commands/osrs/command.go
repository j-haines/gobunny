package osrs

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"gobunny/commands"
)

type command struct {
	log *log.Logger
}

const wikiURL = "https://oldschool.runescape.wiki"

// NewCommand returns a Command for looking up Old School RuneScape wiki pages
func NewCommand(logger *log.Logger) commands.Command {
	return &command{
		log: logger,
	}
}

func (c *command) Aliases() []string {
	return []string{}
}

func (c *command) Name() string {
	return "osrs"
}

func (c *command) Handle(args commands.Arguments, response http.ResponseWriter, request *http.Request) error {
	var redirectURL string
	if len(args) == 0 {
		redirectURL = wikiURL
	} else {
		joined := strings.Join(args, " ")
		redirectURL = fmt.Sprintf("%s/w/%s", wikiURL, joined)
	}

	http.Redirect(response, request, redirectURL, http.StatusSeeOther)
	return nil
}

func (c *command) Help() string {
	return fmt.Sprintf(
		"usage: gobunny %s <search query>",
		c.Name(),
	)
}

func (c *command) Readme() string {
	return fmt.Sprintf(
		"'gobunny %s' allows looking up Old School RuneScape wiki articles",
		c.Name(),
	)
}
