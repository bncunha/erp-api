package middleware

import (
	"context"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/labstack/echo/v4"
)

func TenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// // Supondo que o tenant_id vem no header
		// tenantIDStr := c.Request().Header.Get("X-Tenant-ID")
		// if tenantIDStr == "" {
		// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tenant ID missing"})
		// }

		// tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
		// if err != nil {
		// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Tenant ID"})
		// }

		// Adiciona no contexto da requisição
		ctx := context.WithValue(c.Request().Context(), constants.TENANT_KEY, 1)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
