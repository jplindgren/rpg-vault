package main

import (
	"fmt"
	"net/http"
)

// Healthcheck returns if the service is available and the port it is running.
// swagger:route GET /healthcheck HealthCheck
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "port: %d \n", app.config.port)
}
