package osrs

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"gobunny/commands"
	"gobunny/handlers"

	"github.com/go-chi/chi"
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

func (c *command) Handle(response http.ResponseWriter, request *http.Request) error {
	query := chi.URLParam(request, "query")
	if len(query) == 0 {
		return handlers.Redirect(response, request, wikiURL)
	}

	params, err := url.ParseQuery(query)
	if err != nil {
		return handlers.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	searchURL := fmt.Sprintf("%s/w/%s", wikiURL, params.Encode())
	return handlers.Redirect(response, request, searchURL)
}

func (c *command) Routes() []commands.Route {
	return []commands.Route{
		{
			Method:   http.MethodGet,
			Patterns: []string{"osrs {query}"},
		},
	}
}
