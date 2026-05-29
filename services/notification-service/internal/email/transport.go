package email

import "gopkg.in/gomail.v2"

type Transport struct {
	host string
	port int
	from string
}

type Email struct {
	To       string
	Subject  string
	HTMLBody string
}

func NewTransport(host string, port int, from string) *Transport {
	return &Transport{
		host: host,
		port: port,
		from: from,
	}
}

func (s *Transport) Transport(email Email) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.from)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)

	m.SetBody("text/html", email.HTMLBody)

	dialer := gomail.NewDialer(
		s.host,
		s.port,
		"",
		"",
	)

	return dialer.DialAndSend(m)
}
