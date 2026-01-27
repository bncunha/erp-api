package controller

import (
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	"github.com/bncunha/erp-api/src/api/viewmodel"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/labstack/echo/v4"
)

type BillingController struct {
	billingService service.BillingService
}

func NewBillingController(billingService service.BillingService) *BillingController {
	return &BillingController{
		billingService: billingService,
	}
}

func (c *BillingController) Status(context echo.Context) error {
	status, err := c.billingService.GetStatus(context.Request().Context())
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, viewmodel.ToBillingStatusViewModel(status))
}

func (c *BillingController) Summary(context echo.Context) error {
	summary, err := c.billingService.GetSummary(context.Request().Context())
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, viewmodel.ToBillingSummaryViewModel(summary))
}

func (c *BillingController) Payments(context echo.Context) error {
	payments, err := c.billingService.GetPayments(context.Request().Context())
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, viewmodel.ToBillingPaymentViewModels(payments))
}
