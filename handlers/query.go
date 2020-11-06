package handlers

import (
	"gobunny/commands/registry"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

var errNotFound = []byte("error: not found")

// GetQueryHandler returns a http.HandlerFunc for handling querying Commands
func GetQueryHandler(r registry.Registry) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		query := chi.URLParam(request, "query")
		split := strings.Split(query, " ")

		if len(split) == 0 {
			response.WriteHeader(http.StatusNotFound)
			response.Write(errNotFound)
			return
		}

		name := split[0]

		command, found := r.Get(name)
		if !found {
			response.WriteHeader(http.StatusNotFound)
			response.Write(errNotFound)
			return
		}

		if len(split) < 2 {
			command.Help(response, request)
			return
		}

		first := strings.ToLower(split[1])
		switch first {
		case "?":
			fallthrough
		case "help":
			command.Help(response, request)
		case "readme":
			command.Readme(response, request)
		default:
			if err := command.Handle(split[1:], response, request); err != nil {
				response.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}
