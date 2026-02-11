package sales_usecase

import (
	"context"

	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type SalesUseCase interface {
	DoSale(ctx context.Context, input DoSaleInput) error
	DoReturn(ctx context.Context, input DoReturnInput) error
}

type salesUseCase struct {
	userRepository          domain.UserRepository
	customerRepository      domain.CustomerRepository
	skuRepository           domain.SkuRepository
	saleRepository          domain.SalesRepository
	inventoryUseCase        inventory_usecase.InventoryUseCase
	inventoryRepository     domain.InventoryRepository
	inventoryItemRepository domain.InventoryItemRepository
	repository              *repository.Repository
}

func NewSalesUseCase(userRepository domain.UserRepository,
	customerRepository domain.CustomerRepository,
	skuRepository domain.SkuRepository,
	saleRepository domain.SalesRepository,
	inventoryUseCase inventory_usecase.InventoryUseCase,
	inventoryRepository domain.InventoryRepository,
	inventoryItemRepository domain.InventoryItemRepository,
	repository *repository.Repository) SalesUseCase {
	return &salesUseCase{
		userRepository:          userRepository,
		customerRepository:      customerRepository,
		skuRepository:           skuRepository,
		saleRepository:          saleRepository,
		inventoryUseCase:        inventoryUseCase,
		inventoryRepository:     inventoryRepository,
		repository:              repository,
		inventoryItemRepository: inventoryItemRepository,
	}
}
