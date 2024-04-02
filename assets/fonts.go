package assets

import (
	"embed"
	"errors"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type FontName string

const (
	OpenSans          FontName = "open-sans"
	CedarvilleCursive FontName = "cedarville-cursive"
)

var (
	//go:embed _fonts/*.ttf
	fontsDir embed.FS

	embededFonts = map[FontName]string{
		OpenSans:          "_fonts/OpenSans-Bold.ttf",
		CedarvilleCursive: "_fonts/CedarvilleCursive-Regular.ttf",
	}

	ErrEmbededFontNotFound = errors.New("embeded font not found")
)

func LoadFont(ttf []byte, size float64) (font.Face, error) {
	f, err := truetype.Parse(ttf)
	if err != nil {
		return nil, err
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size: size,
	})
	return face, nil
}

func LoadEmbededFont(font FontName, size float64) (font.Face, error) {
	fontPath, ok := embededFonts[font]
	if !ok {
		return nil, ErrEmbededFontNotFound
	}
	fileContent, err := fontsDir.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	return LoadFont(fileContent, size)

}
