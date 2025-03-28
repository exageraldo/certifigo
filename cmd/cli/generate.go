package main

import (
	"github.com/exageraldo/certifigo"
	"github.com/spf13/cobra"
)

var (
	EventFromCLI     certifigo.Event
	AttendeeFromCLI  certifigo.Attendee
	SpeakerFromCLI   certifigo.Speaker
	EventFileFromCLI string
)

func init() {
	// attendee subcommand flags
	generateAttendeeCmd.Flags().StringVar(&AttendeeFromCLI.Name, "name", "", "Name of the attendee")
	generateAttendeeCmd.Flags().StringVar(&AttendeeFromCLI.Email, "email", "", "Email of the attendee")
	generateAttendeeCmd.Flags().BoolVar(&AttendeeFromCLI.Notify, "notify", false, "Send email notification")
	generateAttendeeCmd.Flags().StringVar(&EventFromCLI.Name, "event", "", "Name of the event")
	generateAttendeeCmd.Flags().StringVar(&EventFromCLI.Location, "loc", "", "Location of the event")
	generateAttendeeCmd.Flags().StringVar((*string)(&EventFromCLI.Date), "date", "", "Date of the event")
	generateAttendeeCmd.Flags().IntVar(&EventFromCLI.Duration, "duration", 0, "Duration of the event")
	generateAttendeeCmd.Flags().StringVar(&EventFromCLI.Signature, "signature", "", "Name of the signature")
	generateAttendeeCmd.Flags().StringVar(&EventFromCLI.SignatureImg, "signature-img", "", "Signature image path")
	generateAttendeeCmd.Flags().StringVar(&EventFromCLI.Logo, "logo", "", "Logo image path")

	generateAttendeeCmd.MarkFlagRequired("name")
	generateAttendeeCmd.MarkFlagRequired("event")
	generateAttendeeCmd.MarkFlagRequired("loc")
	generateAttendeeCmd.MarkFlagRequired("date")
	generateAttendeeCmd.MarkFlagRequired("duration")
	generateAttendeeCmd.MarkFlagRequired("signature")
	generateCmd.AddCommand(generateAttendeeCmd)

	// speaker subcommand flags
	generateSpeakerCmd.Flags().StringVar(&SpeakerFromCLI.Name, "name", "", "Name of the speaker")
	generateSpeakerCmd.Flags().StringVar(&SpeakerFromCLI.Email, "email", "", "Email of the speaker")
	generateSpeakerCmd.Flags().BoolVar(&SpeakerFromCLI.Notify, "notify", false, "Send email notification")
	generateSpeakerCmd.Flags().StringVar(&SpeakerFromCLI.TalkTitle, "talk-title", "", "Title of the talk")
	generateSpeakerCmd.Flags().IntVar(&SpeakerFromCLI.TalkDuration, "talk-duration", 0, "Duration of the talk")
	generateSpeakerCmd.Flags().BoolVar(&SpeakerFromCLI.Attendee, "attendee", false, "")
	generateSpeakerCmd.Flags().StringVar(&EventFromCLI.Name, "event", "", "Name of the event")
	generateSpeakerCmd.Flags().StringVar(&EventFromCLI.Location, "loc", "", "Location of the event")
	generateSpeakerCmd.Flags().StringVar((*string)(&EventFromCLI.Date), "date", "", "Date of the event")
	generateSpeakerCmd.Flags().IntVar(&EventFromCLI.Duration, "duration", 0, "Duration of the event")
	generateSpeakerCmd.Flags().StringVar(&EventFromCLI.Signature, "signature", "", "Name of the signature")
	generateSpeakerCmd.Flags().StringVar(&EventFromCLI.Logo, "logo", "", "Logo image path")

	generateSpeakerCmd.MarkFlagRequired("name")
	generateSpeakerCmd.MarkFlagRequired("talk-title")
	generateSpeakerCmd.MarkFlagRequired("talk-duration")
	generateSpeakerCmd.MarkFlagRequired("event")
	generateSpeakerCmd.MarkFlagRequired("loc")
	generateSpeakerCmd.MarkFlagRequired("date")
	generateSpeakerCmd.MarkFlagRequired("duration")
	generateSpeakerCmd.MarkFlagRequired("signature")
	generateCmd.AddCommand(generateSpeakerCmd)

	// from-file subcommand flags
	generateFromFileCmd.Flags().StringVar(&EventFileFromCLI, "file", "", "Event file")

	generateFromFileCmd.MarkFlagRequired("file")
	generateCmd.AddCommand(generateFromFileCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate certificates for events.",
}

var generateAttendeeCmd = &cobra.Command{
	Use:   "attendee",
	Short: "Generate certificates for attendees.",
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := certifigo.NewEnvCredentials()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		if AttendeeFromCLI.Notify && !credentials.CheckEmailCredentials() {
			cmd.PrintErr("Email credentials not set.\n")
			return
		}

		var certificateConfigFile certifigo.CertificateConfigFile
		defaultCfgFile, err := certifigo.LoadDefaultCertificateConfigFile(
			map[string]any{"Event": EventFromCLI},
		)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		if ConfigFileFromCLI != "" {
			cliCfgFile, err := certifigo.LoadCertificateConfigFile(
				ConfigFileFromCLI,
				map[string]any{
					"Event":  EventFromCLI,
					"Config": defaultCfgFile,
				},
			)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			certificateConfigFile = certifigo.Merge(*defaultCfgFile, *cliCfgFile)
		} else {
			certificateConfigFile = *defaultCfgFile
		}

		certPath, err := certifigo.NewCertificateDrawer(
			certifigo.AttendanceCertification,
			EventFromCLI,
			certificateConfigFile,
		).DrawAndSave(AttendeeFromCLI.Name)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if AttendeeFromCLI.Notify {
			sender, err := certifigo.NewGMailSender(
				credentials.EmailSender,
				credentials.EmailPassword,
			)
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			email := certifigo.Email{
				Subject:     certificateConfigFile.Attendee.EmailSubject,
				Body:        certificateConfigFile.Attendee.EmailBody,
				To:          AttendeeFromCLI.Email,
				Attachments: []string{certPath},
			}
			if err := sender.Send(email); err != nil {
				cmd.PrintErr(err)
				return
			}
		}
	},
}

var generateSpeakerCmd = &cobra.Command{
	Use:   "speaker",
	Short: "Generate certificates for speakers.",
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := certifigo.NewEnvCredentials()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if SpeakerFromCLI.Notify && !credentials.CheckEmailCredentials() {
			cmd.PrintErr("Email credentials not set.")
			return
		}

		var certificateConfigFile certifigo.CertificateConfigFile
		defaultCfgFile, err := certifigo.LoadDefaultCertificateConfigFile(
			map[string]any{"Event": EventFromCLI},
		)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		if ConfigFileFromCLI != "" {
			cliCfgFile, err := certifigo.LoadCertificateConfigFile(
				ConfigFileFromCLI,
				map[string]any{
					"Event":  EventFromCLI,
					"Config": defaultCfgFile,
				},
			)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			certificateConfigFile = certifigo.Merge(*defaultCfgFile, *cliCfgFile)
		} else {
			certificateConfigFile = *defaultCfgFile
		}

		sCertPath, err := certifigo.NewCertificateDrawer(
			certifigo.SpeakerCertification,
			EventFromCLI,
			certificateConfigFile,
		).DrawAndSave(SpeakerFromCLI.Name)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		certificationsPath := []string{sCertPath}
		if SpeakerFromCLI.Attendee {
			aCertPath, err := certifigo.NewCertificateDrawer(
				certifigo.AttendanceCertification,
				EventFromCLI,
				certificateConfigFile,
			).DrawAndSave(SpeakerFromCLI.Name)
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			certificationsPath = append(certificationsPath, aCertPath)
		}

		if SpeakerFromCLI.Notify {
			sender, err := certifigo.NewGMailSender(
				credentials.EmailSender,
				credentials.EmailPassword,
			)
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			email := certifigo.Email{
				Subject:     certificateConfigFile.Speaker.EmailSubject,
				Body:        certificateConfigFile.Speaker.EmailBody,
				To:          SpeakerFromCLI.Email,
				Attachments: certificationsPath,
			}
			if err := sender.Send(email); err != nil {
				cmd.PrintErr(err)
				return
			}
		}
	},
}

var generateFromFileCmd = &cobra.Command{
	Use:   "from-file",
	Short: "Generate certificates from a configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		eventFile, err := certifigo.LoadEventFile(EventFileFromCLI)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		var certificateConfigFile certifigo.CertificateConfigFile
		defaultCfgFile, err := certifigo.LoadDefaultCertificateConfigFile(
			map[string]any{"Event": eventFile.Event},
		)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		if ConfigFileFromCLI != "" {
			cliCfgFile, err := certifigo.LoadCertificateConfigFile(
				ConfigFileFromCLI,
				map[string]any{
					"Event":  eventFile.Event,
					"Config": defaultCfgFile,
				},
			)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			certificateConfigFile = certifigo.Merge(*defaultCfgFile, *cliCfgFile)
		} else {
			certificateConfigFile = *defaultCfgFile
		}

		var emails []certifigo.Email
		for _, attendee := range eventFile.Attendees {
			certPath, err := certifigo.NewCertificateDrawer(
				certifigo.AttendanceCertification,
				eventFile.Event,
				certificateConfigFile,
			).DrawAndSave(attendee.Name)
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			if attendee.Notify {
				emails = append(emails, certifigo.Email{
					Subject:     certificateConfigFile.Attendee.EmailSubject,
					Body:        certificateConfigFile.Attendee.EmailBody,
					To:          attendee.Email,
					Attachments: []string{certPath},
				})
			}
		}

		for _, speaker := range eventFile.Speakers {
			sCertPath, err := certifigo.NewCertificateDrawer(
				certifigo.SpeakerCertification,
				eventFile.Event,
				certificateConfigFile,
			).DrawAndSave(speaker.Name)
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			certificationsPath := []string{sCertPath}
			if speaker.Attendee {
				aCertPath, err := certifigo.NewCertificateDrawer(
					certifigo.AttendanceCertification,
					eventFile.Event,
					certificateConfigFile,
				).DrawAndSave(speaker.Name)
				if err != nil {
					cmd.PrintErr(err)
					return
				}

				certificationsPath = append(certificationsPath, aCertPath)
			}

			if speaker.Notify {
				emails = append(emails, certifigo.Email{
					Subject:     certificateConfigFile.Speaker.EmailSubject,
					Body:        certificateConfigFile.Speaker.EmailBody,
					To:          speaker.Email,
					Attachments: certificationsPath,
				})
			}
		}

		if len(emails) == 0 {
			return
		}

		credentials, err := certifigo.NewEnvCredentials()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if len(emails) > 0 && !credentials.CheckEmailCredentials() {
			cmd.PrintErr("All certificates were generated, but no email was sent because the email credentials were not set.")
			return
		}

		sender, err := certifigo.NewGMailSender(
			credentials.EmailSender,
			credentials.EmailPassword,
		)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if err := sender.BulkSend(emails); err != nil {
			cmd.PrintErr(err)
			return
		}
	},
}
