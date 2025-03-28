package certifigo

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

func NewCertificateDrawer(cType CertificateType, event Event, config CertificateConfigFile) *CertificateDrawer {
	return &CertificateDrawer{
		Type:  cType,
		Event: event,

		canva: gg.NewContext(
			config.CanvaSize.Width,
			config.CanvaSize.Height,
		),
		config: config,
	}
}

type CertificateDrawer struct {
	Type  CertificateType
	Event Event

	canva  *gg.Context
	config CertificateConfigFile
}

// Width returns the width of the canvas as a float64 value.
// It acts as an alias to the Width method of the underlying canvas object.
func (c *CertificateDrawer) Width() float64 {
	return float64(c.canva.Width())
}

// Height returns the height of the canvas as a float64 value.
// It acts as an alias to the Height method of the underlying canvas object.
func (c *CertificateDrawer) Height() float64 {
	return float64(c.canva.Height())
}

// useColor sets the drawing color for the CertificateDrawer's canvas.
// It takes a HexColor as input, which contains the RGBA values to define the color.
// The color is then applied to the canvas for subsequent drawing operations.
//
// Parameters:
//   - hColor: A HexColor struct containing the red (R), green (G), blue (B),
//     and alpha (A) values of the color to be set.
func (c *CertificateDrawer) useColor(hColor HexColor) {
	c.canva.SetColor(color.RGBA{
		R: hColor.R,
		G: hColor.G,
		B: hColor.B,
		A: hColor.A,
	})
}

// useFont sets the font for the CertificateDrawer's canvas.
// It loads the specified font by name and size, and applies it to the canvas.
//
// Parameters:
//   - fontName: The name of the font to be loaded.
//   - size: The size of the font to be loaded.
//
// Returns:
//   - error: An error if the font could not be loaded, or nil if successful.
func (c *CertificateDrawer) useFont(fontName string, size float64) error {
	f, err := LoadFont(fontName, size)
	if err != nil {
		return err
	}

	c.canva.SetFontFace(f)
	return nil
}

func (c *CertificateDrawer) drawBackground() {
	// background
	c.canva.DrawRectangle(0, 0, c.Width(), c.Height())
	c.useColor(c.config.Background.BorderColor)
	c.canva.Fill()

	// semi-transparent overlay
	m := c.config.Background.BorderSize
	c.useColor(c.config.Background.Color)
	c.canva.DrawRectangle(m, m, c.Width()-(2.0*m), c.Height()-(2.0*m))
	c.canva.Fill()
}

// drawLogoImg draws the logo image onto the canvas if a logo path is provided.
// It first checks if the logo path is empty and returns early if so. Otherwise,
// it resolves the absolute path of the logo file and attempts to load the image.
// The loaded image is resized to a height of 200 pixels while maintaining its
// aspect ratio. The resized image is then drawn onto the canvas, anchored at
// the center horizontally and at one-fifth of the canvas height vertically.
//
// Returns an error if the logo path is invalid, the image cannot be loaded, or
// any other issue occurs during the process.
func (c *CertificateDrawer) drawLogoImg() error {
	if c.Event.Logo == "" {
		return nil
	}
	logoPath, err := filepath.Abs(c.Event.Logo)
	if err != nil {
		return err
	}
	img, err := gg.LoadImage(logoPath)
	if err != nil {
		return err
	}

	resizedSignature := imaging.Resize(img, 0, 200, imaging.Lanczos)
	c.canva.DrawImageAnchored(
		resizedSignature,
		int(c.Width()/2),
		int(c.Height()/5),
		0.5,
		0.5,
	)
	return nil
}

func (c *CertificateDrawer) drawCertificationTitle() error {
	var title string
	switch c.Type {
	case AttendanceCertification:
		title = c.config.Attendee.Title
	case SpeakerCertification:
		title = c.config.Speaker.Title
	default:
		return fmt.Errorf("invalid certificate type: %v", c.Type)
	}

	if err := c.useFont(OpenSans, c.config.Text.TitleTextSize); err != nil {
		return err
	}
	c.useColor(c.config.Text.TitleTextColor)

	height := c.Height() / 4
	if c.Event.Logo != "" {
		height = c.Height() / 3
	}

	c.canva.DrawStringAnchored(
		title,
		c.Width()/2,
		height,
		0.5,
		0.5,
	)
	return nil
}

func (c *CertificateDrawer) drawPersonName(name string) error {
	if err := c.useFont(OpenSans, c.config.Text.PersonTextSize); err != nil {
		return err
	}
	c.useColor(c.config.Text.TextColor)
	c.canva.DrawStringAnchored(
		name,
		c.Width()/2,
		c.Height()/2,
		0.5,
		0.5,
	)
	return nil
}

func (c *CertificateDrawer) drawEventInfo() error {
	var info string
	switch c.Type {
	case AttendanceCertification:
		info = c.config.Attendee.Body
	case SpeakerCertification:
		info = c.config.Speaker.Body
	default:
		return fmt.Errorf("invalid certificate type: %v", c.Type)
	}

	if err := c.useFont(OpenSans, c.config.Text.TextSize); err != nil {
		return err
	}
	c.useColor(c.config.Text.TextColor)

	for idx, line := range strings.Split(info, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		_, h := c.canva.MeasureString(line)
		c.canva.DrawStringAnchored(
			line,
			c.Width()/2,
			(5*c.Height()/9)+((2+float64(idx))*h),
			0.5,
			0.5,
		)
	}

	return nil
}

func (c *CertificateDrawer) drawImgSignature() error {
	imgHeight := c.config.Signature.ImgSize
	imgPath, err := c.config.MountSignaturePath(
		strings.ToLower(strings.ReplaceAll(c.Event.Signature, " ", "-")) + ".png",
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

	if err := c.useFont(OpenSans, c.config.Signature.TitleSize); err != nil {
		return err
	}
	c.useColor(c.config.Signature.TextColor)
	c.canva.DrawStringAnchored(
		strings.Repeat("_", c.config.Signature.LineLength),
		c.Width()/2,
		(5*c.Height()/6)+float64(imgHeight/2),
		0.5,
		0.5,
	)

	_, h := c.canva.MeasureString(c.Event.Signature)
	c.canva.DrawStringAnchored(
		c.Event.Signature,
		c.Width()/2,
		5*c.Height()/6+2*h+float64(imgHeight/2),
		0.5,
		0.5,
	)
	return nil
}

func (c *CertificateDrawer) drawTextSignature() error {
	if err := c.useFont(CedarvilleCursive, c.config.Signature.TextSize); err != nil {
		return err
	}
	c.useColor(c.config.Signature.TextColor)
	c.canva.DrawStringAnchored(
		c.Event.Signature,
		c.Width()/2,
		5*c.Height()/6,
		0.5,
		0.5,
	)
	_, signatureHeight := c.canva.MeasureString(c.Event.Signature)
	if err := c.useFont(OpenSans, c.config.Signature.TitleSize); err != nil {
		return err
	}
	c.useColor(c.config.Signature.TitleColor)
	c.canva.DrawStringAnchored(
		strings.Repeat("_", c.config.Signature.LineLength),
		c.Width()/2,
		(5*c.Height()/6)+float64(signatureHeight/2),
		0.5,
		0.5,
	)

	_, h := c.canva.MeasureString(c.Event.Signature)
	c.canva.DrawStringAnchored(
		c.Event.Signature,
		c.Width()/2,
		5*c.Height()/6+2*h+float64(signatureHeight/2),
		0.5,
		0.5,
	)

	return nil
}

func (c *CertificateDrawer) drawSignature() error {
	if c.Event.SignatureImg != "" {
		return c.drawImgSignature()
	}
	return c.drawTextSignature()
}

func (c *CertificateDrawer) DrawAndSave(personName string) (string, error) {
	c.drawBackground()
	if err := c.drawLogoImg(); err != nil {
		return "", err
	}
	if err := c.drawCertificationTitle(); err != nil {
		return "", err
	}
	if err := c.drawPersonName(personName); err != nil {
		return "", err
	}
	if err := c.drawEventInfo(); err != nil {
		return "", err
	}
	if err := c.drawSignature(); err != nil {
		return "", err
	}

	fileName := strings.ToLower(strings.ReplaceAll(fmt.Sprintf(
		"%s-%s-%s.png",
		c.Event.Name,
		string(c.Type),
		personName,
	), " ", "-"))

	outputPath, err := c.config.MountOutputPath(fileName)
	if err != nil {
		return "", err
	}

	// Create directory if it doesn't exist
	err = os.MkdirAll(
		filepath.Dir(outputPath),
		os.ModePerm,
	)
	if err != nil {
		return "", err
	}

	if err := c.canva.SavePNG(outputPath); err != nil {
		return "", err
	}
	return outputPath, nil
}
