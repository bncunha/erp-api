package controller

import (
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/api/viewmodel"
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

func (c *InventoryController) GetAllInventory(context echo.Context) error {
	inventoryItems, err := c.inventoryService.GetAllInventory(context.Request().Context())
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	var inventoryItemViewModels []viewmodel.GetInventoryItemsViewModel
	for _, inventoryItem := range inventoryItems {
		inventoryItemViewModels = append(inventoryItemViewModels, viewmodel.ToGetInventoryItemsViewModel(inventoryItem))
	}

	return context.JSON(_http.StatusOK, inventoryItemViewModels)
}

func (c *InventoryController) GetAllInventoryTransactions(context echo.Context) error {
	inventoryTransactions, err := c.inventoryService.GetAllInventoryTransactions(context.Request().Context())
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	var inventoryTransactionViewModels []viewmodel.GetInventoryTransactionsViewModel
	for _, inventoryTransaction := range inventoryTransactions {
		inventoryTransactionViewModels = append(inventoryTransactionViewModels, viewmodel.ToGetInventoryTransactionsViewModel(inventoryTransaction))
	}

	return context.JSON(_http.StatusOK, inventoryTransactionViewModels)
}
