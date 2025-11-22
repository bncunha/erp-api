package email_brevo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/bncunha/erp-api/src/application/ports"
	"github.com/bncunha/erp-api/src/infrastructure/logs"
)

type EmailBrevoConfig struct {
	ApiKey string
}

type emailBrevo struct {
	config EmailBrevoConfig
}

func NewEmailBrevo(config EmailBrevoConfig) ports.EmailPort {
	return &emailBrevo{
		config: config,
	}
}

func (e *emailBrevo) Send(senderEmail string, senderName string, toEmail string, toName string, subject string, body string) error {
	url := "https://api.brevo.com/v3/smtp/email"

	request := NewEmailBrevoRequest(senderEmail, senderName, toEmail, toName, subject, body)

	jsonBody, err := json.Marshal(request)
	if err != nil {
		logs.Logger.Errorf("Erro ao serializar email brevo: %v", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		logs.Logger.Errorf("Erro ao criar requisição para email brevo: %v", err)
		return err
	}

	req.Header.Add("api-key", e.config.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.Logger.Errorf("Erro ao enviar requisição para email brevo: %v", err)
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.Logger.Errorf("Erro ao ler resposta da requisição para email brevo: %v", err)
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		logs.Logger.Errorf("Erro ao enviar email via brevo. Status: %d, Response: %s", resp.StatusCode, string(responseBody))
		return err
	}

	return nil
}
