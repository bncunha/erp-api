package controller

import (
    "errors"
    _http "net/http"

    "github.com/bncunha/erp-api/src/api/http"
    request "github.com/bncunha/erp-api/src/api/requests"
    "github.com/bncunha/erp-api/src/api/viewmodel"
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

func (c *CustomerController) Create(context echo.Context) error {
    var req request.CreateCustomerRequest
    if err := context.Bind(&req); err != nil {
        return context.JSON(_http.StatusBadRequest, http.HandleError(errors.New("parametros invalidos")))
    }

    id, err := c.customerService.Create(context.Request().Context(), req)
    if err != nil {
        return context.JSON(_http.StatusBadRequest, http.HandleError(err))
    }
    return context.JSON(_http.StatusCreated, id)
}

func (c *CustomerController) GetAll(context echo.Context) error {
    customers, err := c.customerService.GetAll(context.Request().Context())
    if err != nil {
        return context.JSON(_http.StatusBadRequest, http.HandleError(err))
    }
    return context.JSON(_http.StatusOK, viewmodel.ToCustomerViewModel(customers))
}
