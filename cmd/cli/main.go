package main

import (
	"fmt"
	"os"

	"github.com/exageraldo/suacuna-cli/certificates"
)

func main() {
	GenerateAttendeeCertificate()
	GenerateSpeakerCertificate()
}

func GenerateAttendeeCertificate() {
	event, err := certificates.NewEvent("21º Meetup", "Empresa Tal", "04/02/2024", 4)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	attendee := certificates.Attendee{
		Name:  "Fulana de Tal",
		Email: "",
	}
	c, err := certificates.NewAttendanceCertificate(
		attendee,
		*event,
		"Cicrano de Tal",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := c.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func GenerateSpeakerCertificate() {
	event, err := certificates.NewEvent("21º Meetup", "Empresa Tal", "04/02/2024", 4)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	speaker := certificates.Speaker{
		Name:         "Fulana de Tal",
		Email:        "",
		TalkTitle:    "Como contribuir para projetos open source",
		TalkDuration: 45,
	}
	c, err := certificates.NewSpeakerCertificate(
		speaker,
		*event,
		"Cicrano de Tal",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := c.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
