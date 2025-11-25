package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/lib/pq"
)

type skuRepository struct {
	db *sql.DB
}

func NewSkuRepository(db *sql.DB) domain.SkuRepository {
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

	query := `
                SELECT s.id, s.code, s.color, s.size, s.cost, s.price, p.name, sum(inv_item.quantity)
                FROM skus s
                INNER JOIN products p ON p.id = s.product_id
                LEFT JOIN inventory_items inv_item ON inv_item.sku_id = s.id
                WHERE s.product_id = $1 AND s.tenant_id = $2 AND s.deleted_at IS NULL
                GROUP BY s.id, p.name
                ORDER BY s.id ASC`

	rows, err := r.db.QueryContext(ctx, query, productId, tenantId)
	if err != nil {
		return skus, err
	}
	defer rows.Close()

	for rows.Next() {
		var sku domain.Sku
		var quantity sql.NullFloat64
		err = rows.Scan(&sku.Id, &sku.Code, &sku.Color, &sku.Size, &sku.Cost, &sku.Price, &sku.Product.Name, &quantity)
		if err != nil {
			return skus, err
		}
		if quantity.Valid {
			sku.Quantity = quantity.Float64
		}
		skus = append(skus, sku)
	}
	return skus, err
}

func (r *skuRepository) Update(ctx context.Context, sku domain.Sku) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `UPDATE skus SET code = $1, color = $2, size = $3, cost = $4, price = $5 WHERE id = $6 AND tenant_id = $7 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, sku.Code, sku.Color, sku.Size, sku.Cost, sku.Price, sku.Id, tenantId)
	return err
}

func (r *skuRepository) GetById(ctx context.Context, id int64) (domain.Sku, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var sku domain.Sku

	query := `SELECT s.id, s.code, s.color, s.size, s.cost, s.price, p.name
	FROM skus s 
	INNER JOIN products p ON p.id = s.product_id
	WHERE s.id = $1 AND s.tenant_id = $2 AND s.deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&sku.Id, &sku.Code, &sku.Color, &sku.Size, &sku.Cost, &sku.Price, &sku.Product.Name)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return sku, errors.New("SKU não encontrada")
		}
		return sku, err
	}
	return sku, nil
}

func (r *skuRepository) GetByManyIds(ctx context.Context, ids []int64) ([]domain.Sku, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var skus []domain.Sku

	query := `SELECT s.id, s.code, s.color, s.size, s.cost, s.price, p.name
	FROM skus s 
	INNER JOIN products p ON p.id = s.product_id
	WHERE s.id = ANY($1) AND s.tenant_id = $2 AND s.deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query, pq.Array(ids), tenantId)
	if err != nil {
		return skus, err
	}
	defer rows.Close()

	for rows.Next() {
		var sku domain.Sku
		err = rows.Scan(&sku.Id, &sku.Code, &sku.Color, &sku.Size, &sku.Cost, &sku.Price, &sku.Product.Name)
		if err != nil {
			return skus, err
		}
		skus = append(skus, sku)
	}
	return skus, err
}

func (r *skuRepository) GetAll(ctx context.Context, input domain.GetSkusInput) ([]domain.Sku, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var skus []domain.Sku

	query := `SELECT s.id, s.code, s.color, s.size, s.cost, s.price, p.name, sum(inv_item.quantity)
	FROM skus s 
	INNER JOIN products p ON p.id = s.product_id
	LEFT JOIN inventory_items inv_item ON inv_item.sku_id = s.id
	LEFT JOIN inventories inv ON inv.id = inv_item.inventory_id
	WHERE s.tenant_id = $1 AND s.deleted_at IS NULL
	AND ($2::bigint IS NULL OR inv.user_id = $2::bigint)
	GROUP BY s.id, p.name
	ORDER BY s.id ASC`
	rows, err := r.db.QueryContext(ctx, query, tenantId, input.SellerId)
	if err != nil {
		return skus, err
	}
	defer rows.Close()

	for rows.Next() {
		var sku domain.Sku
		var quantity sql.NullFloat64
		err = rows.Scan(&sku.Id, &sku.Code, &sku.Color, &sku.Size, &sku.Cost, &sku.Price, &sku.Product.Name, &quantity)
		if err != nil {
			return skus, err
		}
		if quantity.Valid {
			sku.Quantity = quantity.Float64
		}
		skus = append(skus, sku)
	}
	return skus, err
}

func (r *skuRepository) Inactivate(ctx context.Context, id int64) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `DELETE FROM skus WHERE id = $1 AND tenant_id = $2`
	result, err := r.db.ExecContext(ctx, query, id, tenantId)
	if err != nil {
		if errors.IsForeignKeyViolation(err) {
			return errors.New("Não é possível deletar o SKU pois existem registros associados.")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("SKU não encontrada")
	}

	return nil
}
