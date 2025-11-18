package input

import "github.com/bncunha/erp-api/src/domain"

type CreateUserTokenInput struct {
	User      domain.User
	CreatedBy domain.User
	Type      domain.UserTokenType
}
