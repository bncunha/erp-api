package controller

import (
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/labstack/echo/v4"
)

type SkuController struct {
	skuService service.SkuService
}

func NewSkuController(skuService service.SkuService) *SkuController {
	return &SkuController{skuService}
}

func (c *SkuController) Create(context echo.Context) error {
	productId := context.Param("id")
	var skuRequest request.CreateSkuRequest
	if err := context.Bind(&skuRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	err := c.skuService.Create(context.Request().Context(), skuRequest, helper.ParseInt64(productId))
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusCreated, nil)
}

func (c *SkuController) Edit(context echo.Context) error {
	skuId := helper.ParseInt64(context.Param("id"))
	var skuRequest request.EditSkuRequest

	if err := context.Bind(&skuRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	err := c.skuService.Update(context.Request().Context(), skuRequest, skuId)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, nil)
}