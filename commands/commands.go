package commands

import (
	"net/http"
)

type Command interface {
	Name() string
	Prefixes() []string

	Handle([]string, http.ResponseWriter, *http.Request) error

	Help(http.ResponseWriter, *http.Request) error
	Readme(http.ResponseWriter, *http.Request) error
}
