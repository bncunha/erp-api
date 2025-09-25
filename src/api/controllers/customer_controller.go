package controller

import (
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/labstack/echo/v4"
)

type CustomerController struct {
	customerService service.CustomerService
}

func NewCustomerController(customerService service.CustomerService) *CustomerController {
	return &CustomerController{
		customerService,
	}
}

func (c *CustomerController) GetAll(context echo.Context) error {
	customers, err := c.customerService.GetAll(context.Request().Context())
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	return context.JSON(_http.StatusOK, customers)
}
