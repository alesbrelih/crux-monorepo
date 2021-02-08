package mail

import "net/smtp"

func NewMailService(host, port, from, password string) MailService {
	return &mailService{
		host:     host,
		port:     port,
		from:     from,
		password: password,
	}
}

type MailService interface {
	Send(to, subject, host string) error
}

type mailService struct {
	host     string
	port     string
	from     string
	password string
}

func (m *mailService) Send(to, subject, host string) error {
	auth := m.authenticate()

}

func (m *mailService) hostWithPort() string {
	return m.host + ":" + m.port
}

func (m *mailService) authenticate() smtp.Auth {
	return smtp.PlainAuth("", m.from, m.password, m.host)
}
