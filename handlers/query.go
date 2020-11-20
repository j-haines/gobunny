package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"gobunny/commands"
	gerrors "gobunny/errors"
	"gobunny/registry"
)

var errNotFound = []byte("error: not found")

// GetQueryHandler returns a http.HandlerFunc for handling querying Commands
func GetQueryHandler(r registry.Registry, logger *log.Logger) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		query := chi.URLParam(request, "query")
		split := strings.Split(query, " ")

		if len(split) == 0 {
			response.WriteHeader(http.StatusNotFound)
			if _, err := response.Write(errNotFound); err != nil {
				logger.Printf("response closed before handler finished: %s", err.Error())
			}
			return
		}

		name := split[0]

		command, found := r.Get(name)
		if !found {
			response.WriteHeader(http.StatusNotFound)
			if _, err := response.Write(errNotFound); err != nil {
				logger.Printf("response closed before handler finished: %s", err.Error())
			}

			return
		}

		if len(split) == 1 {
			if err := command.Handle(commands.Arguments{}, response, request); err != nil {
				response.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		first := strings.ToLower(split[1])
		var err error
		switch first {
		case "?":
			fallthrough
		case "help":
			if _, err = response.Write([]byte(command.Help())); err != nil {
				err = gerrors.NewErrResponseClosed(err)
			}
		case "readme":
			if _, err = response.Write([]byte(command.Readme())); err != nil {
				err = gerrors.NewErrResponseClosed(err)
			}
		default:
			err = command.Handle(split[1:], response, request)
		}

		var gerr *gerrors.ErrResponseClosed
		if errors.As(err, &gerr) {
			logger.Println(err.Error())
			return
		}

		if err != nil {
			logger.Println(err.Error())
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
