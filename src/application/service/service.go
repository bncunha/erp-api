package service

import (
	"github.com/bncunha/erp-api/src/application/usecase"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type ApplicationService struct {
	ProductService   ProductService
	SkuService       SkuService
	CategoryService  CategoryService
	AuthService      AuthService
	UserService      UserService
	InventoryService InventoryService
	SalesService     SalesService
	repositories     *repository.Repository
	useCases         *usecase.ApplicationUseCase
}

func NewApplicationService(repositories *repository.Repository, useCases *usecase.ApplicationUseCase) *ApplicationService {
	return &ApplicationService{repositories: repositories, useCases: useCases}
}

func (s *ApplicationService) SetupServices() {
	s.ProductService = NewProductService(s.repositories.ProductRepository, s.repositories.CategoryRepository, s.repositories.SkuRepository)
	s.SkuService = NewSkuService(s.repositories.SkuRepository, s.useCases.InventoryUseCase, s.repositories.ProductRepository)
	s.CategoryService = NewCategoryService(s.repositories.CategoryRepository)
	s.AuthService = NewAuthService(s.repositories.UserRepository)
	s.UserService = NewUserService(s.repositories.UserRepository, s.repositories.InventoryRepository)
	s.InventoryService = NewInventoryService(s.useCases.InventoryUseCase, s.repositories.InventoryItemRepository, s.repositories.InventoryTransactionRepository, s.repositories.InventoryRepository)
	s.SalesService = NewSalesService(s.useCases.SalesUsecase)
}
