package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/domain"
	config "github.com/bncunha/erp-api/src/main"
)

type ProductRepository interface {
	Create(ctx context.Context, product domain.Product) (int64, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(ctx context.Context, product domain.Product) (int64, error) {
	var insertedId int64
	var query string
	var args []any
	if product.Category.Id == 0 {
		query = `INSERT INTO products (name, description, company_id) VALUES ($1, $2, $3) RETURNING id`
		args = []any{product.Name, product.Description, config.COMPANY_TEST_ID}
	} else {
		query = `INSERT INTO products (name, description, company_id, category_id) VALUES ($1, $2, $3, $4) RETURNING id`
		args = []any{product.Name, product.Description, config.COMPANY_TEST_ID, product.Category.Id}
	}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&insertedId)
	return insertedId, err
}