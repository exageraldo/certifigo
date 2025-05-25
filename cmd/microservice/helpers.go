package main

import (
	"net/http"
)

func (app *application) httpError(w http.ResponseWriter, r *http.Request, error string, code int) {
	w.WriteHeader(code)
	app.logger.Error(
		"http error",
		"code", code,
		"error", error,
		"ip", r.RemoteAddr,
		"proto", r.Proto,
		"method", r.Method,
		"uri", r.URL.RequestURI(),
		"x_request_id", GetReqID(r.Context()),
	)
	w.Write([]byte(`{"error": "` + error + `"}`))
}
