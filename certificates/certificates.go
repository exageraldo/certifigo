package certificates

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/exageraldo/suacuna-cli/fonts"
	"github.com/fogleman/gg"
	"github.com/spf13/viper"
)

type CertificationType string

const (
	AttendanceCertification CertificationType = "attendee"
	SpeakerCertification    CertificationType = "speaker"
)

func newCertificate(event Event, signature string, certificationType CertificationType) (*certificate, error) {
	var (
		width  = viper.GetInt("canva.size.w")
		height = viper.GetInt("canva.size.h")
	)

	fmt.Println(width, height)

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
		viper.GetString("SIGNATURES_DIR"),
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

func (c *certificate) getCertificateTitle() string {
	return viper.GetString(map[CertificationType]string{
		AttendanceCertification: "canva.attendance_title",
		SpeakerCertification:    "canva.speaker_title",
	}[c.Type])
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

	// solid color background
	bgColor := viper.GetStringMap("canva.background_color")
	c.canva.SetColor(color.RGBA{
		R: uint8(bgColor["r"].(int64)),
		G: uint8(bgColor["g"].(int64)),
		B: uint8(bgColor["b"].(int64)),
		A: uint8(bgColor["a"].(int64)),
	})
	c.canva.Fill()

	// semi-transparent overlay
	margin := viper.GetFloat64("canva.overlay_margin_size")
	x := margin
	y := margin
	w := c.size.width - (2.0 * margin)
	h := c.size.height - (2.0 * margin)

	overlayColor := viper.GetStringMap("canva.overlay_color")
	c.canva.SetColor(color.RGBA{
		R: uint8(overlayColor["r"].(int64)),
		G: uint8(overlayColor["g"].(int64)),
		B: uint8(overlayColor["b"].(int64)),
		A: uint8(overlayColor["a"].(int64)),
	})
	c.canva.DrawRectangle(x, y, w, h)
	c.canva.Fill()
}

func (c *certificate) setCertificationType() error {
	if err := c.loadDefaultFont(80); err != nil {
		return err
	}
	c.setTextColor()

	c.canva.DrawStringAnchored(
		c.getCertificateTitle(),
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
	c.setTextColor()

	c.canva.DrawStringAnchored(
		name,
		c.size.width/2,
		c.size.height/2,
		0.5,
		0.5,
	)

	return nil
}

func (c *certificate) setTextColor() {
	// txtColor := color.RGBA{}
	// if err := viper.UnmarshalKey("CANVA.TEXT_COLOR", &txtColor); err != nil {
	// 	fmt.Fprintf(os.Stderr, "%s\n", err)
	// 	os.Exit(1)
	// }
	// c.canva.SetColor(txtColor)

	textColor := viper.GetStringMap("canva.text_color")
	fmt.Println(textColor)
	c.canva.SetColor(color.RGBA{
		R: uint8(textColor["r"].(int64)),
		G: uint8(textColor["g"].(int64)),
		B: uint8(textColor["b"].(int64)),
		A: uint8(textColor["a"].(int64)),
	})
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
	c.setTextColor()

	c.canva.DrawStringAnchored(
		strings.Repeat("_", viper.GetInt("canva.signature_line_length")),
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
	if err := c.loadSignatureFont(); err != nil {
		return err
	}
	c.setTextColor()

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
		strings.Repeat("_", viper.GetInt("canva.signature_line_length")),
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
