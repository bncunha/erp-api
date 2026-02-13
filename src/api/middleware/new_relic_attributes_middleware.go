package middleware

import (
	"fmt"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
)

type transactionAttributeAdder interface {
	AddAttribute(string, interface{})
}

func NewRelicTransactionAttributes() echo.MiddlewareFunc {
	return newRelicTransactionAttributesWithExtractor(func(c echo.Context) transactionAttributeAdder {
		return nrecho.FromContext(c)
	})
}

func newRelicTransactionAttributesWithExtractor(extractTxn func(echo.Context) transactionAttributeAdder) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			txn := extractTxn(c)
			if txn != nil {
				addContextAttribute(txn, "tenant_id", c.Request().Context().Value(constants.TENANT_KEY))
				addContextAttribute(txn, "user_id", c.Request().Context().Value(constants.USERID_KEY))
				addContextAttribute(txn, "role", c.Request().Context().Value(constants.ROLE_KEY))
				addContextAttribute(txn, "request_id", c.Get(CtxKeyReqID))
			}
			return next(c)
		}
	}
}

func addContextAttribute(txn transactionAttributeAdder, key string, value interface{}) {
	if value == nil {
		return
	}

	switch v := value.(type) {
	case string:
		if v == "" {
			return
		}
		txn.AddAttribute(key, v)
	case fmt.Stringer:
		txn.AddAttribute(key, v.String())
	case int:
		txn.AddAttribute(key, v)
	case int32:
		txn.AddAttribute(key, v)
	case int64:
		txn.AddAttribute(key, v)
	case float64:
		txn.AddAttribute(key, int64(v))
	case bool:
		txn.AddAttribute(key, v)
	}
}
