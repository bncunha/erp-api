package service

import (
	"context"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type CategoryService interface {
	Create(ctx context.Context, input request.CreateCategoryRequest) error
	Edit(ctx context.Context, input request.EditCategoryRequest) error
	GetById(ctx context.Context, id int64) (domain.Category, error)
	GetAll(ctx context.Context) ([]domain.Category, error)
	Inactivate(ctx context.Context, id int64) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository}
}

func (s *categoryService) Create(ctx context.Context, input request.CreateCategoryRequest) error {
	err := input.Validate()
	if err != nil {
		return err
	}

	_, err = s.categoryRepository.Create(ctx, domain.Category{
		Name: input.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *categoryService) Edit(ctx context.Context, input request.EditCategoryRequest) error {
	err := input.Validate()
	if err != nil {
		return err
	}

	err = s.categoryRepository.Update(ctx, domain.Category{
		Id: input.Id,
		Name: input.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *categoryService) GetById(ctx context.Context, id int64) (domain.Category, error) {
	category, err := s.categoryRepository.GetById(ctx, id)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (s *categoryService) GetAll(ctx context.Context) ([]domain.Category, error) {
	categories, err := s.categoryRepository.GetAll(ctx)
	if err != nil {
		return categories, err
	}
	return categories, nil
}

func (s *categoryService) Inactivate(ctx context.Context, id int64) error {
	return s.categoryRepository.Delete(ctx, id)
}