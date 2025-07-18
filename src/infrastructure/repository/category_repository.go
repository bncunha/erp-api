package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
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

	query := `INSERT INTO categories (name, tenant_id) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, category.Name, ctx.Value(constants.TENANT_KEY)).Scan(&insertedID)
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

	query := `SELECT id, name FROM categories WHERE id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, id, ctx.Value(constants.TENANT_KEY)).Scan(&category.Id, &category.Name)
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

	query := `SELECT id, name FROM categories WHERE name = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, name, ctx.Value(constants.TENANT_KEY)).Scan(&category.Id, &category.Name)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return category, errors.New("Categoria não encontrada")
		}
		return category, err
	}
	return category, nil
}