package repository

import (
	"context"
	"database/sql"
)

type Repository struct {
	db                             *sql.DB
	ProductRepository              ProductRepository
	CategoryRepository             CategoryRepository
	SkuRepository                  SkuRepository
	UserRepository                 UserRepository
	InventoryRepository            InventoryRepository
	InventoryItemRepository        InventoryItemRepository
	InventoryTransactionRepository InventoryTransactionRepository
	SalesRepository                SalesRepository
	CustomerRepository             CustomerRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SetupRepositories() {
	r.ProductRepository = NewProductRepository(r.db)
	r.CategoryRepository = NewCategoryRepository(r.db)
	r.SkuRepository = NewSkuRepository(r.db)
	r.UserRepository = NewUserRepository(r.db)
	r.InventoryRepository = NewInventoryRepository(r.db)
	r.InventoryItemRepository = NewInventoryItemRepository(r.db)
	r.InventoryTransactionRepository = NewInventoryTransactionRepository(r.db, r.InventoryItemRepository)
	r.SalesRepository = NewSalesRepository(r.db)
	r.CustomerRepository = NewCustomerRepository(r.db)
}

func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
