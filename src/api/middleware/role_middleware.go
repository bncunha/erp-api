package middleware

import (
	"fmt"
	_http "net/http"
	"strings"

	"github.com/bncunha/erp-api/src/api/http"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/labstack/echo/v4"
)

func RoleMiddleware(roles []domain.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				return c.JSON(_http.StatusForbidden, http.HandleError(fmt.Errorf("Token inválido")))
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			_, _, role, _, err := helper.ParseJWT(token)
			if err != nil {
				return c.JSON(_http.StatusForbidden, http.HandleError(fmt.Errorf("Token inválido")))
			}

			for _, r := range roles {
				if r == domain.Role(role) {
					return next(c)
				}
			}

			return c.JSON(_http.StatusForbidden, http.HandleError(fmt.Errorf("Acesso negado")))
		}
	}
}
