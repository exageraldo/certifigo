package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/exageraldo/suacuna-cli/assets"
	"github.com/spf13/viper"
)

type CertificateType string

const (
	AttendanceCertification CertificateType = "ATTENDEE"
	SpeakerCertification    CertificateType = "SPEAKER"
)

type ColorConfig struct {
	R, G, B, A float64
}

type SizeConfig struct {
	W, H int
}

type Certificate struct {
	// General
	CanvaSize         SizeConfig
	BackgroundColor   ColorConfig
	OverlayMarginSize float64
	OverlayColor      ColorConfig

	// fonts
	FontsDir string

	// text
	TextSize  float64
	TextColor ColorConfig

	// title
	AttendanceTitle string
	SpeakerTitle    string
	TitleTextSize   float64
	TitleTextColor  ColorConfig

	// person
	PersonTextSize float64

	// valdiator id
	ValidatorMinLength int
	ValidatorMaxLength int
	ValidatorTextSize  float64
	ValidatorTextColor ColorConfig

	// signature
	SignatureDir        string
	SignatureLineLength int
	SignatureImgSize    int
	SignatureTextSize   float64
	SignatureTextColor  ColorConfig
	SignatureTitleSize  float64
	SignatureTitleColor ColorConfig

	// output
	OutputDir       string
	DefaultFileName string
}

func (c Certificate) GetTitleFromType(t CertificateType) (string, error) {
	switch t {
	case AttendanceCertification:
		return c.AttendanceTitle, nil
	case SpeakerCertification:
		return c.SpeakerTitle, nil
	default:
		return "", ErrInvalidCertificateType
	}
}

func (c Certificate) MountSignaturePath(signature string) (string, error) {
	path, err := filepath.Abs(filepath.Join(c.SignatureDir, signature))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return "", ErrInvalidPath
	}
	return path, nil
}

func (c Certificate) MountFontPath(font string) (string, error) {
	path, err := filepath.Abs(filepath.Join(c.SignatureDir, font))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return "", ErrInvalidPath
	}
	return path, nil
}

func (c Certificate) MountOutputPath(out string) (string, error) {
	path, err := filepath.Abs(filepath.Join(c.OutputDir, out))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return "", ErrInvalidPath
	}
	return path, nil
}

func NewCertificateFromConfig() (*Certificate, error) {
	cfg := &Certificate{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadCertificateDefaults() error {
	viper.SetConfigType("toml")

	cfg, err := assets.LoadConfig("certificate.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't load default config: %s\n", err)
		return err
	}

	if err := viper.ReadConfig(bytes.NewReader(cfg)); err != nil {
		fmt.Fprintf(os.Stderr, "Can't read default config: %s\n", err)
		return err
	}

	return nil
}

func LoadCertificateFromFile(filename string) error {
	viper.SetConfigType("toml")

	_, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "Config file does not exist: %s\n", err)
			return ErrInvalidPath
		}
		fmt.Fprintf(os.Stderr, "Can't read passed config: %s\n", err)
		return err
	}

	viper.SetConfigFile(filename)
	if err := viper.MergeInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Can't merge config: %s\n", err)
		return err
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Can't read passed config: %s\n", err)
		return err
	}

	return nil
}
