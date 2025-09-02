// pkg/httpmw/logging.go
package middleware

import (
	"time"

	"github.com/bncunha/erp-api/src/infrastructure/logs"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const CtxKeyReqID = "req_id"

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			reqID := c.Request().Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.NewString()
			}
			c.Set(CtxKeyReqID, reqID)
			c.Response().Header().Set("X-Request-ID", reqID)

			err := next(c)
			latency := time.Since(start)

			// Campos úteis
			fields := map[string]any{
				"req_id":     reqID,
				"method":     c.Request().Method,
				"path":       c.Path(),
				"status":     c.Response().Status,
				"ip":         c.RealIP(),
				"latency_ms": latency.Milliseconds(),
			}

			// (opcional) propague usuário/tenant no contexto e adicione aqui
			// if u := c.Get("user_id"); u != nil { fields["user_id"] = u }

			log := logs.Logger.With(fields)

			if err != nil {
				log.With(map[string]any{
					"body": c.Request().Body,
					"err":  err.Error(),
				}).Errorf("request completed with error: %v", err)
				return err
			}

			// níveis por status
			switch {
			case c.Response().Status >= 500:
				log.Errorf("request completed")
			case c.Response().Status >= 400:
				log.Warnf("request completed")
			default:
				log.Infof("request completed")
			}
			return nil
		}
	}
}

// Recovery que loga panic com req_id
func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					reqID, _ := c.Get(CtxKeyReqID).(string)
					logs.Logger.With(map[string]any{
						"req_id": reqID, "path": c.Path(),
					}).Errorf("panic: %v", r)
					err = echo.NewHTTPError(500, "internal server error")
				}
			}()
			return next(c)
		}
	}
}
