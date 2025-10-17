package service

import (
	"context"
	"errors"
	"testing"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

func TestProductServiceCreate(t *testing.T) {
	productRepo := &stubProductRepository{}
	categoryRepo := &stubCategoryRepository{getById: domain.Category{Id: 1, Name: "Cat"}}
	skuRepo := &stubSkuRepository{}

	service := &productService{productRepository: productRepo, categoryRepository: categoryRepo, skuRepository: skuRepo}
	cost := 10.0
	price := 15.0
	req := request.CreateProductRequest{
		Name:         "Product",
		Description:  "Desc",
		CategoryName: "Cat",
		Skus: []request.CreateSkuRequest{{
			Code:  "code",
			Color: "red",
			Size:  "M",
			Cost:  &cost,
			Price: price,
		}},
	}

	if _, err := service.Create(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if productRepo.created.Name != "Product" {
		t.Fatalf("expected product to be created")
	}
	if len(skuRepo.created) == 0 {
		t.Fatalf("expected skus to be inserted")
	}
}

func TestProductServiceEdit(t *testing.T) {
	productRepo := &stubProductRepository{}
	categoryRepo := &stubCategoryRepository{getById: domain.Category{Id: 1, Name: "Cat"}}
	service := &productService{productRepository: productRepo, categoryRepository: categoryRepo}
	req := request.EditProductRequest{Id: 1, Name: "New", CategoryID: 1}

	if err := service.Edit(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if productRepo.created.Name != "New" {
		t.Fatalf("expected product to be updated")
	}
}

func TestProductServiceGetById(t *testing.T) {
	productRepo := &stubProductRepository{getById: domain.Product{Id: 1}}
	skuRepo := &stubSkuRepository{getByProduct: []domain.Sku{{Id: 2}}}
	service := &productService{productRepository: productRepo, skuRepository: skuRepo}

	product, err := service.GetById(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(product.Skus) != 1 {
		t.Fatalf("expected sku to be loaded")
	}
}

func TestProductServiceGetAll(t *testing.T) {
	productRepo := &stubProductRepository{getAll: []output.GetAllProductsOutput{{}}}
	service := &productService{productRepository: productRepo}

	ctx := context.WithValue(context.Background(), constants.ROLE_KEY, string(domain.UserRoleAdmin))
	products, err := service.GetAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 1 {
		t.Fatalf("expected products")
	}
}

func TestProductServiceInactivate(t *testing.T) {
	productRepo := &stubProductRepository{}
	service := &productService{productRepository: productRepo}

	if err := service.Inactivate(context.Background(), 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if productRepo.created.Id != 3 {
		t.Fatalf("expected product to be marked inactive")
	}
}

func TestProductServiceGetSkus(t *testing.T) {
	skuRepo := &stubSkuRepository{getByProduct: []domain.Sku{{Id: 1}}}
	service := &productService{skuRepository: skuRepo}

	skus, err := service.GetSkus(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skus) != 1 {
		t.Fatalf("expected skus")
	}
}

func TestProductServiceGetCategoryByName(t *testing.T) {
	categoryRepo := &stubCategoryRepository{
		createErr:    nil,
		getById:      domain.Category{Id: 1, Name: "Cat"},
		getByName:    domain.Category{Id: 1, Name: "Cat"},
		getByNameErr: nil,
	}
	service := &productService{categoryRepository: categoryRepo}

	category, err := service.getCategory(context.Background(), 0, "Cat")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if category.Name != "Cat" {
		t.Fatalf("expected category by name")
	}
}

func TestProductServiceGetCategoryExistingId(t *testing.T) {
	categoryRepo := &stubCategoryRepository{getById: domain.Category{Id: 2}}
	service := &productService{categoryRepository: categoryRepo}

	category, err := service.getCategory(context.Background(), 2, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if category.Id != 2 {
		t.Fatalf("expected category by id")
	}
}

func TestProductServiceGetCategoryDuplicated(t *testing.T) {
	categoryRepo := &stubCategoryRepository{
		createErr: errors.New("Categoria já existe no banco de dados"),
		getByName: domain.Category{Id: 3, Name: "Cat"},
	}
	service := &productService{categoryRepository: categoryRepo}

	category, err := service.getCategory(context.Background(), 0, "Cat")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if category.Id != 3 {
		t.Fatalf("expected category from name lookup")
	}
}

func TestProductServiceInsertSkus(t *testing.T) {
	skuRepo := &stubSkuRepository{}
	service := &productService{skuRepository: skuRepo}
	cost := 1.0
	price := 2.0
	skus, err := service.insertSkus(context.Background(), []request.CreateSkuRequest{{Code: "c", Color: "c", Size: "s", Cost: &cost, Price: price}}, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skus) != 1 || skus[0].Code != "c" {
		t.Fatalf("expected sku conversion")
	}
}

func TestProductServiceCreateValidationError(t *testing.T) {
	service := &productService{}
	if _, err := service.Create(context.Background(), request.CreateProductRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestProductServiceCreateInsertSkusError(t *testing.T) {
	skuRepo := &stubSkuRepository{createManyErr: errors.New("fail")}
	productRepo := &stubProductRepository{}
	categoryRepo := &stubCategoryRepository{getById: domain.Category{Id: 1}}
	service := &productService{productRepository: productRepo, categoryRepository: categoryRepo, skuRepository: skuRepo}
	cost := 1.0
	price := 2.0
	req := request.CreateProductRequest{Name: "Product", CategoryID: 1, Skus: []request.CreateSkuRequest{{Code: "c", Color: "c", Size: "s", Cost: &cost, Price: price}}}

	if _, err := service.Create(context.Background(), req); err == nil {
		t.Fatalf("expected error from insert skus")
	}
}

func TestProductServiceGetCategoryCreateError(t *testing.T) {
	categoryRepo := &stubCategoryRepository{createErr: errors.New("fail")}
	service := &productService{categoryRepository: categoryRepo}

	if _, err := service.getCategory(context.Background(), 0, "Cat"); err == nil {
		t.Fatalf("expected error from category create")
	}
}

func TestProductServiceEditError(t *testing.T) {
	categoryRepo := &stubCategoryRepository{getById: domain.Category{Id: 1}}
	productRepo := &stubProductRepository{editErr: errors.New("fail")}
	service := &productService{productRepository: productRepo, categoryRepository: categoryRepo}
	req := request.EditProductRequest{Id: 1, Name: "Name", CategoryID: 1}

	if err := service.Edit(context.Background(), req); err == nil {
		t.Fatalf("expected error from update")
	}
}

func TestProductServiceGetByIdError(t *testing.T) {
	productRepo := &stubProductRepository{getByIdErr: errors.New("fail")}
	service := &productService{productRepository: productRepo}

	if _, err := service.GetById(context.Background(), 1); err == nil {
		t.Fatalf("expected error")
	}
}

func TestProductServiceGetSkusError(t *testing.T) {
	skuRepo := &stubSkuRepository{getByProductErr: errors.New("fail")}
	service := &productService{skuRepository: skuRepo}

	if _, err := service.GetSkus(context.Background(), 1); err == nil {
		t.Fatalf("expected error")
	}
}

func TestProductServiceGetAllError(t *testing.T) {
	productRepo := &stubProductRepository{getAllErr: errors.New("fail")}
	service := &productService{productRepository: productRepo}

	ctx := context.WithValue(context.Background(), constants.ROLE_KEY, string(domain.UserRoleAdmin))
	if _, err := service.GetAll(ctx); err == nil {
		t.Fatalf("expected error")
	}
}

func TestProductServiceGetCategoryEmpty(t *testing.T) {
	service := &productService{}
	category, err := service.getCategory(context.Background(), 0, "")
	if err != nil || category != (domain.Category{}) {
		t.Fatalf("expected empty category")
	}
}

func TestProductServiceCreateRepositoryError(t *testing.T) {
	productRepo := &stubProductRepository{createErr: errors.New("fail")}
	service := &productService{productRepository: productRepo, categoryRepository: &stubCategoryRepository{getById: domain.Category{Id: 1}}}
	req := request.CreateProductRequest{Name: "Name", CategoryID: 1}
	if _, err := service.Create(context.Background(), req); err == nil || err.Error() != "fail" {
		t.Fatalf("expected repository error")
	}
}
