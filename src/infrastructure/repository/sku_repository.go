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