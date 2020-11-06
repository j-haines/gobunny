package handlers

import "net/http"

// HealthCheckHandler returns a http.HandlerFunc for determining service health
func HealthCheckHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusOK)
	}
}
