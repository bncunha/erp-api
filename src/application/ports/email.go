package ports

type EmailPort interface {
	Send(senderEmail string, senderName string, toEmail string, toName string, subject string, body string) error
}
