package product_service

import (
	"context"

	"github.com/bncunha/erp-api/src/application/service/product/input"
	"github.com/bncunha/erp-api/src/application/validator"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type ProductService interface {
	Create(ctx context.Context, input input.CreateProductInput) error
}

type productService struct{
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository}
}

func (s *productService) Create(ctx context.Context, input input.CreateProductInput) error {
	err := validator.Validate(input)
	if err != nil {
		return err
	}

	err = s.productRepository.Create(ctx, domain.Product{
		Name: input.Name,
		Description: input.Description,
	})
	if err != nil {
		return err
	}	
	return nil
}
