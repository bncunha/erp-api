package observability

import (
	config "github.com/bncunha/erp-api/src/main"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type NewRelicObservability struct {
	app *newrelic.Application
}

func NewNewRelicObservability() Observability {
	return &NewRelicObservability{}
}

func (n *NewRelicObservability) SetupObservability(config *config.Config) error {
	app, err := newrelic.NewApplication(
		newrelic.ConfigEnabled(config.NR_ENABLED),
		newrelic.ConfigAppName("erp-api"),
		newrelic.ConfigLicense(config.NR_LICENSE_KEY),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	n.app = app
	return err
}

func (n *NewRelicObservability) SetEchoMiddleware(e *echo.Echo) {
	e.Use(nrecho.Middleware(n.app, nrecho.WithSkipper(func(c echo.Context) bool {
		return c.Request().URL.Path == "/health"
	})))
}

func (n *NewRelicObservability) GetApp() any {
	return n.app
}

type NewRelicTransaction struct {
	txn *newrelic.Transaction
}

func (n *NewRelicTransaction) End() error {
	n.txn.End()
	return nil
}
