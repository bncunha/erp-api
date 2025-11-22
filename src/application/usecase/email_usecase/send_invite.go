package emailusecase

import (
	"context"
	"fmt"

	"github.com/bncunha/erp-api/src/domain"
)

const (
	InviteSubject      = "Trinus | Você foi convidado para o Trinus!"
	InviteBodyTemplate = "<html><body><p>Olá %s,</p><p>Você foi convidado para o Trinus! Para acessar o sistema, clique no link abaixo:</p><p><a href=\"%s\">%s</a></p><p>Se você não puder clicar no link, copie e cole o endereço abaixo:</p><p>%s</p><p>Atenciosamente,</p><p>Equipe Trinus</p></body></html>"
)

func (e *emailUseCase) SendInvite(ctx context.Context, user domain.User, code string, uuid string) error {
	frontEndLink := fmt.Sprintf("%s/redefinir-senha?code=%s&uuid=%s&new_user=true", e.config.FRONTEND_URL, code, uuid)
	body := fmt.Sprintf(InviteBodyTemplate, user.Name, frontEndLink, "Clique aqui para definir sua senha de acesso!", frontEndLink)
	return e.emailPort.Send(SenderEmail, SenderName, user.Email, user.Name, InviteSubject, body)
}
