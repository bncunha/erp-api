package service

import (
	"context"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type InventoryService interface {
	DoTransaction(ctx context.Context, request request.CreateInventoryTransactionRequest) error
	GetAllInventory(ctx context.Context) ([]output.GetInventoryItemsOutput, error)
}

type inventoryService struct {
	inventoryUseCase        inventory_usecase.InventoryUseCase
	inventoryItemRepository repository.InventoryItemRepository
}

func NewInventoryService(inventoryUseCase inventory_usecase.InventoryUseCase, inventoryItemRepository repository.InventoryItemRepository) InventoryService {
	return &inventoryService{inventoryUseCase, inventoryItemRepository}
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

func (s *inventoryService) GetAllInventory(ctx context.Context) ([]output.GetInventoryItemsOutput, error) {
	return s.inventoryItemRepository.GetAll(ctx)
}
