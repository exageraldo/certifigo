package fonts

import (
	_ "embed"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	//go:embed _assets/OpenSans-Bold.ttf
	DefaultFont []byte

	//go:embed _assets/CedarvilleCursive-Regular.ttf
	SignatureFont []byte
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
