package service

type EmailService interface {
	Send(to string, subject string, body string) error
}

