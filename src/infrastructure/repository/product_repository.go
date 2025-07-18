package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product domain.Product) (int64, error)
	Edit(ctx context.Context, product domain.Product, id int64) (int64, error)
	GetById(ctx context.Context, id int64) (domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	Inactivate(ctx context.Context, id int64) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(ctx context.Context, product domain.Product) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedId int64
	var query string
	var args []any
	if product.Category.Id == 0 {
		query = `INSERT INTO products (name, description, tenant_id) VALUES ($1, $2, $3) RETURNING id`
		args = []any{product.Name, product.Description, tenantId}
	} else {
		query = `INSERT INTO products (name, description, tenant_id, category_id) VALUES ($1, $2, $3, $4) RETURNING id`
		args = []any{product.Name, product.Description, tenantId, product.Category.Id}
	}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&insertedId)
	return insertedId, err
}

func (r *productRepository) Edit(ctx context.Context, product domain.Product, id int64) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var updatedId int64
	var query string
	var args []any

	if product.Category.Id == 0 {
		query = `UPDATE products SET name = $1, description = $2 WHERE id = $4 AND tenant_id = $5 AND deleted_at IS NULL RETURNING id`
		args = []any{product.Name, product.Description, id, tenantId}
	} else {
		query = `UPDATE products SET name = $1, description = $2, category_id = $4 WHERE id = $5 AND tenant_id = $6 AND deleted_at IS NULL RETURNING id`
		args = []any{product.Name, product.Description, product.Category.Id, id, tenantId}
	}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&updatedId)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return updatedId, errors.New("Produto não encontrado")
		}
		return updatedId, err
	}
	return updatedId, err
}

func (r *productRepository) GetById(ctx context.Context, id int64) (domain.Product, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var product domain.Product
	var categoryID sql.NullInt64

	query := `SELECT id, name, description, category_id FROM products WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&product.Id, &product.Name, &product.Description, &categoryID)
	if err != nil {
		return product, err
	}

	if categoryID.Valid {
		product.Category.Id = categoryID.Int64
	}
	return product, nil
}

func (r *productRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var products []domain.Product

	query := `SELECT id, name, description FROM products WHERE tenant_id = $1 AND deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query, tenantId)
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var product domain.Product
		err = rows.Scan(&product.Id, &product.Name, &product.Description)
		if err != nil {
			return products, err
		}
		products = append(products, product)
	}
	return products, err
}

func (r *productRepository) Inactivate(ctx context.Context, id int64) error {
	query := `UPDATE products SET deleted_at = now() WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, ctx.Value(constants.TENANT_KEY))
	return err
}