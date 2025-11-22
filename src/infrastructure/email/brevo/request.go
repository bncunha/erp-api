package email_brevo

type EmailBrevoRequest struct {
	Sender      EmailBrevoSender `json:"sender"`
	To          []EmailBrevoTo   `json:"to"`
	Subject     string           `json:"subject"`
	HtmlContent string           `json:"htmlContent"`
}

func NewEmailBrevoRequest(senderEmail string, senderName string, toEmail string, toName string, subject string, body string) EmailBrevoRequest {
	return EmailBrevoRequest{
		Sender: EmailBrevoSender{
			Name:  senderName,
			Email: senderEmail,
		},
		To: []EmailBrevoTo{{
			Name:  toName,
			Email: toEmail,
		}},
		Subject:     subject,
		HtmlContent: body,
	}
}

type EmailBrevoSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type EmailBrevoTo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
