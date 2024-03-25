package cmd

import (
	"github.com/exageraldo/suacuna-cli/certificates"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

func setEventFlags(cmd *cobra.Command) {
	cmd.Flags().String("event", "", "Name of the event")
	cmd.Flags().String("loc", "", "Location of the event")
	cmd.Flags().String("date", "", "Date of the event")
	cmd.Flags().Int("duration", 0, "Duration of the event")

	cmd.MarkFlagRequired("event")
	cmd.MarkFlagRequired("loc")
	cmd.MarkFlagRequired("date")
	cmd.MarkFlagRequired("duration")
}

func setSignatureFlag(cmd *cobra.Command) {
	cmd.Flags().String("signature", "", "Name of the signature")
	cmd.MarkFlagRequired("signature")
}

func setNotificationFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("notify", false, "Send email notification")
}

func eventFromCmd(cmd *cobra.Command) (*certificates.Event, error) {
	event, err := cmd.Flags().GetString("event")
	if err != nil {
		return nil, err
	}
	loc, err := cmd.Flags().GetString("loc")
	if err != nil {
		return nil, err
	}
	date, err := cmd.Flags().GetString("date")
	if err != nil {
		return nil, err
	}
	duration, err := cmd.Flags().GetInt("duration")
	if err != nil {
		return nil, err
	}

	return certificates.NewEvent(event, loc, date, duration)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate certificates for events.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
