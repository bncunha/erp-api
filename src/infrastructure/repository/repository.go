package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/domain"
)

type Repository struct {
	db                             *sql.DB
	ProductRepository              domain.ProductRepository
	CategoryRepository             domain.CategoryRepository
	SkuRepository                  domain.SkuRepository
	UserRepository                 domain.UserRepository
	UserTokenRepository            domain.UserTokenRepository
	InventoryRepository            domain.InventoryRepository
	InventoryItemRepository        domain.InventoryItemRepository
	InventoryTransactionRepository domain.InventoryTransactionRepository
	SalesRepository                domain.SalesRepository
	CustomerRepository             domain.CustomerRepository
	CompanyRepository              domain.CompanyRepository
	AddressRepository              domain.AddressRepository
	LegalDocumentRepository        domain.LegalDocumentRepository
	LegalAcceptanceRepository      domain.LegalAcceptanceRepository
	DashboardRepository            domain.DashboardRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SetupRepositories() {
	r.ProductRepository = NewProductRepository(r.db)
	r.CategoryRepository = NewCategoryRepository(r.db)
	r.SkuRepository = NewSkuRepository(r.db)
	r.UserRepository = NewUserRepository(r.db)
	r.UserTokenRepository = NewUserTokenRepository(r.db)
	r.InventoryRepository = NewInventoryRepository(r.db)
	r.InventoryItemRepository = NewInventoryItemRepository(r.db)
	r.InventoryTransactionRepository = NewInventoryTransactionRepository(r.db, r.InventoryItemRepository)
	r.SalesRepository = NewSalesRepository(r.db)
	r.CustomerRepository = NewCustomerRepository(r.db)
	r.CompanyRepository = NewCompanyRepository(r.db)
	r.AddressRepository = NewAddressRepository(r.db)
	r.LegalDocumentRepository = NewLegalDocumentRepository(r.db)
	r.LegalAcceptanceRepository = NewLegalAcceptanceRepository(r.db)
	r.DashboardRepository = NewDashboardRepository(r.db)
}

func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
