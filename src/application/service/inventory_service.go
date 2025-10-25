package service

import (
	"context"
	"database/sql"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type transactionManager interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type InventoryService interface {
	DoTransaction(ctx context.Context, request request.CreateInventoryTransactionRequest) error
	GetAllInventoryItems(ctx context.Context) ([]output.GetInventoryItemsOutput, error)
	GetInventoryItemsByInventoryId(ctx context.Context, id int64) ([]output.GetInventoryItemsOutput, error)
	GetInventoryTransactionsByInventoryId(ctx context.Context, id int64) ([]output.GetInventoryTransactionsOutput, error)
	GetAllInventories(ctx context.Context) ([]domain.Inventory, error)
	GetInventoriesSummary(ctx context.Context) ([]output.GetInventorySummaryOutput, error)
	GetInventorySummaryById(ctx context.Context, id int64) (output.GetInventorySummaryByIdOutput, error)
}

type inventoryService struct {
	inventoryUseCase         inventory_usecase.InventoryUseCase
	inventoryItemRepository  repository.InventoryItemRepository
	inventoryTransactionRepo repository.InventoryTransactionRepository
	inventoryRepository      repository.InventoryRepository
	txManager                transactionManager
}

func NewInventoryService(inventoryUseCase inventory_usecase.InventoryUseCase, inventoryItemRepository repository.InventoryItemRepository, inventoryTransactionRepo repository.InventoryTransactionRepository, inventoryRepository repository.InventoryRepository, txManager transactionManager) InventoryService {
	return &inventoryService{inventoryUseCase, inventoryItemRepository, inventoryTransactionRepo, inventoryRepository, txManager}
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

	var tx *sql.Tx
	if s.txManager != nil {
		tx, err = s.txManager.BeginTx(ctx)
		if err != nil {
			return err
		}
		defer func() {
			if err != nil && tx != nil {
				tx.Rollback()
			}
		}()
	}

	err = s.inventoryUseCase.DoTransaction(ctx, tx, inventory_usecase.DoTransactionInput{
		Type:                   domain.InventoryTransactionType(request.Type),
		InventoryOriginId:      request.InventoryOriginId,
		InventoryDestinationId: request.InventoryDestinationId,
		Justification:          request.Justification,
		Skus:                   inputSkus,
	})
	if err != nil {
		return err
	}
	if tx != nil {
		return tx.Commit()
	}
	return nil
}

func (s *inventoryService) GetAllInventoryItems(ctx context.Context) ([]output.GetInventoryItemsOutput, error) {
	return s.inventoryItemRepository.GetAll(ctx)
}

func (s *inventoryService) GetInventoryItemsByInventoryId(ctx context.Context, id int64) ([]output.GetInventoryItemsOutput, error) {
	return s.inventoryItemRepository.GetByInventoryId(ctx, id)
}

func (s *inventoryService) GetInventoryTransactionsByInventoryId(ctx context.Context, id int64) ([]output.GetInventoryTransactionsOutput, error) {
	return s.inventoryTransactionRepo.GetByInventoryId(ctx, id)
}

func (s *inventoryService) GetAllInventories(ctx context.Context) ([]domain.Inventory, error) {
	return s.inventoryRepository.GetAll(ctx)
}

func (s *inventoryService) GetInventoriesSummary(ctx context.Context) ([]output.GetInventorySummaryOutput, error) {
	return s.inventoryRepository.GetSummary(ctx)
}

func (s *inventoryService) GetInventorySummaryById(ctx context.Context, id int64) (output.GetInventorySummaryByIdOutput, error) {
	return s.inventoryRepository.GetSummaryById(ctx, id)
}
