package controller

import (
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/labstack/echo/v4"
)

type InventoryController struct {
	inventoryService service.InventoryService
}

func NewInventoryController(inventoryService service.InventoryService) *InventoryController {
	return &InventoryController{
		inventoryService,
	}
}

func (c *InventoryController) DoTransaction(context echo.Context) error {
	var inventoryTransactionRequest request.CreateInventoryTransactionRequest
	if err := context.Bind(&inventoryTransactionRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	err := c.inventoryService.DoTransaction(context.Request().Context(), inventoryTransactionRequest)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	return context.JSON(_http.StatusOK, nil)
}
