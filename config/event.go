package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Event struct {
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
		Name   string `toml:"name"`
		Email  string `toml:"email"`
		Notify bool   `toml:"notify"`
	} `toml:"attendees"`
}

func NewEventFromFile(path string) (*Event, error) {
	cfgFilePath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return nil, ErrInvalidPath
	}
	_, err = os.Stat(cfgFilePath)
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "Config file does not exist: %v\n", err)
		return nil, ErrInvalidPath
	}

	fileContent, err := os.ReadFile(cfgFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var cfg Event
	err = toml.Unmarshal(fileContent, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &cfg, nil
}
