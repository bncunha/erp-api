package service

import (
	"context"
	"errors"
	"testing"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

func TestSkuServiceCreate(t *testing.T) {
	skuRepo := &stubSkuRepository{}
	inventoryUseCase := &stubInventoryUseCase{}
	productRepo := &stubProductRepository{getById: domain.Product{Id: 1}}
	service := &skuService{skuRepository: skuRepo, inventoryUseCase: inventoryUseCase, productRepository: productRepo, txManager: &stubRepository{}}

	qty := 5.0
	dest := int64(1)
	cost := 10.0
	price := 15.0
	req := request.CreateSkuRequest{Code: "code", Color: "red", Size: "M", Quantity: &qty, DestinationId: &dest, Cost: &cost, Price: price}

	if err := service.Create(context.Background(), req, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skuRepo.created) == 0 {
		t.Fatalf("expected sku creation")
	}
	if inventoryUseCase.receivedInput.Type != domain.InventoryTransactionTypeIn {
		t.Fatalf("expected inventory transaction to be created")
	}
}

func TestSkuServiceCreateDuplicated(t *testing.T) {
	skuRepo := &stubSkuRepository{createErr: errors.New("duplicate key value violates unique constraint")}
	service := &skuService{skuRepository: skuRepo, productRepository: &stubProductRepository{getById: domain.Product{Id: 1}}, txManager: &stubRepository{}}
	cost := 1.0
	price := 2.0
	req := request.CreateSkuRequest{Code: "code", Color: "red", Size: "M", Cost: &cost, Price: price}

	err := service.Create(context.Background(), req, 1)
	if err == nil || err.Error() != "Código já cadastrado!" {
		t.Fatalf("expected duplicated error, got %v", err)
	}
}

func TestSkuServiceUpdate(t *testing.T) {
	skuRepo := &stubSkuRepository{}
	service := &skuService{skuRepository: skuRepo}
	cost := 1.0
	price := 2.0
	req := request.EditSkuRequest{CreateSkuRequest: request.CreateSkuRequest{Code: "code", Color: "red", Size: "M", Cost: &cost, Price: price}}

	if err := service.Update(context.Background(), req, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skuRepo.created) == 0 {
		t.Fatalf("expected update to occur")
	}
}

func TestSkuServiceUpdateDuplicated(t *testing.T) {
	skuRepo := &stubSkuRepository{updateErr: errors.New("duplicate key value violates unique constraint")}
	service := &skuService{skuRepository: skuRepo}
	cost := 1.0
	price := 2.0
	req := request.EditSkuRequest{CreateSkuRequest: request.CreateSkuRequest{Code: "code", Color: "red", Size: "M", Cost: &cost, Price: price}}

	err := service.Update(context.Background(), req, 1)
	if err == nil || err.Error() != "Código já cadastrado!" {
		t.Fatalf("expected duplicated error")
	}
}

func TestSkuServiceCreateInventoryError(t *testing.T) {
	skuRepo := &stubSkuRepository{}
	inventoryUseCase := &stubInventoryUseCase{err: errors.New("fail")}
	productRepo := &stubProductRepository{getById: domain.Product{Id: 1}}
	service := &skuService{skuRepository: skuRepo, inventoryUseCase: inventoryUseCase, productRepository: productRepo, txManager: &stubRepository{}}
	qty := 1.0
	dest := int64(1)
	cost := 1.0
	price := 2.0
	req := request.CreateSkuRequest{Code: "code", Color: "red", Size: "M", Quantity: &qty, DestinationId: &dest, Cost: &cost, Price: price}

	err := service.Create(context.Background(), req, 1)
	if err == nil || err.Error() != "Operação realizada parcialmente! Erro ao atualizar a quantidade de itens no estoque!" {
		t.Fatalf("expected partial operation error")
	}
}

func TestSkuServiceCreateValidationError(t *testing.T) {
	service := &skuService{}
	if err := service.Create(context.Background(), request.CreateSkuRequest{}, 1); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestSkuServiceUpdateValidationError(t *testing.T) {
	service := &skuService{}
	if err := service.Update(context.Background(), request.EditSkuRequest{}, 1); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestSkuServiceCreateRepositoryError(t *testing.T) {
	skuRepo := &stubSkuRepository{createErr: errors.New("other")}
	service := &skuService{skuRepository: skuRepo, productRepository: &stubProductRepository{getById: domain.Product{Id: 1}}, txManager: &stubRepository{}}
	cost := 1.0
	price := 2.0
	req := request.CreateSkuRequest{Code: "code", Color: "red", Size: "M", Cost: &cost, Price: price}

	if err := service.Create(context.Background(), req, 1); err == nil || err.Error() != "other" {
		t.Fatalf("expected repository error")
	}
}

func TestSkuServiceUpdateRepositoryError(t *testing.T) {
	skuRepo := &stubSkuRepository{updateErr: errors.New("other")}
	service := &skuService{skuRepository: skuRepo}
	cost := 1.0
	price := 2.0
	req := request.EditSkuRequest{CreateSkuRequest: request.CreateSkuRequest{Code: "code", Color: "red", Size: "M", Cost: &cost, Price: price}}

	if err := service.Update(context.Background(), req, 1); err == nil || err.Error() != "other" {
		t.Fatalf("expected repository error")
	}
}

func TestSkuServiceGetById(t *testing.T) {
	skuRepo := &stubSkuRepository{getById: domain.Sku{Id: 1}}
	service := &skuService{skuRepository: skuRepo}

	sku, err := service.GetById(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sku.Id != 1 {
		t.Fatalf("expected sku")
	}
}

func TestSkuServiceGetByIdError(t *testing.T) {
	skuRepo := &stubSkuRepository{getByIdErr: errors.New("fail")}
	service := &skuService{skuRepository: skuRepo}
	if _, err := service.GetById(context.Background(), 1); err == nil {
		t.Fatalf("expected error")
	}
}

func TestSkuServiceGetAll(t *testing.T) {
	skuRepo := &stubSkuRepository{getAll: []domain.Sku{{Id: 1}}}
	service := &skuService{skuRepository: skuRepo}

	ctx := context.WithValue(context.Background(), constants.ROLE_KEY, string(domain.UserRoleAdmin))
	skus, err := service.GetAll(ctx, GetSkusFilters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skus) != 1 {
		t.Fatalf("expected skus")
	}
}

func TestSkuServiceInactivate(t *testing.T) {
	skuRepo := &stubSkuRepository{}
	service := &skuService{skuRepository: skuRepo}

	if err := service.Inactivate(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSkuServiceGetAllError(t *testing.T) {
	skuRepo := &stubSkuRepository{getAllErr: errors.New("fail")}
	service := &skuService{skuRepository: skuRepo}
	ctx := context.WithValue(context.Background(), constants.ROLE_KEY, string(domain.UserRoleAdmin))
	if _, err := service.GetAll(ctx, GetSkusFilters{}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestSkuServiceGetAllAdminFilter(t *testing.T) {
	skuRepo := &stubSkuRepository{}
	service := &skuService{skuRepository: skuRepo}

	ctx := context.WithValue(context.Background(), constants.ROLE_KEY, string(domain.UserRoleAdmin))
	sellerId := 10.0

	if _, err := service.GetAll(ctx, GetSkusFilters{SellerId: &sellerId}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if skuRepo.getAllInput.SellerId == nil || *skuRepo.getAllInput.SellerId != sellerId {
		t.Fatalf("expected seller filter to be forwarded")
	}
}

func TestSkuServiceGetAllNonAdminIgnoresFilter(t *testing.T) {
	skuRepo := &stubSkuRepository{}
	service := &skuService{skuRepository: skuRepo}

	ctx := context.WithValue(context.Background(), constants.ROLE_KEY, string(domain.UserRoleReseller))
	userId := 5.0
	ctx = context.WithValue(ctx, constants.USERID_KEY, userId)
	sellerId := 2.0

	if _, err := service.GetAll(ctx, GetSkusFilters{SellerId: &sellerId}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if skuRepo.getAllInput.SellerId == nil || *skuRepo.getAllInput.SellerId != userId {
		t.Fatalf("expected seller filter to use logged user id")
	}
}

func TestSkuServiceInactivateError(t *testing.T) {
	skuRepo := &stubSkuRepository{inactivateErr: errors.New("fail")}
	service := &skuService{skuRepository: skuRepo}
	if err := service.Inactivate(context.Background(), 1); err == nil || err.Error() != "fail" {
		t.Fatalf("expected error")
	}
}

func TestSkuServiceCreateProductLookupError(t *testing.T) {
	service := &skuService{skuRepository: &stubSkuRepository{}, productRepository: &stubProductRepository{getByIdErr: errors.New("fail")}, txManager: &stubRepository{}}
	cost := 1.0
	price := 2.0
	req := request.CreateSkuRequest{Code: "code", Color: "red", Size: "M", Cost: &cost, Price: price}
	if err := service.Create(context.Background(), req, 1); err == nil || err.Error() != "fail" {
		t.Fatalf("expected product lookup error")
	}
}

func TestSkuServiceCreateTxBeginError(t *testing.T) {
	service := &skuService{
		skuRepository:     &stubSkuRepository{},
		productRepository: &stubProductRepository{getById: domain.Product{Id: 1}},
		txManager:         &stubTxManager{err: errors.New("begin fail")},
	}
	cost := 1.0
	price := 2.0
	req := request.CreateSkuRequest{Code: "code", Color: "red", Size: "M", Cost: &cost, Price: price}
	if err := service.Create(context.Background(), req, 1); err == nil || err.Error() != "begin fail" {
		t.Fatalf("expected begin tx error")
	}
}

func TestSkuServiceCreateCommitsTransaction(t *testing.T) {
	sqlTx, fakeTx, cleanup := newTestSQLTx()
	defer cleanup()

	service := &skuService{
		skuRepository:     &stubSkuRepository{},
		productRepository: &stubProductRepository{getById: domain.Product{Id: 1}},
		txManager:         &stubTxManager{tx: sqlTx},
	}
	cost := 1.0
	price := 2.0
	req := request.CreateSkuRequest{Code: "code", Color: "red", Size: "M", Cost: &cost, Price: price}
	if err := service.Create(context.Background(), req, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fakeTx.committed {
		t.Fatalf("expected commit to be called")
	}
	if fakeTx.rolledBack {
		t.Fatalf("did not expect rollback on success")
	}
}

func TestSkuServiceCreateRollsBackOnInventoryError(t *testing.T) {
	sqlTx, fakeTx, cleanup := newTestSQLTx()
	defer cleanup()

	qty := 1.0
	dest := int64(1)
	cost := 1.0
	price := 2.0
	service := &skuService{
		skuRepository:     &stubSkuRepository{},
		productRepository: &stubProductRepository{getById: domain.Product{Id: 1}},
		inventoryUseCase:  &stubInventoryUseCase{err: errors.New("fail")},
		txManager:         &stubTxManager{tx: sqlTx},
	}
	req := request.CreateSkuRequest{Code: "code", Color: "red", Size: "M", Quantity: &qty, DestinationId: &dest, Cost: &cost, Price: price}
	err := service.Create(context.Background(), req, 1)
	if err == nil || err.Error() != "Operação realizada parcialmente! Erro ao atualizar a quantidade de itens no estoque!" {
		t.Fatalf("expected partial operation error")
	}
	if !fakeTx.rolledBack {
		t.Fatalf("expected rollback when inventory update fails")
	}
	if fakeTx.committed {
		t.Fatalf("did not expect commit on failure")
	}
}

func TestSkuServiceGetInventory(t *testing.T) {
	skuRepo := &stubSkuRepository{getById: domain.Sku{Id: 3}}
	itemRepo := &stubInventoryItemRepository{getBySku: []output.GetSkuInventoryOutput{{InventoryId: 3}}}
	service := &skuService{skuRepository: skuRepo, inventoryItemRepository: itemRepo}

	items, err := service.GetInventory(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0].InventoryId != 3 {
		t.Fatalf("expected inventory items")
	}
}

func TestSkuServiceGetInventorySkuNotFound(t *testing.T) {
	skuRepo := &stubSkuRepository{getByIdErr: errors.New("not found")}
	service := &skuService{skuRepository: skuRepo, inventoryItemRepository: &stubInventoryItemRepository{}}

	if _, err := service.GetInventory(context.Background(), 3); err == nil {
		t.Fatalf("expected sku lookup error")
	}
}

func TestSkuServiceGetTransactions(t *testing.T) {
	skuRepo := &stubSkuRepository{getById: domain.Sku{Id: 5}}
	trxRepo := &stubInventoryTransactionRepository{getBySkuId: []output.GetInventoryTransactionsOutput{{Id: 7}}}
	service := &skuService{skuRepository: skuRepo, inventoryTransactionRepository: trxRepo}

	items, err := service.GetTransactions(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0].Id != 7 {
		t.Fatalf("expected transaction items")
	}
}

func TestSkuServiceGetTransactionsSkuNotFound(t *testing.T) {
	skuRepo := &stubSkuRepository{getByIdErr: errors.New("not found")}
	service := &skuService{skuRepository: skuRepo, inventoryTransactionRepository: &stubInventoryTransactionRepository{}}

	if _, err := service.GetTransactions(context.Background(), 5); err == nil {
		t.Fatalf("expected sku lookup error")
	}
}
