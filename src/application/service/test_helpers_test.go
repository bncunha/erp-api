package service

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/domain"
)

type stubProductRepository struct {
	created    domain.Product
	createErr  error
	editErr    error
	getById    domain.Product
	getByIdErr error
	getAll     []output.GetAllProductsOutput
	getAllErr  error
}

func (s *stubProductRepository) Create(ctx context.Context, product domain.Product) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = product
	return 1, nil
}

func (s *stubProductRepository) Edit(ctx context.Context, product domain.Product, id int64) (int64, error) {
	if s.editErr != nil {
		return 0, s.editErr
	}
	s.created = product
	return id, nil
}

func (s *stubProductRepository) GetById(ctx context.Context, id int64) (domain.Product, error) {
	return s.getById, s.getByIdErr
}

func (s *stubProductRepository) GetAll(ctx context.Context) ([]output.GetAllProductsOutput, error) {
	return s.getAll, s.getAllErr
}

func (s *stubProductRepository) Inactivate(ctx context.Context, id int64) error {
	s.created = domain.Product{Id: id}
	return nil
}

type stubCategoryRepository struct {
	created      domain.Category
	createErr    error
	getById      domain.Category
	getByIdErr   error
	getByName    domain.Category
	getByNameErr error
	updateErr    error
	deleteErr    error
	getAll       []domain.Category
	getAllErr    error
}

func (s *stubCategoryRepository) Create(ctx context.Context, category domain.Category) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = category
	return 1, nil
}

func (s *stubCategoryRepository) GetById(ctx context.Context, id int64) (domain.Category, error) {
	return s.getById, s.getByIdErr
}

func (s *stubCategoryRepository) GetByName(ctx context.Context, name string) (domain.Category, error) {
	return s.getByName, s.getByNameErr
}

func (s *stubCategoryRepository) Update(ctx context.Context, category domain.Category) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.created = category
	return nil
}

func (s *stubCategoryRepository) Delete(ctx context.Context, id int64) error {
	return s.deleteErr
}

func (s *stubCategoryRepository) GetAll(ctx context.Context) ([]domain.Category, error) {
	return s.getAll, s.getAllErr
}

type stubSkuRepository struct {
	created         []domain.Sku
	createErr       error
	createManyErr   error
	updateErr       error
	getById         domain.Sku
	getByIdErr      error
	getByProduct    []domain.Sku
	getByProductErr error
	getAll          []domain.Sku
	getAllErr       error
	inactivateErr   error
}

func (s *stubSkuRepository) Create(ctx context.Context, sku domain.Sku, productId int64) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = []domain.Sku{sku}
	return 1, nil
}

func (s *stubSkuRepository) CreateMany(ctx context.Context, skus []domain.Sku, productId int64) ([]int64, error) {
	if s.createManyErr != nil {
		return nil, s.createManyErr
	}
	s.created = append(s.created, skus...)
	return []int64{1}, nil
}

func (s *stubSkuRepository) GetByProductId(ctx context.Context, productId int64) ([]domain.Sku, error) {
	return s.getByProduct, s.getByProductErr
}

func (s *stubSkuRepository) Update(ctx context.Context, sku domain.Sku) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.created = []domain.Sku{sku}
	return nil
}

func (s *stubSkuRepository) GetById(ctx context.Context, id int64) (domain.Sku, error) {
	return s.getById, s.getByIdErr
}

func (s *stubSkuRepository) GetByManyIds(ctx context.Context, ids []int64) ([]domain.Sku, error) {
	return nil, nil
}

func (s *stubSkuRepository) GetAll(ctx context.Context) ([]domain.Sku, error) {
	return s.getAll, s.getAllErr
}

func (s *stubSkuRepository) Inactivate(ctx context.Context, id int64) error {
	return s.inactivateErr
}

type stubUserRepository struct {
	created          domain.User
	createErr        error
	updateErr        error
	getById          domain.User
	getByIdErr       error
	getAll           []domain.User
	getAllErr        error
	inactivateErr    error
	getByUsername    domain.User
	getByUsernameErr error
}

func (s *stubUserRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = user
	return 1, nil
}

func (s *stubUserRepository) Update(ctx context.Context, user domain.User) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.created = user
	return nil
}

func (s *stubUserRepository) GetById(ctx context.Context, id int64) (domain.User, error) {
	return s.getById, s.getByIdErr
}

func (s *stubUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.getAll, s.getAllErr
}

func (s *stubUserRepository) Inactivate(ctx context.Context, id int64) error {
	return s.inactivateErr
}

func (s *stubUserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	return s.getByUsername, s.getByUsernameErr
}

type stubInventoryRepository struct {
	createErr    error
	created      domain.Inventory
	getAll       []domain.Inventory
	getAllErr    error
	getByUser    domain.Inventory
	getByUserErr error
	getById      domain.Inventory
	getByIdErr   error
}

func (s *stubInventoryRepository) Create(ctx context.Context, inventory domain.Inventory) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = inventory
	return 1, nil
}

func (s *stubInventoryRepository) GetById(ctx context.Context, id int64) (domain.Inventory, error) {
	return s.getById, s.getByIdErr
}

func (s *stubInventoryRepository) GetAll(ctx context.Context) ([]domain.Inventory, error) {
	return s.getAll, s.getAllErr
}

func (s *stubInventoryRepository) GetByUserId(ctx context.Context, userId int64) (domain.Inventory, error) {
	return s.getByUser, s.getByUserErr
}

type stubInventoryItemRepository struct {
	getAll            []output.GetInventoryItemsOutput
	getAllErr         error
	getByInventory    []output.GetInventoryItemsOutput
	getByInventoryErr error
}

func (s *stubInventoryItemRepository) GetAll(ctx context.Context) ([]output.GetInventoryItemsOutput, error) {
	return s.getAll, s.getAllErr
}

func (s *stubInventoryItemRepository) GetByInventoryId(ctx context.Context, id int64) ([]output.GetInventoryItemsOutput, error) {
	return s.getByInventory, s.getByInventoryErr
}

func (s *stubInventoryItemRepository) Create(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) (int64, error) {
	return 0, nil
}

func (s *stubInventoryItemRepository) UpdateQuantity(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) error {
	return nil
}

func (s *stubInventoryItemRepository) GetById(ctx context.Context, id int64) (domain.InventoryItem, error) {
	return domain.InventoryItem{}, nil
}

func (s *stubInventoryItemRepository) GetByIdWithTransaction(ctx context.Context, tx *sql.Tx, id int64) (domain.InventoryItem, error) {
	return domain.InventoryItem{}, nil
}

func (s *stubInventoryItemRepository) GetByManySkuIdsAndInventoryId(ctx context.Context, skuIds []int64, inventoryId int64) ([]domain.InventoryItem, error) {
	return nil, nil
}

type stubInventoryTransactionRepository struct {
	getAll    []output.GetInventoryTransactionsOutput
	getAllErr error
}

func (s *stubInventoryTransactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction domain.InventoryTransaction) (int64, error) {
	return 0, nil
}

func (s *stubInventoryTransactionRepository) GetAll(ctx context.Context) ([]output.GetInventoryTransactionsOutput, error) {
	return s.getAll, s.getAllErr
}

type stubRepository struct {
	beginErr error
}

func (s *stubRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return nil, s.beginErr
}

type stubInventoryUseCase struct {
	receivedInput inventory_usecase.DoTransactionInput
	err           error
}

func (s *stubInventoryUseCase) DoTransaction(ctx context.Context, input inventory_usecase.DoTransactionInput) error {
	s.receivedInput = input
	return s.err
}
