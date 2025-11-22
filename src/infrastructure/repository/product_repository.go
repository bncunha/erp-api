package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
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
		query = `UPDATE products SET name = $1, description = $2 WHERE id = $3 AND tenant_id = $4 AND deleted_at IS NULL RETURNING id`
		args = []any{product.Name, product.Description, id, tenantId}
	} else {
		query = `UPDATE products SET name = $1, description = $2, category_id = $3 WHERE id = $4 AND tenant_id = $5 AND deleted_at IS NULL RETURNING id`
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

func (r *productRepository) GetAll(ctx context.Context, input domain.GetProductsInput) ([]domain.GetAllProductsOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var products []domain.GetAllProductsOutput

	query := `
		SELECT p.id, p.name, p.description, c.name AS category_name, c.id AS category_id, sum(inv_item.quantity) 
		FROM products p 
		LEFT JOIN skus sku ON sku.product_id = p.id 
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN inventory_items inv_item ON sku.id = inv_item.sku_id
		LEFT JOIN inventories inv ON inv.id = inv_item.inventory_id
		WHERE p.tenant_id = $1 AND p.deleted_at IS NULL

		AND ($2::bigint IS NULL OR inv.user_id = $2::bigint)

		GROUP BY p.id, c.id
		ORDER BY p.id ASC`

	rows, err := r.db.QueryContext(ctx, query, tenantId, input.SellerId)
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var product domain.Product
		var categoryName sql.NullString
		var categoryId sql.NullInt64
		var quantity sql.NullFloat64

		err = rows.Scan(&product.Id, &product.Name, &product.Description, &categoryName, &categoryId, &quantity)
		if err != nil {
			return products, err
		}
		if categoryId.Valid {
			product.Category.Id = categoryId.Int64
		}
		if categoryName.Valid {
			product.Category.Name = categoryName.String
		}
		products = append(products, domain.GetAllProductsOutput{Product: product, Quantity: quantity.Float64})
	}
	return products, err
}

func (r *productRepository) Inactivate(ctx context.Context, id int64) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `DELETE FROM products WHERE id = $1 AND tenant_id = $2`
	result, err := r.db.ExecContext(ctx, query, id, tenantId)
	if err != nil {
		if errors.IsForeignKeyViolation(err) {
			return errors.New("Não é possível deletar o produto pois existem registros associados.")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Produto não encontrado")
	}

	return nil
}
