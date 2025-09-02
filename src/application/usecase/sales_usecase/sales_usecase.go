package sales_usecase

import (
	"context"

	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type SalesUseCase interface {
	DoSale(ctx context.Context, input DoSaleInput) error
}

type salesUseCase struct {
	userRepository      repository.UserRepository
	customerRepository  repository.CustomerRepository
	skuRepository       repository.SkuRepository
	saleRepository      repository.SalesRepository
	inventoryUseCase    inventory_usecase.InventoryUseCase
	inventoryRepository repository.InventoryRepository
	repository          *repository.Repository
}

func NewSalesUseCase(userRepository repository.UserRepository,
	customerRepository repository.CustomerRepository,
	skuRepository repository.SkuRepository,
	saleRepository repository.SalesRepository,
	inventoryUseCase inventory_usecase.InventoryUseCase,
	inventoryRepository repository.InventoryRepository,
	repository *repository.Repository) SalesUseCase {
	return &salesUseCase{
		userRepository:      userRepository,
		customerRepository:  customerRepository,
		skuRepository:       skuRepository,
		saleRepository:      saleRepository,
		inventoryUseCase:    inventoryUseCase,
		inventoryRepository: inventoryRepository,
		repository:          repository,
	}
}
