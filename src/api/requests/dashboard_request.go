package request

import (
	"time"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/validator"
)

type DashboardWidgetPeriodRequest struct {
	From string `json:"from" validate:"required"`
	To   string `json:"to" validate:"required"`
}

type DashboardWidgetFiltersRequest struct {
	ResellerId *int64 `json:"reseller_id"`
	ProductId  *int64 `json:"product_id"`
}

type DashboardWidgetDataRequest struct {
	Enum    string                         `json:"enum" validate:"required"`
	Period  DashboardWidgetPeriodRequest   `json:"period" validate:"required"`
	Filters *DashboardWidgetFiltersRequest `json:"filters"`
}

func (r *DashboardWidgetDataRequest) Validate() error {
	if err := validator.Validate(r); err != nil {
		return err
	}

	from, err := time.Parse(time.DateOnly, r.Period.From)
	if err != nil {
		return errors.New("Período inicial inválido")
	}
	to, err := time.Parse(time.DateOnly, r.Period.To)
	if err != nil {
		return errors.New("Período final inválido")
	}
	if from.After(to) {
		return errors.New("Per inválido")
	}

	return nil
}
