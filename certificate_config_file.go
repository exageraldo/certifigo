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
	Color       HexColor `toml:"color" json:"color"`
	BorderSize  float64  `toml:"border_size" json:"border_size"`
	BorderColor HexColor `toml:"border_color" json:"border_color"`
}

type TextConfig struct {
	// fonts
	FontsDir string `toml:"fonts_dir" json:"-"`

	// text
	TextSize  float64  `toml:"text_size" json:"text_size"`
	TextColor HexColor `toml:"text_color" json:"text_color"`

	// title
	TitleTextSize  float64  `toml:"title_text_size" json:"title_text_size"`
	TitleTextColor HexColor `toml:"title_text_color" json:"title_text_color"`

	// person
	PersonTextSize float64 `toml:"person_text_size" json:"person_text_size"`
}

type ValidatorConfig struct {
	MinLength int      `toml:"min_length" json:"min_length"`
	MaxLength int      `toml:"max_length" json:"max_length"`
	TextSize  float64  `toml:"text_size" json:"text_size"`
	TextColor HexColor `toml:"text_color" json:"text_color"`
}

type SignatureConfig struct {
	LineLength int      `toml:"line_length" json:"line_length"`
	ImgSize    int      `toml:"img_size" json:"img_size"`
	TextSize   float64  `toml:"text_size" json:"text_size"`
	TextColor  HexColor `toml:"text_color" json:"text_color"`
	TitleSize  float64  `toml:"title_size" json:"title_size"`
	TitleColor HexColor `toml:"title_color" json:"title_color"`
	Folder     string   `toml:"folder" json:"-"`
}

type OutputConfig struct {
	Folder          string `toml:"folder" json:"folder"`
	DefaultFileName string `toml:"default_file_name" json:"default_file_name"`
}

type TemplateConfig struct {
	Title        string `toml:"title" json:"title"`
	Body         string `toml:"body" json:"body"`
	EmailSubject string `toml:"email_subject" json:"email_subject"`
	EmailBody    string `toml:"email_body" json:"email_body"`
}

type CertificateConfigFile struct {
	CanvaSize WxHSize `toml:"certification_size" json:"certification_size"`

	Background BackgroundConfig `toml:"background" json:"background"`
	Text       TextConfig       `toml:"text" json:"text"`
	Validator  ValidatorConfig  `toml:"validator" json:"validator"`
	Signature  SignatureConfig  `toml:"signature" json:"signature"`
	Output     OutputConfig     `toml:"output" json:"-"`

	Attendee TemplateConfig `toml:"attendee" json:"attendee"`
	Speaker  TemplateConfig `toml:"speaker" json:"speaker"`
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
