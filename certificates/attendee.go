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
	c.setTextColor()

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
