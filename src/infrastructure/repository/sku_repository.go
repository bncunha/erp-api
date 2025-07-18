package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bncunha/erp-api/src/domain"
)

type SkuRepository interface {
	CreateMany(ctx context.Context, skus []domain.Sku, productId int64) ([]int64, error)
	GetByProductId(ctx context.Context, productId int64) ([]domain.Sku, error)
}

type skuRepository struct {
	db *sql.DB
}

func NewSkuRepository(db *sql.DB) SkuRepository {
	return &skuRepository{db}
}

func (r *skuRepository) CreateMany(ctx context.Context, skus []domain.Sku, productId int64) ([]int64, error) {
	var insertedIDs []int64
	var values []any
	var placeholders []string

	for i, sku := range skus {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
		values = append(values, sku.Code, sku.Color, sku.Size, sku.Cost, sku.Price, productId)
	}

	query := fmt.Sprintf(`
		INSERT INTO skus (code, color, size, cost, price, product_id)
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
	var skus []domain.Sku = make([]domain.Sku, 0)

	query := `SELECT id, code, color, size, cost, price FROM skus WHERE product_id = $1`
	rows, err := r.db.QueryContext(ctx, query, productId)
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