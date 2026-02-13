package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type fakeTransaction struct {
	attributes map[string]interface{}
}

func (f *fakeTransaction) AddAttribute(key string, value interface{}) {
	if f.attributes == nil {
		f.attributes = map[string]interface{}{}
	}
	f.attributes[key] = value
}

func TestNREchoMiddlewareHasTransactionInContext(t *testing.T) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigEnabled(false),
		newrelic.ConfigAppName("test-app"),
		newrelic.ConfigLicense("1234567890123456789012345678901234567890"),
	)
	if err != nil {
		t.Fatalf("failed to create new relic app: %v", err)
	}

	e := echo.New()
	e.Use(nrecho.Middleware(app))
	e.GET("/instrumented", func(c echo.Context) error {
		if nrecho.FromContext(c) == nil {
			return c.String(http.StatusInternalServerError, "missing transaction")
		}
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/instrumented", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestNewRelicTransactionAttributes_OnlyAllowedAttributes(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctx := context.WithValue(c.Request().Context(), constants.TENANT_KEY, float64(42))
	ctx = context.WithValue(ctx, constants.USERID_KEY, int64(7))
	ctx = context.WithValue(ctx, constants.ROLE_KEY, "admin")
	ctx = context.WithValue(ctx, constants.USERNAME_KEY, "should-not-be-sent")
	c.SetRequest(c.Request().WithContext(ctx))
	c.Set(CtxKeyReqID, "req-123")

	fakeTxn := &fakeTransaction{attributes: map[string]interface{}{}}
	mw := newRelicTransactionAttributesWithExtractor(func(_ echo.Context) transactionAttributeAdder {
		return fakeTxn
	})

	err := mw(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})(c)
	if err != nil {
		t.Fatalf("unexpected middleware error: %v", err)
	}

	if len(fakeTxn.attributes) != 4 {
		t.Fatalf("expected 4 attributes, got %d", len(fakeTxn.attributes))
	}
	if _, ok := fakeTxn.attributes["tenant_id"]; !ok {
		t.Fatalf("expected tenant_id attribute")
	}
	if _, ok := fakeTxn.attributes["user_id"]; !ok {
		t.Fatalf("expected user_id attribute")
	}
	if _, ok := fakeTxn.attributes["role"]; !ok {
		t.Fatalf("expected role attribute")
	}
	if _, ok := fakeTxn.attributes["request_id"]; !ok {
		t.Fatalf("expected request_id attribute")
	}
	if _, ok := fakeTxn.attributes["username"]; ok {
		t.Fatalf("username should not be included as transaction attribute")
	}
}
