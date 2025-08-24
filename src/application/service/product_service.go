package service

import (
	"context"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/application/validator"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type ProductService interface {
	Create(ctx context.Context, input request.CreateProductRequest) (int64, error)
	Edit(ctx context.Context, input request.EditProductRequest) error
	GetById(ctx context.Context, id int64) (domain.Product, error)
	GetAll(ctx context.Context) ([]output.GetAllProductsOutput, error)
	Inactivate(ctx context.Context, id int64) error
	GetSkus(ctx context.Context, id int64) ([]domain.Sku, error)
}

type productService struct {
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
	skuRepository      repository.SkuRepository
}

func NewProductService(productRepository repository.ProductRepository, categoryRepository repository.CategoryRepository, skuRepositoy repository.SkuRepository) ProductService {
	return &productService{productRepository, categoryRepository, skuRepositoy}
}

func (s *productService) Create(ctx context.Context, input request.CreateProductRequest) (int64, error) {
	err := input.Validate()
	if err != nil {
		return 0, err
	}

	category, err := s.getCategory(ctx, input.CategoryID, input.CategoryName)
	if err != nil {
		return 0, err
	}

	productId, err := s.productRepository.Create(ctx, domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Category:    category,
	})
	if err != nil {
		return 0, err
	}

	_, err = s.insertSkus(ctx, input.Skus, productId)
	if err != nil {
		return 0, err
	}

	return productId, nil
}

func (s *productService) Edit(ctx context.Context, input request.EditProductRequest) error {
	err := validator.Validate(input)
	if err != nil {
		return err
	}

	category, err := s.getCategory(ctx, input.CategoryID, input.CategoryName)
	if err != nil {
		return err
	}

	_, err = s.productRepository.Edit(ctx, domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Category:    category,
	}, input.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) GetById(ctx context.Context, id int64) (domain.Product, error) {
	product, err := s.productRepository.GetById(ctx, id)
	if err != nil {
		return product, err
	}

	product.Skus, err = s.skuRepository.GetByProductId(ctx, id)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *productService) GetAll(ctx context.Context) ([]output.GetAllProductsOutput, error) {
	products, err := s.productRepository.GetAll(ctx)
	if err != nil {
		return products, err
	}
	return products, nil
}

func (s *productService) Inactivate(ctx context.Context, id int64) error {
	return s.productRepository.Inactivate(ctx, id)
}

func (s *productService) GetSkus(ctx context.Context, id int64) ([]domain.Sku, error) {
	skus, err := s.skuRepository.GetByProductId(ctx, id)
	if err != nil {
		return skus, err
	}
	return skus, nil
}

func (s *productService) insertSkus(ctx context.Context, skus []request.CreateSkuRequest, productId int64) ([]domain.Sku, error) {
	var skusDomain []domain.Sku
	for _, sku := range skus {
		skusDomain = append(skusDomain, domain.Sku{
			Code:  sku.Code,
			Color: sku.Color,
			Size:  sku.Size,
			Cost:  sku.Cost,
			Price: sku.Price,
		})
	}

	if len(skusDomain) > 0 {
		_, err := s.skuRepository.CreateMany(ctx, skusDomain, productId)
		if err != nil {
			return skusDomain, err
		}
	}
	return skusDomain, nil
}

func (s *productService) getCategory(ctx context.Context, categoryId int64, categoryName string) (domain.Category, error) {
	if categoryId == 0 && categoryName == "" {
		return domain.Category{}, nil
	}

	if categoryId != 0 {
		return s.categoryRepository.GetById(ctx, categoryId)
	}

	categoryId, err := s.categoryRepository.Create(ctx, domain.Category{
		Name: categoryName,
	})
	if err != nil && err.Error() != "Categoria já existe no banco de dados" {
		return domain.Category{}, err
	} else if err != nil {
		return s.categoryRepository.GetByName(ctx, categoryName)
	}
	return s.categoryRepository.GetById(ctx, categoryId)
}
