package certifigo

type Event struct {
	Name     string     `toml:"name" json:"name"`
	Location string     `toml:"location" json:"location"`
	Date     StringDate `toml:"date" json:"date"`
	Duration int        `toml:"duration" json:"duration"`

	Signature    string `toml:"signature" json:"signature"`
	SignatureImg string `toml:"signature_img" json:"signature_img"`
	Logo         string `toml:"logo" json:"logo"`
}

type Speaker struct {
	Name         string `toml:"name" json:"name"`
	Email        string `toml:"email" json:"email"`
	TalkTitle    string `toml:"talk_title" json:"talk_title"`
	TalkDuration int    `toml:"talk_duration" json:"talk_duration"`
	Attendee     bool   `toml:"attendee" json:"attendee"`
	Notify       bool   `toml:"notify" json:"notify"`
}

type Attendee struct {
	Name   string `toml:"name" json:"name"`
	Email  string `toml:"email" json:"email"`
	Notify bool   `toml:"notify" json:"notify"`
}

type EventFile struct {
	Event     Event      `toml:"event" json:"event"`
	Speakers  []Speaker  `toml:"speakers" json:"speakers"`
	Attendees []Attendee `toml:"attendees" json:"attendees"`
}
