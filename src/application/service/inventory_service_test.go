package service

import (
	"context"
	"errors"
	"testing"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

func TestInventoryServiceDoTransaction(t *testing.T) {
	useCase := &stubInventoryUseCase{}
	service := &inventoryService{inventoryUseCase: useCase}
	req := request.CreateInventoryTransactionRequest{
		Type:                   domain.InventoryTransactionTypeIn,
		InventoryDestinationId: 1,
		Skus:                   []request.CreateInventoryTransactionSkusRequest{{SkuId: 1, Quantity: 2}},
	}

	if err := service.DoTransaction(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(useCase.receivedInput.Skus) != 1 || useCase.receivedInput.Skus[0].Quantity != 2 {
		t.Fatalf("expected sku quantity to be passed")
	}
}

func TestInventoryServiceGetAllItems(t *testing.T) {
	repo := &stubInventoryItemRepository{getAll: []output.GetInventoryItemsOutput{{}}}
	service := &inventoryService{inventoryItemRepository: repo}

	items, err := service.GetAllInventoryItems(context.Background())
	if err != nil || len(items) != 1 {
		t.Fatalf("unexpected items result")
	}
}

func TestInventoryServiceGetByInventoryId(t *testing.T) {
	repo := &stubInventoryItemRepository{getByInventory: []output.GetInventoryItemsOutput{{}}}
	service := &inventoryService{inventoryItemRepository: repo}

	items, err := service.GetInventoryItemsByInventoryId(context.Background(), 1)
	if err != nil || len(items) != 1 {
		t.Fatalf("unexpected items result")
	}
}

func TestInventoryServiceGetTransactions(t *testing.T) {
	repo := &stubInventoryTransactionRepository{getAll: []output.GetInventoryTransactionsOutput{{}}}
	service := &inventoryService{inventoryTransactionRepo: repo}

	txs, err := service.GetAllInventoryTransactions(context.Background())
	if err != nil || len(txs) != 1 {
		t.Fatalf("unexpected transactions result")
	}
}

func TestInventoryServiceGetInventories(t *testing.T) {
	repo := &stubInventoryRepository{getAll: []domain.Inventory{{Id: 1}}}
	service := &inventoryService{inventoryRepository: repo}

	inventories, err := service.GetAllInventories(context.Background())
	if err != nil || len(inventories) != 1 {
		t.Fatalf("unexpected inventories result")
	}
}

func TestInventoryServiceDoTransactionValidationError(t *testing.T) {
	service := &inventoryService{}
	if err := service.DoTransaction(context.Background(), request.CreateInventoryTransactionRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestInventoryServiceDoTransactionUseCaseError(t *testing.T) {
	useCase := &stubInventoryUseCase{err: errors.New("fail")}
	service := &inventoryService{inventoryUseCase: useCase}
	req := request.CreateInventoryTransactionRequest{
		Type:                   domain.InventoryTransactionTypeIn,
		InventoryDestinationId: 1,
		Skus:                   []request.CreateInventoryTransactionSkusRequest{{SkuId: 1, Quantity: 1}},
	}
	if err := service.DoTransaction(context.Background(), req); err == nil || err.Error() != "fail" {
		t.Fatalf("expected error from use case")
	}
}
