package service

import (
	product_service "github.com/bncunha/erp-api/src/application/service/product"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type ApplicationService struct {
	ProductService product_service.ProductService
	repositories *repository.Repository
}

func NewApplicationService(repositories *repository.Repository) *ApplicationService {
	return &ApplicationService{repositories: repositories}
}

func (s *ApplicationService) SetupServices() {
	s.ProductService = product_service.NewProductService(s.repositories.ProductRepository)
}