package cmd

import (
	"fmt"
	"os"

	"github.com/exageraldo/suacuna-cli/certificates"
	"github.com/exageraldo/suacuna-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	setSpeakerFlags(generateSpeakerCmd)
	setEventFlags(generateSpeakerCmd)
	setSignatureFlag(generateSpeakerCmd)
	setNotificationFlag(generateSpeakerCmd)
	setLogoFlag(generateSpeakerCmd)

	// generate attendee certificate
	generateSpeakerCmd.Flags().Bool("attendee", false, "Generate attendee certificate")

	generateCmd.AddCommand(generateSpeakerCmd)
}

func setSpeakerFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Name of the speaker")
	cmd.Flags().String("email", "", "Email of the speaker")
	cmd.Flags().String("talk-title", "", "Title of the talk")
	cmd.Flags().Int("talk-duration", 0, "Duration of the talk")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("talk-title")
	cmd.MarkFlagRequired("talk-duration")
}

func speakerFromCmd(cmd *cobra.Command) (*certificates.Speaker, error) {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return nil, err
	}
	email, err := cmd.Flags().GetString("email")
	if err != nil {
		return nil, err
	}
	talk, err := cmd.Flags().GetString("talk-title")
	if err != nil {
		return nil, err
	}
	talkDuration, err := cmd.Flags().GetInt("talk-duration")
	if err != nil {
		return nil, err
	}

	return &certificates.Speaker{
		Name:         name,
		Email:        email,
		TalkTitle:    talk,
		TalkDuration: talkDuration,
	}, nil
}

var generateSpeakerCmd = &cobra.Command{
	Use:   "speaker",
	Short: "Generate certificates for speakers.",
	Run: func(cmd *cobra.Command, args []string) {
		event, err := eventFromCmd(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		speaker, err := speakerFromCmd(cmd)
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
		c := certificates.NewSpeakerCertificate(
			*speaker,
			*event,
			signature,
			logo,
			*cfg,
		)

		if err := c.Generate(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		generateAttendee, err := cmd.Flags().GetBool("attendee")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if generateAttendee {
			attendee, err := attendeeFromCmd(cmd)
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
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}

			if err := c.Generate(); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		}
	},
}
