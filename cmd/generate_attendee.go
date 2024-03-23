package cmd

import (
	"fmt"
	"os"

	"github.com/exageraldo/suacuna-cli/certificates"
	"github.com/spf13/cobra"
)

func init() {
	generateCmd.AddCommand(generateAttendeeCmd)
}

var generateAttendeeCmd = &cobra.Command{
	Use:   "attendee",
	Short: "Generate certificates for attendees.",
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}
