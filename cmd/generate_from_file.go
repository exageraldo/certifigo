package cmd

import (
	"fmt"
	"os"

	"github.com/exageraldo/suacuna-cli/certificates"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func init() {
	generateFromFileCmd.Flags().String("file", "", "Configuration file")
	generateFromFileCmd.MarkFlagRequired("file")

	generateCmd.AddCommand(generateFromFileCmd)
}

type tomlFileCfg struct {
	Event struct {
		Name      string `toml:"name"`
		Location  string `toml:"location"`
		Date      string `toml:"date"`
		Duration  int    `toml:"duration"`
		Signature string `toml:"signature"`
	} `toml:"event"`

	Speakers []struct {
		Name         string `toml:"name"`
		Email        string `toml:"email"`
		TalkTitle    string `toml:"talkTitle"`
		TalkDuration int    `toml:"talkDuration"`
	} `toml:"speakers"`

	Attendees []struct {
		Name  string `toml:"name"`
		Email string `toml:"email"`
	} `toml:"attendees"`
}

func (c tomlFileCfg) GetEvent() (*certificates.Event, error) {
	event, err := certificates.NewEvent(c.Event.Name, c.Event.Location, c.Event.Date, c.Event.Duration)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (c tomlFileCfg) GetAttendees() ([]*certificates.Attendee, error) {
	attendees := make([]*certificates.Attendee, 0, len(c.Attendees))
	for _, a := range c.Attendees {
		attendees = append(attendees, &certificates.Attendee{
			Name:  a.Name,
			Email: a.Email,
		})

	}
	return attendees, nil
}

func (c tomlFileCfg) GetSpeakers() ([]*certificates.Speaker, error) {
	speakers := make([]*certificates.Speaker, 0, len(c.Speakers))
	for _, s := range c.Speakers {
		speakers = append(speakers, &certificates.Speaker{
			Name:         s.Name,
			Email:        s.Email,
			TalkTitle:    s.TalkTitle,
			TalkDuration: s.TalkDuration,
		})
	}
	return speakers, nil
}

var generateFromFileCmd = &cobra.Command{
	Use:   "from-file",
	Short: "Generate certificates from a configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if filePath == "" {
			fmt.Fprintf(os.Stderr, "Config file is required\n")
			os.Exit(1)
		}

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		var cfg tomlFileCfg
		err = toml.Unmarshal(fileContent, &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		event, err := cfg.GetEvent()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		speakers, err := cfg.GetSpeakers()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		attendees, err := cfg.GetAttendees()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		for _, a := range attendees {
			c, err := certificates.NewAttendanceCertificate(
				*a,
				*event,
				cfg.Event.Signature,
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

		for _, s := range speakers {
			c, err := certificates.NewSpeakerCertificate(
				*s,
				*event,
				cfg.Event.Signature,
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
