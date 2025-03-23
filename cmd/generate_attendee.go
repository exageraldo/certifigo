package cmd

import (
	"fmt"
	"os"

	"github.com/exageraldo/suacuna-cli/certificates"
	"github.com/exageraldo/suacuna-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	setAttendeeFlags(generateAttendeeCmd)
	setEventFlags(generateAttendeeCmd)
	setSignatureFlag(generateAttendeeCmd)
	setNotificationFlag(generateAttendeeCmd)
	setLogoFlag(generateAttendeeCmd)

	generateCmd.AddCommand(generateAttendeeCmd)
}

func setAttendeeFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Name of the attendee")
	cmd.Flags().String("email", "", "Email of the attendee")
	cmd.MarkFlagRequired("name")
}

func attendeeFromCmd(cmd *cobra.Command) (*certificates.Attendee, error) {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return nil, err
	}

	email, err := cmd.Flags().GetString("email")
	if err != nil {
		return nil, err
	}

	return &certificates.Attendee{
		Name:  name,
		Email: email,
	}, nil
}

var generateAttendeeCmd = &cobra.Command{
	Use:   "attendee",
	Short: "Generate certificates for attendees.",
	Run: func(cmd *cobra.Command, args []string) {
		event, err := eventFromCmd(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		attendee, err := attendeeFromCmd(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		signature, err := cmd.Flags().GetString("signature")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		logo, err := cmd.Flags().GetString("logo")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		cfg, err := config.NewCertificateFromConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		c := certificates.NewAttendanceCertificate(
			*attendee,
			*event,
			signature,
			logo,
			*cfg,
		)

		if err := c.Generate(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	},
}
