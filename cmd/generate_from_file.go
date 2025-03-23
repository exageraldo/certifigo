package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/exageraldo/suacuna-cli/certificates"
	"github.com/exageraldo/suacuna-cli/config"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var (
	ErrConfigFileRequired = errors.New("config file is required")
)

func init() {
	generateFromFileCmd.Flags().String("file", "", "Configuration file")
	generateFromFileCmd.MarkFlagRequired("file")

	generateCmd.AddCommand(generateFromFileCmd)
}

type eventConfig struct {
	Event struct {
		Name      string         `toml:"name"`
		Location  string         `toml:"location"`
		Date      toml.LocalDate `toml:"date"`
		Duration  int            `toml:"duration"`
		Signature string         `toml:"signature"`
		Folder    string         `toml:"folder"`
	} `toml:"event"`

	Speakers []struct {
		Name         string `toml:"name"`
		Email        string `toml:"email"`
		TalkTitle    string `toml:"talkTitle"`
		TalkDuration int    `toml:"talkDuration"`
		Attendee     bool   `toml:"attendee"`
		Notify       bool   `toml:"notify"`
	} `toml:"speakers"`

	Attendees []struct {
		Name  string `toml:"name"`
		Email string `toml:"email"`
	} `toml:"attendees"`
}

func (c eventConfig) GetEvent() *certificates.Event {
	event := &certificates.Event{
		Name:     c.Event.Name,
		Loc:      c.Event.Location,
		Date:     c.Event.Date.AsTime(time.Local),
		Duration: c.Event.Duration,
	}
	return event
}

func (c eventConfig) GetAttendees() []*certificates.Attendee {
	attendees := make([]*certificates.Attendee, 0, len(c.Attendees))
	for _, a := range c.Attendees {
		attendees = append(attendees, &certificates.Attendee{
			Name:  a.Name,
			Email: a.Email,
		})

	}
	return attendees
}

func (c eventConfig) GetSpeakers() []*certificates.Speaker {
	speakers := make([]*certificates.Speaker, 0, len(c.Speakers))
	for _, s := range c.Speakers {
		speakers = append(speakers, &certificates.Speaker{
			Name:         s.Name,
			Email:        s.Email,
			TalkTitle:    s.TalkTitle,
			TalkDuration: s.TalkDuration,
		})
	}
	return speakers
}

func getFileContentFromCmd(cmd *cobra.Command) ([]byte, error) {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, err
	}
	if filePath == "" {
		return nil, ErrConfigFileRequired
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

var generateFromFileCmd = &cobra.Command{
	Use:   "from-file",
	Short: "Generate certificates from a configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		fileContent, err := getFileContentFromCmd(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		var eventCfg eventConfig
		err = toml.Unmarshal(fileContent, &eventCfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		certCfg, err := config.NewCertificateFromConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		event := eventCfg.GetEvent()
		for _, s := range eventCfg.GetSpeakers() {
			c := certificates.NewSpeakerCertificate(
				*s,
				*event,
				eventCfg.Event.Signature,
				"",
				*certCfg,
			)
			if err := c.Generate(); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		}

		for _, a := range eventCfg.GetAttendees() {
			c := certificates.NewAttendanceCertificate(
				*a,
				*event,
				eventCfg.Event.Signature,
				"",
				*certCfg,
			)
			if err := c.Generate(); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		}
	},
}
