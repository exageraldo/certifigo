package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	protected := alice.New(app.requireAuthenticatedUser)
	mux.Handle("GET /healthcheck", protected.ThenFunc(app.healthcheckGet))
	mux.Handle("GET /certificates/default_config", protected.ThenFunc(app.defaultConfigCertificatesGet))
	mux.Handle("POST /certificates/generate", protected.ThenFunc(app.generateCertificatesPost))
	mux.Handle("POST /certificates/validate", protected.ThenFunc(app.validateCertificatesPost))

	// Handle the not found route.
	// It will catch all unmatched routes and return a 404 error.
	// This should be the last route in the mux.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.httpError(w, r, "Not found", http.StatusNotFound)
	})

	standard := alice.New(app.recoverPanic, app.requestID, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
