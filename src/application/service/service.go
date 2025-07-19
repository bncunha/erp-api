package service

import "github.com/bncunha/erp-api/src/infrastructure/repository"

type ApplicationService struct {
	ProductService ProductService
	SkuService SkuService
	CategoryService CategoryService
	repositories *repository.Repository
}

func NewApplicationService(repositories *repository.Repository) *ApplicationService {
	return &ApplicationService{repositories: repositories}
}

func (s *ApplicationService) SetupServices() {
	s.ProductService = NewProductService(s.repositories.ProductRepository, s.repositories.CategoryRepository, s.repositories.SkuRepository)
	s.SkuService = NewSkuService(s.repositories.SkuRepository)
	s.CategoryService = NewCategoryService(s.repositories.CategoryRepository)
}