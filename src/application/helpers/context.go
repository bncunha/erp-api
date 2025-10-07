package helper

import (
	"context"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
)

func GetRole(ctx context.Context) domain.Role {
	role := ctx.Value(constants.ROLE_KEY).(string)
	return domain.Role(role)
}
