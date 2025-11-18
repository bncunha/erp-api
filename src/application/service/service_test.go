package service

import (
	"testing"

	"github.com/bncunha/erp-api/src/application/ports"
	"github.com/bncunha/erp-api/src/application/usecase"
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

func TestNewApplicationService(t *testing.T) {
	repos := &repository.Repository{}
	useCases := &usecase.ApplicationUseCase{}
	service := NewApplicationService(repos, useCases, ports.NewPorts(&stubEncrypto{}, &stubEmailPort{}))
	if service == nil {
		t.Fatalf("expected service to be created")
	}
}

func TestApplicationServiceSetup(t *testing.T) {
	repos := &repository.Repository{
		ProductRepository:              &stubProductRepository{},
		CategoryRepository:             &stubCategoryRepository{},
		SkuRepository:                  &stubSkuRepository{},
		UserRepository:                 &stubUserRepository{},
		UserTokenRepository:            &stubUserTokenRepository{},
		InventoryRepository:            &stubInventoryRepository{},
		InventoryItemRepository:        &stubInventoryItemRepository{},
		InventoryTransactionRepository: &stubInventoryTransactionRepository{},
	}
	useCases := &usecase.ApplicationUseCase{
		InventoryUseCase: inventory_usecase.NewInventoryUseCase(nil, repos.InventoryRepository, repos.InventoryItemRepository, repos.InventoryTransactionRepository, repos.SkuRepository),
	}

	service := NewApplicationService(repos, useCases, ports.NewPorts(&stubEncrypto{}, &stubEmailPort{}))
	service.SetupServices()

	if service.ProductService == nil || service.InventoryService == nil {
		t.Fatalf("expected services to be initialized")
	}
}
