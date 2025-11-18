package emailusecase

import (
	"context"
	"fmt"

	"github.com/bncunha/erp-api/src/domain"
)

const (
	RecoverPassword             = "Trinus | Recuperação de Senha"
	RecoverPasswordBodyTemplate = "<html><body><p>Olá %s,</p><p>Recebemos uma solicitação para redefinir sua senha. Para redefinir sua senha, clique no link abaixo:</p><p><a href=\"%s\">%s</a></p><p>Se você não puder clicar no link, copie e cole o endereço abaixo:</p><p>%s</p><p>Se você não solicitou a redefinição de senha, ignore este e-mail.</p><p>Atenciosamente,</p><p>Equipe Trinus</p></body></html>"
)

func (e *emailUseCase) SendRecoverPassword(ctx context.Context, user domain.User, code string, uuid string) error {
	frontEndLink := fmt.Sprintf("%s/redefinir-senha?code=%s&uuid=%s", e.config.FRONTEND_URL, code, uuid)
	body := fmt.Sprintf(RecoverPasswordBodyTemplate, user.Name, frontEndLink, "Clique aqui para redefinir sua senha de acesso!", frontEndLink)
	return e.emailPort.Send(user.Email, RecoverPassword, body)
}
