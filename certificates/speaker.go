package certificates

import "fmt"

func NewSpeakerCertificate(speaker Speaker, event Event, signature string) (*SpeakerCertificate, error) {
	c, err := newCertificate(event, signature, SpeakerCertification)
	if err != nil {
		return nil, err
	}

	return &SpeakerCertificate{
		Speaker:     &speaker,
		certificate: *c,
	}, nil
}

type Speaker struct {
	Name         string
	Email        string
	TalkTitle    string
	TalkDuration int // in minutes
}

type SpeakerCertificate struct {
	Speaker *Speaker
	certificate
}

func (c *SpeakerCertificate) Generate() error {
	if err := c.generate(c.Speaker.Name); err != nil {
		return err
	}

	if err := c.loadDefaultFont(30); err != nil {
		return err
	}
	c.setTextColor()

	line := fmt.Sprintf(
		"participou como palestrante no %s, realizado no dia %s,",
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
			"com a palestra %s, com duração de %.1f horas.",
			c.Speaker.TalkTitle, float64(c.Speaker.TalkDuration)/60,
		),
		c.Width()/2,
		(5*c.Height()/8)+2*h,
		0.5,
		0.5,
	)

	if err := c.save(c.Speaker.Name); err != nil {
		return err
	}

	return nil
}
