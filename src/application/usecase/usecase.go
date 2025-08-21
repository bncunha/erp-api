package usecase

import (
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type ApplicationUseCase struct {
	repositories     *repository.Repository
	InventoryUseCase inventory_usecase.InventoryUseCase
}

func NewApplicationUseCase(repositories *repository.Repository) *ApplicationUseCase {
	return &ApplicationUseCase{repositories: repositories}
}

func (s *ApplicationUseCase) SetupUseCases() {
	s.InventoryUseCase = inventory_usecase.NewInventoryUseCase(s.repositories, s.repositories.InventoryRepository, s.repositories.InventoryItemRepository, s.repositories.InventoryTransactionRepository, s.repositories.SkuRepository)
}
