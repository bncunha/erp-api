package usecase

import (
	"github.com/bncunha/erp-api/src/application/ports"
	emailusecase "github.com/bncunha/erp-api/src/application/usecase/email_usecase"
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/application/usecase/sales_usecase"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
	config "github.com/bncunha/erp-api/src/main"
)

type ApplicationUseCase struct {
	ports            *ports.Ports
	repositories     *repository.Repository
	config           *config.Config
	InventoryUseCase inventory_usecase.InventoryUseCase
	SalesUsecase     sales_usecase.SalesUseCase
	EmailUseCase     emailusecase.EmailUseCase
}

func NewApplicationUseCase(repositories *repository.Repository, config *config.Config, ports *ports.Ports) *ApplicationUseCase {
	return &ApplicationUseCase{repositories: repositories, config: config, ports: ports}
}

func (s *ApplicationUseCase) SetupUseCases() {
	s.InventoryUseCase = inventory_usecase.NewInventoryUseCase(s.repositories, s.repositories.InventoryRepository, s.repositories.InventoryItemRepository, s.repositories.InventoryTransactionRepository, s.repositories.SkuRepository)
	s.EmailUseCase = emailusecase.NewEmailUseCase(s.config, s.ports.EmailPort)
	s.SalesUsecase = sales_usecase.NewSalesUseCase(
		s.repositories.UserRepository,
		s.repositories.CustomerRepository,
		s.repositories.SkuRepository,
		s.repositories.SalesRepository,
		s.InventoryUseCase,
		s.repositories.InventoryRepository,
		s.repositories.InventoryItemRepository,
		s.repositories,
	)
}
