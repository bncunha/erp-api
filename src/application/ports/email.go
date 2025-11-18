package ports

type EmailPort interface {
	Send(to string, subject string, body string) error
}
