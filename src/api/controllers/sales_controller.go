package controller

import (
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/labstack/echo/v4"
)

type SalesController struct {
	salesService service.SalesService
}

func NewSalesController(salesService service.SalesService) *SalesController {
	return &SalesController{salesService}
}

func (c *SalesController) Create(context echo.Context) error {
	var salesRequeste request.CreateSaleRequest
	if err := context.Bind(&salesRequeste); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	err := c.salesService.CreateSales(context.Request().Context(), salesRequeste)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusCreated, nil)
}
