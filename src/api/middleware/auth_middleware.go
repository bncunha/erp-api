package middleware

import (
	"context"
	"fmt"
	_http "net/http"
	"strings"

	"github.com/bncunha/erp-api/src/api/http"
	"github.com/bncunha/erp-api/src/application/constants"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.JSON(_http.StatusUnauthorized, http.HandleError(fmt.Errorf("Token inválido")))
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		username, tenant_id, role, user_id, billing, err := helper.ParseJWTWithBilling(token)
		if err != nil {
			return c.JSON(_http.StatusUnauthorized, http.HandleError(fmt.Errorf("Token inválido")))
		}

		ctx := context.WithValue(c.Request().Context(), constants.TENANT_KEY, tenant_id)
		ctx = context.WithValue(ctx, constants.USERNAME_KEY, username)
		ctx = context.WithValue(ctx, constants.USERID_KEY, user_id)
		ctx = context.WithValue(ctx, constants.ROLE_KEY, role)
		ctx = context.WithValue(ctx, constants.BILLING_CAN_WRITE_KEY, billing.CanWrite)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
