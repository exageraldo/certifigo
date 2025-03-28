package certifigo

type Event struct {
	Name     string     `toml:"name"`
	Location string     `toml:"location"`
	Date     StringDate `toml:"date"`
	Duration int        `toml:"duration"`

	Signature    string `toml:"signature"`
	SignatureImg string `toml:"signature_img"`
	Folder       string `toml:"folder"`
	Logo         string `toml:"logo"`
}

type Speaker struct {
	Name         string `toml:"name"`
	Email        string `toml:"email"`
	TalkTitle    string `toml:"talk_title"`
	TalkDuration int    `toml:"talk_duration"`
	Attendee     bool   `toml:"attendee"`
	Notify       bool   `toml:"notify"`
}

type Attendee struct {
	Name   string `toml:"name"`
	Email  string `toml:"email"`
	Notify bool   `toml:"notify"`
}

type EventFile struct {
	Event     Event      `toml:"event"`
	Speakers  []Speaker  `toml:"speakers"`
	Attendees []Attendee `toml:"attendees"`
}
