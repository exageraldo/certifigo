package certificates

import (
	"fmt"

	"github.com/exageraldo/suacuna-cli/assets"
	"github.com/exageraldo/suacuna-cli/config"
	"github.com/fogleman/gg"
)

func NewAttendanceCertificate(attendee Attendee, event Event, signature string, cfg config.Certificate) *AttendanceCertificate {
	return &AttendanceCertificate{
		Attendee: &attendee,
		Certificate: Certificate{
			Type:  config.AttendanceCertification,
			Event: &event,

			canva: gg.NewContext(
				cfg.CanvaSize.W,
				cfg.CanvaSize.H,
			),
			config:    &cfg,
			signature: signature,
		},
	}
}

type Attendee struct {
	Name  string
	Email string
}

type AttendanceCertificate struct {
	Attendee *Attendee
	Certificate
}

func (c *AttendanceCertificate) Generate() error {
	if err := c.generate(c.Attendee.Name); err != nil {
		return err
	}

	if err := c.setFont(assets.OpenSans, c.config.TextSize); err != nil {
		return err
	}
	c.setColorConfig(c.config.TextColor)

	line := fmt.Sprintf(
		"participou do %s, realizado no dia %s,",
		c.Event.Name, c.Event.Date.Format("02/01/2006"),
	)
	c.canva.DrawStringAnchored(
		line,
		c.Width()/2,
		5*c.Height()/8,
		0.5,
		0.5,
	)
	_, h := c.canva.MeasureString(line)
	c.canva.DrawStringAnchored(
		fmt.Sprintf(
			"nas instalações da %s, com carga horária total de %d horas.",
			c.Event.Loc, c.Event.Duration,
		),
		c.Width()/2,
		(5*c.Height()/8)+2*h,
		0.5,
		0.5,
	)

	if err := c.save(c.Attendee.Name); err != nil {
		return err
	}

	return nil
}
