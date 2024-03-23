package certificates

import (
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/exageraldo/suacuna-cli/fonts"
	"github.com/fogleman/gg"
)

type CertificationType string

const (
	AttendanceCertification CertificationType = "CERTIFICADO DE PARTICIPAÇÃO"
	SpeakerCertification    CertificationType = "CERTIFICADO DE PALESTRANTE"
)

func newCertificate(event Event, signature string, certificationType CertificationType) (*certificate, error) {
	const width, height = 1600, 800

	return &certificate{
		Type:  certificationType,
		Event: &event,
		canva: gg.NewContext(width, height),
		signature: &Signature{
			Name: signature,
		},
		size: struct {
			width  float64
			height float64
		}{
			width:  float64(width),
			height: float64(height),
		},
	}, nil
}

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

type Signature struct {
	Name string
}

func (s *Signature) GetPath() string {
	return filepath.Join(
		"signatures",
		strings.Replace(s.Name, " ", "-", -1)+".png",
	)
}

type certificate struct {
	Type  CertificationType
	Event *Event

	canva     *gg.Context
	signature *Signature
	size      struct {
		width, height float64
	}
}

func (c *certificate) loadSignatureFont() error {
	f, err := fonts.LoadFont(fonts.SignatureFont, 60)
	if err != nil {
		return err
	}

	c.canva.SetFontFace(f)
	return nil
}

func (c *certificate) loadDefaultFont(points float64) error {
	f, err := fonts.LoadFont(fonts.DefaultFont, points)
	if err != nil {
		return err
	}

	c.canva.SetFontFace(f)
	return nil
}

func (c *certificate) setBackground() {
	c.canva.DrawRectangle(0, 0, c.size.width, c.size.height)

	// solid color
	c.canva.SetColor(color.RGBA{92, 225, 230, 255})
	c.canva.Fill()

	// semi-transparent overlay
	margin := 20.0
	x := margin
	y := margin
	w := c.size.width - (2.0 * margin)
	h := c.size.height - (2.0 * margin)
	c.canva.SetColor(color.RGBA{0, 0, 0, 180})
	c.canva.DrawRectangle(x, y, w, h)
	c.canva.Fill()
}

func (c *certificate) setCertificationType() error {
	if err := c.loadDefaultFont(80); err != nil {
		return err
	}
	c.canva.SetColor(color.White)
	c.canva.DrawStringAnchored(
		string(c.Type),
		c.size.width/2,
		c.size.height/4,
		0.5,
		0.5,
	)
	return nil
}

func (c *certificate) setPersonName(name string) error {
	if err := c.loadDefaultFont(70); err != nil {
		return err
	}

	c.canva.SetColor(color.White)
	c.canva.DrawStringAnchored(
		name,
		c.size.width/2,
		c.size.height/2,
		0.5,
		0.5,
	)

	return nil
}

func (c *certificate) setImgSignature() error {
	signatureHeight := 100
	signatureImg, err := gg.LoadImage(c.signature.GetPath())
	if err != nil {
		return err
	}
	resizedSignature := imaging.Resize(signatureImg, 0, signatureHeight, imaging.Lanczos)

	c.canva.DrawImageAnchored(
		resizedSignature,
		int(c.size.width/2),
		int(5*c.size.height/6),
		0.5,
		0.5,
	)

	if err := c.loadDefaultFont(15); err != nil {
		return err
	}

	c.canva.SetColor(color.White)

	c.canva.DrawStringAnchored(
		"______________________",
		c.size.width/2,
		(5*c.size.height/6)+float64(signatureHeight/2),
		0.5,
		0.5,
	)

	_, h := c.canva.MeasureString(c.signature.Name)
	c.canva.DrawStringAnchored(
		c.signature.Name,
		c.size.width/2,
		5*c.size.height/6+2*h+float64(signatureHeight/2),
		0.5,
		0.5,
	)

	return nil
}

func (c *certificate) setTextSignature() error {
	c.canva.SetColor(color.White)
	if err := c.loadSignatureFont(); err != nil {
		return err
	}
	c.canva.DrawStringAnchored(
		c.signature.Name,
		c.size.width/2,
		5*c.size.height/6,
		0.5,
		0.5,
	)

	_, signatureHeight := c.canva.MeasureString(c.signature.Name)

	if err := c.loadDefaultFont(15); err != nil {
		return err
	}

	c.canva.DrawStringAnchored(
		"______________________",
		c.size.width/2,
		(5*c.size.height/6)+float64(signatureHeight/2),
		0.5,
		0.5,
	)

	_, h := c.canva.MeasureString(c.signature.Name)
	c.canva.DrawStringAnchored(
		c.signature.Name,
		c.size.width/2,
		5*c.size.height/6+2*h+float64(signatureHeight/2),
		0.5,
		0.5,
	)

	return nil
}

func (c *certificate) setSignature() error {
	err := c.setImgSignature()
	if err != nil {
		if os.IsNotExist(err) {
			c.setTextSignature()
		} else {
			return err
		}
	}
	return nil
}

func (c *certificate) generate(name string) error {
	c.setBackground()
	if err := c.setCertificationType(); err != nil {
		return err
	}
	if err := c.setPersonName(name); err != nil {
		return err
	}
	if err := c.setSignature(); err != nil {
		return err
	}

	return nil
}

func (c *certificate) save() error {
	filename := c.Event.Date.Format("02-01-2006-") + strings.ReplaceAll(
		string(c.Type)+"-"+c.Event.Name+"-"+time.Now().Format(time.DateTime),
		" ",
		"-",
	)
	if err := c.canva.SavePNG(filepath.Join(filename + ".png")); err != nil {
		return err
	}
	return nil
}
