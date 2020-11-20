package commands

import (
	"net/http"
)

// Arguments are passed to the Command Handle function
type Arguments []string

// Route is a Method and Route pattern pair for calling the Command handler
type Route struct {
	Method   string
	Patterns []string
}

// Command is the interface for implementing new GoBunny commands
type Command interface {
	// Handle is called by the GoBunny HTTP Handler
	Handle(http.ResponseWriter, *http.Request) error

	// Routes returns the URL path routing rules for calling a Command
	Routes() []Route
}
