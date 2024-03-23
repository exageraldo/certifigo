package cmd

import (
	"fmt"
	"os"

	"github.com/exageraldo/suacuna-cli/certificates"
	"github.com/spf13/cobra"
)

func init() {
	generateCmd.AddCommand(generateSpeakerCmd)
}

var generateSpeakerCmd = &cobra.Command{
	Use:   "speaker",
	Short: "Generate certificates for speakers.",
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}
