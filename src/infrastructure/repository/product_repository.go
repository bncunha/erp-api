package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/domain"
	config "github.com/bncunha/erp-api/src/main"
)

type ProductRepository interface {
	Create(ctx context.Context, product domain.Product) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(ctx context.Context, product domain.Product) error {
	query := `INSERT INTO products (name, description, company_id) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Description, config.COMPANY_TEST_ID)
	return err
}