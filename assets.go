package certifigo

import (
	"embed"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	// Font names
	OpenSans          string = "open-sans"
	CedarvilleCursive string = "cedarville-cursive"
)

var (
	//go:embed _assets/configs/*.toml
	//go:embed _assets/fonts/*.ttf
	assetsDir embed.FS

	embededFonts = map[string]string{
		OpenSans:          "_assets/fonts/OpenSans-Bold.ttf",
		CedarvilleCursive: "_assets/fonts/CedarvilleCursive-Regular.ttf",
	}
)

func LoadDefaultCertificateConfigFile(data map[string]any) (*CertificateConfigFile, error) {
	var certFile CertificateConfigFile
	if err := ParseTOMLTemplateFS(
		"_assets/configs/default_certificate.toml",
		&certFile,
		data,
	); err != nil {
		return nil, err
	}

	return &certFile, nil
}

func LoadCertificateConfigFile(filePath string, data map[string]any) (*CertificateConfigFile, error) {
	var certFile CertificateConfigFile
	if err := ParseTOMLTemplate(filePath, &certFile, data); err != nil {
		return nil, err
	}
	return &certFile, nil
}

func LoadEventFile(filePath string) (*EventFile, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var eventFile EventFile
	if err := ParseTOMLFile(fileContent, &eventFile); err != nil {
		return nil, err
	}
	return &eventFile, nil
}

func LoadFont(fontPath string, size float64) (font.Face, error) {
	var fileContent []byte
	var err error
	if embFontPath, ok := embededFonts[fontPath]; ok {
		fileContent, err = assetsDir.ReadFile(embFontPath)
	} else {
		fileContent, err = os.ReadFile(fontPath)
	}

	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(fileContent)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(f, &truetype.Options{
		Size: size,
	}), nil
}
