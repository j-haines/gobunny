package google

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

const baseURL = "https://google.com"

// NewCommand returns a Command which provides support for Google searches
func NewCommand(logger *log.Logger) commands.Command {
	return &command{
		log: logger,
	}
}

// Handle implements commands.Command
func (c *command) Handle(response http.ResponseWriter, request *http.Request) error {
	query := chi.URLParam(request, "query")
	if len(query) == 0 {
		return handlers.Redirect(response, request, baseURL)
	}

	params, err := url.ParseQuery(query)
	if err != nil {
		return handlers.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	searchURL := fmt.Sprintf("%s/search?q=%s", baseURL, params.Encode())
	return handlers.Redirect(response, request, searchURL)
}

func (c *command) Routes() []commands.Route {
	return []commands.Route{
		{
			Method: http.MethodGet,
			Patterns: []string{
				"google {query}",
				"goog {query}",
				"g {query}",
			},
		},
	}
}
