package email_smtp

import (
	"fmt"
	"net/smtp"

	"github.com/bncunha/erp-api/src/application/ports"
	"github.com/bncunha/erp-api/src/infrastructure/logs"
)

type EmailSmtpConfig struct {
	Host     string
	Port     int64
	Username string
	Password string
}

type emailSmtp struct {
	config EmailSmtpConfig
}

func NewEmailSmtp(config EmailSmtpConfig) ports.EmailPort {
	return &emailSmtp{
		config: config,
	}
}

func (e *emailSmtp) authenticate() smtp.Auth {
	return smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Host)
}

func (e *emailSmtp) Send(senderEmail string, senderName string, toEmail string, toName string, subject string, body string) error {
	auth := e.authenticate()
	host := fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)
	msg := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
		"\r\n"+
		"%s\r\n", toEmail, subject, body)
	logs.Logger.Infof("%s - %s - %s", host, e.config.Username, msg)
	return smtp.SendMail(host, auth, e.config.Username, []string{toEmail}, []byte(msg))
}
