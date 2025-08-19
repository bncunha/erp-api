package service

import "github.com/bncunha/erp-api/src/infrastructure/repository"

type ApplicationService struct {
	ProductService  ProductService
	SkuService      SkuService
	CategoryService CategoryService
	AuthService     AuthService
	UserService     UserService
	repositories    *repository.Repository
}

func NewApplicationService(repositories *repository.Repository) *ApplicationService {
	return &ApplicationService{repositories: repositories}
}

func (s *ApplicationService) SetupServices() {
	s.ProductService = NewProductService(s.repositories.ProductRepository, s.repositories.CategoryRepository, s.repositories.SkuRepository)
	s.SkuService = NewSkuService(s.repositories.SkuRepository)
	s.CategoryService = NewCategoryService(s.repositories.CategoryRepository)
	s.AuthService = NewAuthService(s.repositories.UserRepository)
	s.UserService = NewUserService(s.repositories.UserRepository, s.repositories.InventoryRepository)
}
