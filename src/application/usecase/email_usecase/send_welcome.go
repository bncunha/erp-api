package emailusecase

import (
	"context"
	"fmt"
)

const (
	WelcomeSubject      = "Trinus | Bem-vindo ao Trinus!"
	WelcomeBodyTemplate = "<html><body><p>Olá %s,</p><p>Bem-vindo ao Trinus! Sua empresa foi cadastrada com sucesso.</p><p>Faça login com seu usuário para começar a usar a plataforma.</p><p><a href=\"%s\" target=\"_blank\">Acessar a plataforma</a></p><p>Atenciosamente,</p><p>Equipe Trinus</p></body></html>"
)

func (e *emailUseCase) SendWelcome(ctx context.Context, email string, name string) error {
	body := fmt.Sprintf(WelcomeBodyTemplate, name, e.config.FRONTEND_URL)
	return e.emailPort.Send(SenderEmail, SenderName, email, name, WelcomeSubject, body)
}
