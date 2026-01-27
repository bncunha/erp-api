package middleware

import (
	_http "net/http"
	"strings"
	"time"

	"github.com/bncunha/erp-api/src/api/responses"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/labstack/echo/v4"
)

func BillingWriteGuard() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			method := ctx.Request().Method
			if method == _http.MethodGet || method == _http.MethodHead || method == _http.MethodOptions {
				return next(ctx)
			}

			path := ctx.Request().URL.Path
			if strings.HasPrefix(path, "/billing") {
				return next(ctx)
			}

			canWriteValue := ctx.Request().Context().Value(constants.BILLING_CAN_WRITE_KEY)
			canWrite, _ := canWriteValue.(bool)
			if canWrite {
				return next(ctx)
			}

			return ctx.JSON(_http.StatusForbidden, response.BillingForbiddenResponse{
				Code:             "BILLING_BLOCKED",
				Message:          "Plano nao permite alteracoes.",
				Plan:             "",
				CurrentPeriodEnd: time.Time{},
			})
		}
	}
}
