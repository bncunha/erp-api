package emailusecase

import (
	"context"

	"github.com/bncunha/erp-api/src/application/ports"
	"github.com/bncunha/erp-api/src/domain"
	config "github.com/bncunha/erp-api/src/main"
)

type EmailUseCase interface {
	SendInvite(ctx context.Context, user domain.User, code string, uuid string) error
	SendRecoverPassword(to string, code string) error
}

type emailUseCase struct {
	config    *config.Config
	emailPort ports.EmailPort
}

func NewEmailUseCase(config *config.Config, emailPort ports.EmailPort) EmailUseCase {
	return &emailUseCase{config, emailPort}
}
