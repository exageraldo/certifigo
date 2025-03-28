package certifigo

import (
	"path/filepath"
)

type CertificateType string

const (
	AttendanceCertification CertificateType = "ATTENDEE"
	SpeakerCertification    CertificateType = "SPEAKER"
)

type BackgroundConfig struct {
	Color       HexColor `toml:"color"`
	BorderSize  float64  `toml:"border_size"`
	BorderColor HexColor `toml:"border_color"`
}

type TextConfig struct {
	// fonts
	FontsDir string `toml:"fonts_dir"`

	// text
	TextSize  float64  `toml:"text_size"`
	TextColor HexColor `toml:"text_color"`

	// title
	TitleTextSize  float64  `toml:"title_text_size"`
	TitleTextColor HexColor `toml:"title_text_color"`

	// person
	PersonTextSize float64 `toml:"person_text_size"`
}

type ValidatorConfig struct {
	MinLength int      `toml:"min_length"`
	MaxLength int      `toml:"max_length"`
	TextSize  float64  `toml:"text_size"`
	TextColor HexColor `toml:"text_color"`
}

type SignatureConfig struct {
	LineLength int      `toml:"line_length"`
	ImgSize    int      `toml:"img_size"`
	TextSize   float64  `toml:"text_size"`
	TextColor  HexColor `toml:"text_color"`
	TitleSize  float64  `toml:"title_size"`
	TitleColor HexColor `toml:"title_color"`
	Folder     string   `toml:"folder"`
}

type OutputConfig struct {
	Folder          string `toml:"folder"`
	DefaultFileName string `toml:"default_file_name"`
}

type TemplateConfig struct {
	Title        string `toml:"title"`
	Body         string `toml:"body"`
	EmailSubject string `toml:"email_subject"`
	EmailBody    string `toml:"email_body"`
}

type CertificateConfigFile struct {
	CanvaSize WxHSize `toml:"certification_size"`

	Background BackgroundConfig `toml:"background"`
	Text       TextConfig       `toml:"text"`
	Validator  ValidatorConfig  `toml:"validator"`
	Signature  SignatureConfig  `toml:"signature"`
	Output     OutputConfig     `toml:"output"`

	Attendee TemplateConfig `toml:"attendee"`
	Speaker  TemplateConfig `toml:"speaker"`
}

func (c CertificateConfigFile) MountOutputPath(out string) (string, error) {
	path, err := filepath.Abs(filepath.Join(c.Output.Folder, out))
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c CertificateConfigFile) MountSignaturePath(signature string) (string, error) {
	path, err := filepath.Abs(filepath.Join(c.Signature.Folder, signature))
	if err != nil {
		return "", err
	}
	return path, nil
}
