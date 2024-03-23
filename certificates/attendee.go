package certificates

import (
	"fmt"
)

func NewAttendanceCertificate(attendee Attendee, event Event, signature string) (*AttendanceCertificate, error) {
	c, err := newCertificate(event, signature, AttendanceCertification)
	if err != nil {
		return nil, err
	}

	return &AttendanceCertificate{
		Attendee:    &attendee,
		certificate: *c,
	}, nil
}

type Attendee struct {
	Name  string
	Email string
}

type AttendanceCertificate struct {
	Attendee *Attendee
	certificate
}

func (c *AttendanceCertificate) Generate() error {
	if err := c.generate(c.Attendee.Name); err != nil {
		return err
	}

	if err := c.loadDefaultFont(30); err != nil {
		return err
	}

	line := fmt.Sprintf(
		"participou do %s, realizado no dia %s,",
		c.Event.Name, c.Event.Date.Format("02/01/2006"),
	)
	c.canva.DrawStringAnchored(
		line,
		c.size.width/2,
		5*c.size.height/8,
		0.5,
		0.5,
	)
	_, h := c.canva.MeasureString(line)
	c.canva.DrawStringAnchored(
		fmt.Sprintf(
			"nas instalações da %s, com carga horária total de %d horas.",
			c.Event.Loc, c.Event.Duration,
		),
		c.size.width/2,
		(5*c.size.height/8)+2*h,
		0.5,
		0.5,
	)

	if err := c.save(); err != nil {
		return err
	}

	return nil
}
