package google

import (
	"fmt"
	"net/http"
	"strings"
)

type Command struct{}

func (c *Command) Name() string {
	return "google"
}

func (c *Command) Prefixes() []string {
	return []string{"g"}
}

func (c *Command) Handle(args []string, response http.ResponseWriter, request *http.Request) error {
	joined := strings.Join(args, " ")
	searchURL := fmt.Sprintf("https://google.com/search?q=%s", joined)

	http.Redirect(response, request, searchURL, http.StatusSeeOther)
	return nil
}

func (c *Command) Help(response http.ResponseWriter, request *http.Request) error {
	response.Write(
		[]byte(fmt.Sprintf(
			"usage: gobunny %s <search query>",
			c.Name(),
		)),
	)

	return nil
}

func (c *Command) Readme(response http.ResponseWriter, request *http.Request) error {
	response.Write(
		[]byte(fmt.Sprintf(
			`"gobunny %s" provides convenient shorthand for performing Google searches\n\n`+
				`aliases: %s`,
			c.Name(),
			strings.Join(c.Prefixes(), " "),
		)),
	)

	return nil
}
