package google

import (
	"fmt"
	"gobunny/commands"
	"net/http"
	"strings"
)

// Command provides Google search support
type Command struct{}

// Aliases implements commands.Command
func (c *Command) Aliases() []string {
	return []string{"g"}
}

// Name implements commands.Command
func (c *Command) Name() string {
	return "google"
}

// Handle implements commands.Command
func (c *Command) Handle(args commands.Arguments, response http.ResponseWriter, request *http.Request) error {
	joined := strings.Join(args, " ")
	searchURL := fmt.Sprintf("https://google.com/search?q=%s", joined)

	http.Redirect(response, request, searchURL, http.StatusSeeOther)
	return nil
}

// Help implements commands.Command
func (c *Command) Help(response http.ResponseWriter, request *http.Request) error {
	response.Write(
		[]byte(fmt.Sprintf(
			"usage: gobunny %s <search query>",
			c.Name(),
		)),
	)

	return nil
}

// Readme implements commands.Command
func (c *Command) Readme(response http.ResponseWriter, request *http.Request) error {
	response.Write(
		[]byte(fmt.Sprintf(
			`"gobunny %s" provides convenient shorthand for performing Google searches\n\n`+
				`aliases: %s`,
			c.Name(),
			strings.Join(c.Aliases(), " "),
		)),
	)

	return nil
}
