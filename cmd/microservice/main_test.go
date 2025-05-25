package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/exageraldo/certifigo"
)

func testApplication() *application {
	cfg = config{
		port:              4000,
		dsn:               "",
		configFileFromCLI: "",
		credentials: certifigo.EnvCredentials{
			EmailSender:   "email@test.com",
			EmailPassword: "some-password-here",
			JWTToken:      "token",
		},
	}

	app := &application{
		logger: slog.New(slog.NewJSONHandler(
			io.Discard,
			&slog.HandlerOptions{},
		)),
		validationModel: &ValidationRecordMock{
			DB: make(map[string]*ValidationRecord),
		},
	}

	return app
}

func Test_ProtectedRoutes(t *testing.T) {
	cases := []struct {
		name   string
		url    string
		method string
	}{
		{"Healthcheck", "/healthcheck", http.MethodGet},
		{"CertificatesGenerate", "/certificates/generate", http.MethodPost},
		{"CertificatesValidate", "/certificates/validate", http.MethodPost},
	}

	app := testApplication()
	handler := app.routes()

	for _, tt := range cases {
		t.Run("MissingToken/"+tt.name, func(t *testing.T) {
			// Arrange
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, nil)
			req.Header.Add("Content-Type", "application/json")

			// Act
			handler.ServeHTTP(resp, req)

			// Assert
			if resp.Result().StatusCode != http.StatusUnauthorized {
				t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.Result().StatusCode)
			}
			if resp.Body.String() != `{"error": "Missing Token"}` {
				t.Errorf("expected error message 'Missing Token', got '%s'", resp.Body.String())
			}
		})
	}

	for _, tt := range cases {
		t.Run("MalformedToken/"+tt.name, func(t *testing.T) {
			// Arrange
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, nil)
			req.Header.Add("Authorization", "Bearer")
			req.Header.Add("Content-Type", "application/json")

			// Act
			handler.ServeHTTP(resp, req)

			// Assert
			if resp.Result().StatusCode != http.StatusUnauthorized {
				t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.Result().StatusCode)
			}
			if resp.Body.String() != `{"error": "Malformed Token"}` {
				t.Errorf("expected error message 'Malformed Token', got '%s'", resp.Body.String())
			}
		})
	}

	for _, tt := range cases {
		t.Run("WrongToken/"+tt.name, func(t *testing.T) {
			// Arrange
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, nil)
			req.Header.Add("Authorization", "Bearer wrong-token")
			req.Header.Add("Content-Type", "application/json")

			// Act
			handler.ServeHTTP(resp, req)

			// Assert
			if resp.Result().StatusCode != http.StatusUnauthorized {
				t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.Result().StatusCode)
			}
			if resp.Body.String() != `{"error": "Unauthorized"}` {
				t.Errorf("expected error message 'Unauthorized', got '%s'", resp.Body.String())
			}
		})
	}
}

func Test_PingGet(t *testing.T) {
	// Arrange
	app := testApplication()
	handler := app.routes()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	resp := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(resp, req)

	// Assert
	if resp.Result().StatusCode != http.StatusNoContent {
		t.Errorf("expected status code %d, got %d", http.StatusNoContent, resp.Result().StatusCode)
	}
	if resp.Body.String() != "" {
		t.Errorf("expected empty body, got '%s'", resp.Body.String())
	}
}

func Test_HealthcheckGet(t *testing.T) {
	cases := []struct {
		name           string
		cfgCredentials certifigo.EnvCredentials
		expectedBody   string
	}{
		{
			"AllAvailable",
			certifigo.EnvCredentials{
				EmailSender:   "email@test.com",
				EmailPassword: "some-password-here",
				JWTToken:      "token",
			},
			`{"email":"available","service":"available"}`,
		},
		{
			"EmailUnavailable/EmptyEmailSender",
			certifigo.EnvCredentials{
				EmailSender:   "",
				EmailPassword: "some-password-here",
				JWTToken:      "token",
			},
			`{"email":"unavailable","service":"available"}`,
		},
		{
			"EmailUnavailable/EmptyEmailPassword",
			certifigo.EnvCredentials{
				EmailSender:   "email@test.com",
				EmailPassword: "",
				JWTToken:      "token",
			},
			`{"email":"unavailable","service":"available"}`,
		},
		{
			"EmailUnavailable/EmptyEmailCredentials",
			certifigo.EnvCredentials{
				EmailSender:   "",
				EmailPassword: "",
				JWTToken:      "token",
			},
			`{"email":"unavailable","service":"available"}`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			app := testApplication()
			handler := app.routes()

			cfg.credentials = tt.cfgCredentials

			resp := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet,
				"/healthcheck",
				nil,
			)
			req.Header.Add("Authorization", "Bearer token")
			req.Header.Add("Content-Type", "application/json")

			handler.ServeHTTP(resp, req)

			if resp.Result().StatusCode != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.Result().StatusCode)
			}
			if resp.Body.String() != tt.expectedBody {
				t.Errorf("expected body '%s', got '%s'", tt.expectedBody, resp.Body.String())
			}
		})
	}
}

func Test_ValidateCertificatesPost(t *testing.T) {
	app := testApplication()
	handler := app.routes()

	{
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/certificates/validate",
			nil,
		)
		req.Header.Add("Authorization", "Bearer token")
		req.Header.Add("Content-Type", "application/json")

		handler.ServeHTTP(resp, req)

		if resp.Result().StatusCode != http.StatusNoContent {
			t.Errorf("expected status code %d, got %d", http.StatusNoContent, resp.Result().StatusCode)
		}
	}
}
