package controller

import (
	"errors"
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/api/viewmodel"
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

func (c *SkuController) GetById(context echo.Context) error {
	skuId := helper.ParseInt64(context.Param("id"))

	sku, err := c.skuService.GetById(context.Request().Context(), skuId)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, viewmodel.ToSkuViewModel(sku))
}

func (c *SkuController) GetAll(context echo.Context) error {
	var sellerId *float64
	if sellerIdParam := context.QueryParam("seller_id"); sellerIdParam != "" {
		parsedSellerId, err := helper.ParseFloat(sellerIdParam)
		if err != nil {
			return context.JSON(_http.StatusBadRequest, http.HandleError(errors.New("seller_id inválido")))
		}
		sellerId = &parsedSellerId
	}

	skus, err := c.skuService.GetAll(context.Request().Context(), service.GetSkusFilters{SellerId: sellerId})
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	var skuViewModels []viewmodel.SkuViewModel = make([]viewmodel.SkuViewModel, 0)
	for _, sku := range skus {
		skuViewModels = append(skuViewModels, viewmodel.ToSkuViewModel(sku))
	}

	return context.JSON(_http.StatusOK, skuViewModels)
}

func (c *SkuController) Inactivate(context echo.Context) error {
	id := helper.ParseInt64(context.Param("id"))

	err := c.skuService.Inactivate(context.Request().Context(), id)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, nil)
}

func (c *SkuController) GetInventory(context echo.Context) error {
	skuId := helper.ParseInt64(context.Param("id"))

	items, err := c.skuService.GetInventory(context.Request().Context(), skuId)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	viewModels := make([]viewmodel.SkuInventoryViewModel, 0, len(items))
	for _, item := range items {
		viewModels = append(viewModels, viewmodel.ToSkuInventoryViewModel(item))
	}

	return context.JSON(_http.StatusOK, viewModels)
}

func (c *SkuController) GetTransactions(context echo.Context) error {
	skuId := helper.ParseInt64(context.Param("id"))

	transactions, err := c.skuService.GetTransactions(context.Request().Context(), skuId)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	viewModels := make([]viewmodel.SkuTransactionViewModel, 0, len(transactions))
	for _, transaction := range transactions {
		viewModels = append(viewModels, viewmodel.ToSkuTransactionViewModel(transaction))
	}

	return context.JSON(_http.StatusOK, viewModels)
}
