package observability

import (
	config "github.com/bncunha/erp-api/src/main"
	"github.com/labstack/echo/v4"
)

type Observability interface {
	SetupObservability(*config.Config) error
	SetEchoMiddleware(e *echo.Echo)
	GetApp() any
}

type Transaction interface {
	AddAttribute(key string, value interface{})
	End() error
}

type observability struct {
	observabilityApp Observability
}

func NewObservability(observabilityApp Observability) *observability {
	return &observability{
		observabilityApp: observabilityApp,
	}
}

func (o *observability) GetApp() any {
	return o.observabilityApp.GetApp()
}

func (o *observability) SetupObservability(config *config.Config) error {
	return o.observabilityApp.SetupObservability(config)
}

func (o *observability) SetEchoMiddleware(e *echo.Echo) {
	o.observabilityApp.SetEchoMiddleware(e)
}
