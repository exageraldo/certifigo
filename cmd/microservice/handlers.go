package main

import (
	"encoding/json"
	"net/http"

	"github.com/exageraldo/certifigo"
)

type generateBody struct {
	Event     certifigo.Event      `json:"event"`
	Speakers  []certifigo.Speaker  `json:"speakers"`
	Attendees []certifigo.Attendee `json:"attendees"`

	Certificate struct {
		CanvaSize  certifigo.WxHSize          `json:"certification_size"`
		Background certifigo.BackgroundConfig `json:"background"`
		Text       certifigo.TextConfig       `json:"text"`
		Validator  certifigo.ValidatorConfig  `json:"validator"`
		Signature  certifigo.SignatureConfig  `json:"signature"`
	} `json:"certificate"`
}

type validateBody struct {
	EventName   string `json:"event_name"`
	PersonName  string `json:"person_name"`
	PersonEmail string `json:"person_email"`
	Hash        string `json:"hash"`
}

func (app *application) defaultConfigCertificatesGet(w http.ResponseWriter, r *http.Request) {
	certCfg, err := cfg.getCertificateConfigFile(certifigo.Event{
		Name:     "<EVENT NAME>",
		Location: "<EVENT LOCATION>",
		Date:     certifigo.StringDate("01/01/2001"),
		Duration: 0,

		Signature: "<SIGNATURE NAME>",
	})
	if err != nil {
		app.logger.Error("Error loading certificate config file", "error", err)
		app.httpError(w, r, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(certCfg)
	if err != nil {
		app.logger.Error("Error marshalling certificate config file to JSON", "error", err)
		app.httpError(w, r, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (app *application) healthcheckGet(w http.ResponseWriter, r *http.Request) {
	const (
		availableStatus   = "available"
		unavailableStatus = "unavailable"
	)

	statusMap := map[string]string{
		"service": availableStatus,
	}

	// Check if the database connection is alive.
	// statusMap["db"] = availableStatus
	// if err := app.db.Ping(); err != nil {
	// 	statusMap["db"] = unavailableStatus
	// }

	// Check if the Email credentials are defined.
	statusMap["email"] = availableStatus
	if !cfg.credentials.CheckEmailCredentials() {
		statusMap["email"] = unavailableStatus
	}

	// Convert the map to JSON and write it to the response.
	jsonResponse, err := json.Marshal(statusMap)
	if err != nil {
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (app *application) generateCertificatesPost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) validateCertificatesPost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
