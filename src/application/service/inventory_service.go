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
	GetAllInventoryItems(ctx context.Context) ([]output.GetInventoryItemsOutput, error)
	GetInventoryItemsByInventoryId(ctx context.Context, id int64) ([]output.GetInventoryItemsOutput, error)
	GetAllInventoryTransactions(ctx context.Context) ([]output.GetInventoryTransactionsOutput, error)
	GetAllInventories(ctx context.Context) ([]domain.Inventory, error)
}

type inventoryService struct {
	inventoryUseCase         inventory_usecase.InventoryUseCase
	inventoryItemRepository  repository.InventoryItemRepository
	inventoryTransactionRepo repository.InventoryTransactionRepository
	inventoryRepository      repository.InventoryRepository
}

func NewInventoryService(inventoryUseCase inventory_usecase.InventoryUseCase, inventoryItemRepository repository.InventoryItemRepository, inventoryTransactionRepo repository.InventoryTransactionRepository, inventoryRepository repository.InventoryRepository) InventoryService {
	return &inventoryService{inventoryUseCase, inventoryItemRepository, inventoryTransactionRepo, inventoryRepository}
}

func (s *inventoryService) DoTransaction(ctx context.Context, request request.CreateInventoryTransactionRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	inputSkus := make([]inventory_usecase.DoTransactionSkusInput, 0)
	for _, sku := range request.Skus {
		inputSkus = append(inputSkus, inventory_usecase.DoTransactionSkusInput{
			SkuId:    sku.SkuId,
			Quantity: sku.Quantity,
		})
	}

	return s.inventoryUseCase.DoTransaction(ctx, inventory_usecase.DoTransactionInput{
		Type:                   domain.InventoryTransactionType(request.Type),
		InventoryOriginId:      request.InventoryOriginId,
		InventoryDestinationId: request.InventoryDestinationId,
		Justification:          request.Justification,
		Skus:                   inputSkus,
	})
}

func (s *inventoryService) GetAllInventoryItems(ctx context.Context) ([]output.GetInventoryItemsOutput, error) {
	return s.inventoryItemRepository.GetAll(ctx)
}

func (s *inventoryService) GetInventoryItemsByInventoryId(ctx context.Context, id int64) ([]output.GetInventoryItemsOutput, error) {
	return s.inventoryItemRepository.GetByInventoryId(ctx, id)
}

func (s *inventoryService) GetAllInventoryTransactions(ctx context.Context) ([]output.GetInventoryTransactionsOutput, error) {
	return s.inventoryTransactionRepo.GetAll(ctx)
}

func (s *inventoryService) GetAllInventories(ctx context.Context) ([]domain.Inventory, error) {
	return s.inventoryRepository.GetAll(ctx)
}
