package service

import (
	"context"
	"errors"
	"testing"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/domain"
)

func TestCategoryServiceCreate(t *testing.T) {
	repo := &stubCategoryRepository{}
	service := &categoryService{categoryRepository: repo}
	req := request.CreateCategoryRequest{Name: "Cat"}

	if err := service.Create(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.created.Name != "Cat" {
		t.Fatalf("expected category creation")
	}
}

func TestCategoryServiceEdit(t *testing.T) {
	repo := &stubCategoryRepository{}
	service := &categoryService{categoryRepository: repo}
	req := request.EditCategoryRequest{Id: 1, CreateCategoryRequest: request.CreateCategoryRequest{Name: "Cat"}}

	if err := service.Edit(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.created.Name != "Cat" {
		t.Fatalf("expected category update")
	}
}

func TestCategoryServiceGetById(t *testing.T) {
	repo := &stubCategoryRepository{getById: domain.Category{Id: 1}}
	service := &categoryService{categoryRepository: repo}

	category, err := service.GetById(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if category.Id != 1 {
		t.Fatalf("expected category")
	}
}

func TestCategoryServiceGetAll(t *testing.T) {
	repo := &stubCategoryRepository{getAll: []domain.Category{{Id: 1}}}
	service := &categoryService{categoryRepository: repo}

	categories, err := service.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(categories) != 1 {
		t.Fatalf("expected categories")
	}
}

func TestCategoryServiceInactivate(t *testing.T) {
	repo := &stubCategoryRepository{}
	service := &categoryService{categoryRepository: repo}

	if err := service.Inactivate(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCategoryServiceCreateValidationError(t *testing.T) {
	service := &categoryService{}
	if err := service.Create(context.Background(), request.CreateCategoryRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestCategoryServiceCreateRepositoryError(t *testing.T) {
	repo := &stubCategoryRepository{createErr: errors.New("fail")}
	service := &categoryService{categoryRepository: repo}
	if err := service.Create(context.Background(), request.CreateCategoryRequest{Name: "Cat"}); err == nil || err.Error() != "fail" {
		t.Fatalf("expected repository error")
	}
}

func TestCategoryServiceEditValidationError(t *testing.T) {
	service := &categoryService{}
	if err := service.Edit(context.Background(), request.EditCategoryRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestCategoryServiceEditRepositoryError(t *testing.T) {
	repo := &stubCategoryRepository{updateErr: errors.New("fail")}
	service := &categoryService{categoryRepository: repo}
	if err := service.Edit(context.Background(), request.EditCategoryRequest{Id: 1, CreateCategoryRequest: request.CreateCategoryRequest{Name: "Cat"}}); err == nil || err.Error() != "fail" {
		t.Fatalf("expected repository error")
	}
}

func TestCategoryServiceGetByIdError(t *testing.T) {
	repo := &stubCategoryRepository{getByIdErr: errors.New("fail")}
	service := &categoryService{categoryRepository: repo}
	if _, err := service.GetById(context.Background(), 1); err == nil {
		t.Fatalf("expected error")
	}
}

func TestCategoryServiceGetAllError(t *testing.T) {
	repo := &stubCategoryRepository{getAllErr: errors.New("fail")}
	service := &categoryService{categoryRepository: repo}
	if _, err := service.GetAll(context.Background()); err == nil {
		t.Fatalf("expected error")
	}
}
