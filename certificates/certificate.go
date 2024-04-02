package certificates

import (
	"errors"
	"fmt"
	"image/color"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/exageraldo/suacuna-cli/assets"
	"github.com/exageraldo/suacuna-cli/config"
	"github.com/fogleman/gg"
)

func NewEvent(name, loc, date string, duration int) (*Event, error) {
	t, err := time.Parse("02/01/2006", date)
	if err != nil {
		return nil, err
	}

	return &Event{
		Name:     name,
		Loc:      loc,
		Date:     t,
		Duration: duration,
	}, nil
}

type Event struct {
	Name     string
	Date     time.Time
	Loc      string
	Duration int
}

type Certificate struct {
	Type  config.CertificateType
	Event *Event

	canva     *gg.Context
	config    *config.Certificate
	signature string
}

func (c *Certificate) Width() float64 {
	return float64(c.canva.Width())
}

func (c *Certificate) Height() float64 {
	return float64(c.canva.Height())
}

func (c *Certificate) setFont(fontName assets.FontName, size float64) error {
	f, err := assets.LoadEmbededFont(fontName, size)
	if err != nil {
		return err
	}

	c.canva.SetFontFace(f)
	return nil
}

func (c *Certificate) rotateCanva(degrees float64) {
	radians := degrees * (math.Pi / 180)
	c.canva.RotateAbout(radians, c.Width()/2, c.Height()/2)
}

func (c *Certificate) setColorConfig(cc config.ColorConfig) {
	c.canva.SetColor(color.RGBA{
		R: uint8(cc.R),
		G: uint8(cc.G),
		B: uint8(cc.B),
		A: uint8(cc.A),
	})
}

func (c *Certificate) setBackground() {
	// background
	c.canva.DrawRectangle(0, 0, c.Width(), c.Height())
	c.setColorConfig(c.config.BackgroundColor)
	c.canva.Fill()

	// semi-transparent overlay
	m := c.config.OverlayMarginSize
	c.setColorConfig(c.config.OverlayColor)
	c.canva.DrawRectangle(m, m, c.Width()-(2.0*m), c.Height()-(2.0*m))
	c.canva.Fill()
}

func (c *Certificate) setCertificationTitle() error {
	if err := c.setFont(assets.OpenSans, c.config.TitleTextSize); err != nil {
		return err
	}
	c.setColorConfig(c.config.TextColor)
	title, err := c.config.GetTitleFromType(c.Type)
	if err != nil {
		return err
	}
	c.canva.DrawStringAnchored(
		title,
		c.Width()/2,
		c.Height()/4,
		0.5,
		0.5,
	)
	return nil
}

func (c *Certificate) setPersonName(name string) error {
	if err := c.setFont(assets.OpenSans, c.config.PersonTextSize); err != nil {
		return err
	}
	c.setColorConfig(c.config.TextColor)
	c.canva.DrawStringAnchored(
		name,
		c.Width()/2,
		c.Height()/2,
		0.5,
		0.5,
	)
	return nil
}

func (c *Certificate) setImgSignature() error {
	imgHeight := c.config.SignatureImgSize
	imgPath, err := c.config.MountSignaturePath(
		strings.ToLower(strings.ReplaceAll(c.signature, " ", "-")) + ".png",
	)
	if err != nil {
		return err
	}
	img, err := gg.LoadImage(imgPath)
	if err != nil {
		return err
	}
	resizedSignature := imaging.Resize(img, 0, imgHeight, imaging.Lanczos)
	c.canva.DrawImageAnchored(
		resizedSignature,
		int(c.Width()/2),
		int(5*c.Height()/6),
		0.5,
		0.5,
	)

	if err := c.setFont(assets.OpenSans, c.config.SignatureTitleSize); err != nil {
		return err
	}
	c.setColorConfig(c.config.SignatureTextColor)
	c.canva.DrawStringAnchored(
		strings.Repeat("_", c.config.SignatureLineLength),
		c.Width()/2,
		(5*c.Height()/6)+float64(imgHeight/2),
		0.5,
		0.5,
	)

	_, h := c.canva.MeasureString(c.signature)
	c.canva.DrawStringAnchored(
		c.signature,
		c.Width()/2,
		5*c.Height()/6+2*h+float64(imgHeight/2),
		0.5,
		0.5,
	)
	return nil
}

func (c *Certificate) setTextSignature() error {
	if err := c.setFont(assets.CedarvilleCursive, c.config.SignatureTextSize); err != nil {
		return err
	}
	c.setColorConfig(c.config.SignatureTextColor)
	c.canva.DrawStringAnchored(
		c.signature,
		c.Width()/2,
		5*c.Height()/6,
		0.5,
		0.5,
	)
	_, signatureHeight := c.canva.MeasureString(c.signature)
	if err := c.setFont(assets.OpenSans, c.config.SignatureTitleSize); err != nil {
		return err
	}
	c.setColorConfig(c.config.SignatureTitleColor)
	c.canva.DrawStringAnchored(
		strings.Repeat("_", c.config.SignatureLineLength),
		c.Width()/2,
		(5*c.Height()/6)+float64(signatureHeight/2),
		0.5,
		0.5,
	)

	_, h := c.canva.MeasureString(c.signature)
	c.canva.DrawStringAnchored(
		c.signature,
		c.Width()/2,
		5*c.Height()/6+2*h+float64(signatureHeight/2),
		0.5,
		0.5,
	)

	return nil
}

func (c *Certificate) setSignature() error {
	err := c.setImgSignature()
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		err = c.setTextSignature()
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *Certificate) setVerificationHash() error {
	if err := c.setFont(assets.OpenSans, c.config.ValidatorTextSize); err != nil {
		return err
	}
	c.setColorConfig(c.config.ValidatorTextColor)

	// Rotate text
	margin := c.config.OverlayMarginSize + 15
	c.rotateCanva(90)

	// Draw text
	c.canva.DrawStringAnchored(
		fmt.Sprintf("ID: %s", "asdasdasdasdasdasdaa:12345678901234567890"),
		c.Width()/2,
		margin-c.Height()/2,
		0.5,
		0.5,
	)

	// Rotate back
	c.rotateCanva(-90)

	return nil
}

func (c *Certificate) generate(name string) error {
	c.setBackground()
	if err := c.setCertificationTitle(); err != nil {
		return err
	}
	if err := c.setPersonName(name); err != nil {
		return err
	}
	if err := c.setSignature(); err != nil {
		return err
	}
	if err := c.setVerificationHash(); err != nil {
		return err
	}

	return nil
}

func (c *Certificate) save(identifier string) error {
	fileName := fmt.Sprintf(
		"%s-%s-%s-%s.png",
		c.Event.Date.Format("02-01-2006"),
		string(c.Type),
		c.Event.Name,
		identifier,
	)
	outputPath, err := c.config.MountOutputPath(
		strings.ToLower(strings.ReplaceAll(fileName, " ", "-")),
	)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	err = os.MkdirAll(
		filepath.Dir(outputPath),
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	if err := c.canva.SavePNG(outputPath); err != nil {
		return err
	}
	return nil
}
