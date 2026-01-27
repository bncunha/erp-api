package helper

import (
	"context"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

func GetRole(ctx context.Context) domain.Role {
	role := ctx.Value(constants.ROLE_KEY).(string)
	return domain.Role(role)
}

func GetTenantId(ctx context.Context) (int64, error) {
	tenantIdValue := ctx.Value(constants.TENANT_KEY)
	switch value := tenantIdValue.(type) {
	case int64:
		return value, nil
	case float64:
		return int64(value), nil
	default:
		return 0, errors.New("tenant id invalido")
	}
}
