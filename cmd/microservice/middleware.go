package main

import (
	"net/http"
	"runtime/debug"
	"strings"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			url    = r.URL.RequestURI()
			reqID  = GetReqID(r.Context())
		)

		app.logger.Info(
			"received request",
			"ip", ip,
			"proto", proto,
			"method", method,
			"url", url,
			"x_reqquest_id", reqID,
		)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		app.logger.Info(
			"response finished",
			"ip", ip,
			"proto", proto,
			"method", method,
			"url", url,
			"x_reqquest_id", reqID,
		)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set(HTTPHeaderNameRequestID, GetReqID(r.Context()))
				app.httpError(w, r, "Internal Server Error", http.StatusInternalServerError)
				app.logger.Error(
					"panic recovered",
					"error", err,
					"x_request_id", GetReqID(r.Context()),
					"stack", string(debug.Stack()),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.logger.Error("Missing Token", "error", authHeader)
			app.httpError(w, r, "Missing Token", http.StatusUnauthorized)
			return
		}
		authToken := strings.Split(authHeader, "Bearer ")
		if len(authToken) != 2 {
			app.logger.Error("Malformed Token", "error", authToken)
			app.httpError(w, r, "Malformed Token", http.StatusUnauthorized)
		} else {
			jwtToken := authToken[1]
			if jwtToken == cfg.credentials.JWTToken {
				next.ServeHTTP(w, r)
			} else {
				app.logger.Error("Invalid Token", "token", jwtToken)
				app.httpError(w, r, "Unauthorized", http.StatusUnauthorized)
			}

		}
	})
}

func (app *application) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := AttachReqID(r.Context())
		r = r.WithContext(ctx)
		w.Header().Set(HTTPHeaderNameRequestID, GetReqID(r.Context()))

		next.ServeHTTP(w, r)
	})
}
