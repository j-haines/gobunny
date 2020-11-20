package handlers

import (
	"fmt"
	"net/http"
	"net/url"
)

// Redirect redirects the response to a destination URL
func Redirect(
	response http.ResponseWriter,
	request *http.Request,
	destination string,
) error {
	if _, err := url.ParseRequestURI(destination); err != nil {
		return NewHTTPError(
			fmt.Sprintf("malformed redirect URL: '%s'", err.Error()),
			http.StatusBadRequest,
		)
	}

	http.Redirect(response, request, destination, http.StatusSeeOther)
	return nil
}
