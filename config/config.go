package config

import (
	_ "embed"
	"errors"
)

var (
	ErrConfigFileRequired     = errors.New("config file is required")
	ErrInvalidPath            = errors.New("invalid path")
	ErrInvalidCertificateType = errors.New("invalid certification type")
)
