package certifigo

import (
	"github.com/wneessen/go-mail"
)

func NewGMailSender(sender, password string) (*EmailSender, error) {
	client, err := mail.NewClient(
		"smtp.gmail.com",
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(sender),
		mail.WithPassword(password),
	)
	if err != nil {
		return nil, err
	}

	return &EmailSender{
		client: client,
		sender: sender,
	}, nil
}

type Email struct {
	Subject     string
	Body        string
	To          string
	Attachments []string
}

type EmailSender struct {
	client *mail.Client
	sender string
}

func (s *EmailSender) mountMsgFromEmail(email Email) (*mail.Msg, error) {
	message := mail.NewMsg()
	if err := message.From(s.sender); err != nil {
		return nil, err
	}
	if err := message.To(email.To); err != nil {
		return nil, err
	}
	message.Subject(email.Subject)
	message.SetBodyString(mail.TypeTextPlain, email.Body)
	for _, certificationPath := range email.Attachments {
		message.AttachFile(certificationPath)
	}
	message.SetMessageID()
	message.SetDate()

	return message, nil
}

func (s *EmailSender) Send(email Email) error {
	message, err := s.mountMsgFromEmail(email)
	if err != nil {
		return err
	}

	if err := s.client.DialAndSend(message); err != nil {
		return err
	}
	return nil
}

func (s *EmailSender) BulkSend(emails []Email) error {
	var messages []*mail.Msg
	for _, email := range emails {
		message, err := s.mountMsgFromEmail(email)
		if err != nil {
			return err
		}
		message.SetBulk()
		messages = append(messages, message)
	}

	if err := s.client.DialAndSend(messages...); err != nil {
		return err
	}
	return nil
}
