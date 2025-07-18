package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
)

type SkuRepository interface {
	Create(ctx context.Context, sku domain.Sku, productId int64) (int64, error)
	CreateMany(ctx context.Context, skus []domain.Sku, productId int64) ([]int64, error)
	GetByProductId(ctx context.Context, productId int64) ([]domain.Sku, error)
}

type skuRepository struct {
	db *sql.DB
}

func NewSkuRepository(db *sql.DB) SkuRepository {
	return &skuRepository{db}
}

func (r *skuRepository) Create(ctx context.Context, sku domain.Sku, productId int64) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedID int64

	query := `INSERT INTO skus (code, color, size, cost, price, product_id, tenant_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, sku.Code, sku.Color, sku.Size, sku.Cost, sku.Price, productId, tenantId).Scan(&insertedID)
	return insertedID, err
}

func (r *skuRepository) CreateMany(ctx context.Context, skus []domain.Sku, productId int64) ([]int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedIDs []int64
	var values []any
	var placeholders []string

	for i, sku := range skus {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7))
		values = append(values, sku.Code, sku.Color, sku.Size, sku.Cost, sku.Price, productId, tenantId)
	}

	query := fmt.Sprintf(`
		INSERT INTO skus (code, color, size, cost, price, product_id, tenant_id)
		VALUES %s
		RETURNING id
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return insertedIDs, err
	}
	defer rows.Close()

	for rows.Next() {
		var insertedID int64
		err = rows.Scan(&insertedID)
		if err != nil {
			return insertedIDs, err
		}
		insertedIDs = append(insertedIDs, insertedID)
	}
	return insertedIDs, err
}

func (r *skuRepository) GetByProductId(ctx context.Context, productId int64) ([]domain.Sku, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var skus []domain.Sku = make([]domain.Sku, 0)

	query := `SELECT id, code, color, size, cost, price FROM skus WHERE product_id = $1 AND tenant_id = $2`
	rows, err := r.db.QueryContext(ctx, query, productId, tenantId)
	if err != nil {
		return skus, err
	}
	defer rows.Close()

	for rows.Next() {
		var sku domain.Sku
		err = rows.Scan(&sku.Id, &sku.Code, &sku.Color, &sku.Size, &sku.Cost, &sku.Price)
		if err != nil {
			return skus, err
		}
		skus = append(skus, sku)
	}
	return skus, err
}