package commands

import (
	"net/http"
)

// Arguments are passed to the Command Handle function
type Arguments []string

// Command is the interface for implementing new GoBunny commands
type Command interface {
	// Aliases returns the set of strings that will invoke the Command
	Aliases() []string

	// Name returns the primary alias used for the Command
	Name() string

	// Handle is called by the GoBunny HTTP Handler
	Handle(Arguments, http.ResponseWriter, *http.Request) error

	// Help is a http.HandlerFunc for displaying a Command help page
	Help(http.ResponseWriter, *http.Request) error

	// Readme is a http.HandlerFunc for displaying a Command readme page
	Readme(http.ResponseWriter, *http.Request) error
}
