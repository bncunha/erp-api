package controller

import (
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/labstack/echo/v4"
)

type DashboardController struct {
	dashboardService service.DashboardService
}

func NewDashboardController(dashboardService service.DashboardService) *DashboardController {
	return &DashboardController{dashboardService}
}

func (c *DashboardController) GetWidgets(context echo.Context) error {
	items, err := c.dashboardService.ListWidgets(context.Request().Context())
	if err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			return context.JSON(_http.StatusForbidden, http.HandleError(err))
		}
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, items)
}

func (c *DashboardController) GetWidgetData(context echo.Context) error {
	var req request.DashboardWidgetDataRequest
	if err := context.Bind(&req); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	output, err := c.dashboardService.GetWidgetData(context.Request().Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			return context.JSON(_http.StatusForbidden, http.HandleError(err))
		}
		if errors.Is(err, service.ErrDashboardWidgetNotFound) {
			return context.JSON(_http.StatusNotFound, http.HandleError(err))
		}
		if errors.Is(err, service.ErrDashboardResellerNotFound) {
			return context.JSON(_http.StatusNotFound, http.HandleError(err))
		}
		if errors.Is(err, service.ErrDashboardInvalidPeriod) {
			return context.JSON(_http.StatusBadRequest, http.HandleError(err))
		}
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, output)
}
