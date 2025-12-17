package controller

import (
    _http "net/http"

    "github.com/bncunha/erp-api/src/api/http"
    request "github.com/bncunha/erp-api/src/api/requests"
    "github.com/bncunha/erp-api/src/application/service"
    "github.com/labstack/echo/v4"
)

type CompanyController struct {
    companyService service.CompanyService
}

func NewCompanyController(companyService service.CompanyService) *CompanyController {
    return &CompanyController{companyService}
}

func (c *CompanyController) Create(context echo.Context) error {
    var companyRequest request.CreateCompanyRequest
    if err := context.Bind(&companyRequest); err != nil {
        return context.JSON(_http.StatusBadRequest, http.HandleError(err))
    }

    if err := c.companyService.Create(context.Request().Context(), companyRequest); err != nil {
        return context.JSON(_http.StatusBadRequest, http.HandleError(err))
    }

    return context.JSON(_http.StatusCreated, nil)
}
