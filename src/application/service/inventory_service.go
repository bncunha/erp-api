package service

import (
	"context"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/domain"
)

type InventoryService interface {
	DoTransaction(ctx context.Context, request request.CreateInventoryTransactionRequest) error
}

type inventoryService struct {
	inventoryUseCase inventory_usecase.InventoryUseCase
}

func NewInventoryService(inventoryUseCase inventory_usecase.InventoryUseCase) InventoryService {
	return &inventoryService{inventoryUseCase}
}

func (s *inventoryService) DoTransaction(ctx context.Context, request request.CreateInventoryTransactionRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.inventoryUseCase.DoTransaction(ctx, inventory_usecase.DoTransactionInput{
		Type:                   domain.InventoryTransactionType(request.Type),
		SkuId:                  request.SkuId,
		InventoryOriginId:      request.InventoryOriginId,
		InventoryDestinationId: request.InventoryDestinationId,
		Quantity:               request.Quantity,
		Justification:          request.Justification,
	})
}
