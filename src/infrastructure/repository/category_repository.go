package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, category domain.Category) (int64, error)
	GetById(ctx context.Context, id int64) (domain.Category, error)
	GetByName(ctx context.Context, name string) (domain.Category, error)
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) Create(ctx context.Context, category domain.Category) (int64, error) {
	var insertedID int64

	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, category.Name).Scan(&insertedID)
	if err != nil {
		if errors.IsUniqueViolation(err) {
			return insertedID, errors.New("Categoria já existe no banco de dados")
		}
		return insertedID, err
	}
	return insertedID, nil
}

func (r *categoryRepository) GetById(ctx context.Context, id int64) (domain.Category, error) {
	var category domain.Category

	query := `SELECT id, name FROM categories WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&category.Id, &category.Name)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return category, errors.New("Categoria não encontrada")
		}
		return category, err
	}
	return category, nil
}

func (r *categoryRepository) GetByName(ctx context.Context, name string) (domain.Category, error) {
	var category domain.Category

	query := `SELECT id, name FROM categories WHERE name = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&category.Id, &category.Name)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return category, errors.New("Categoria não encontrada")
		}
		return category, err
	}
	return category, nil
}