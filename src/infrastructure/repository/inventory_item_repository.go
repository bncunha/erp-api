package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

var (
	ErrInventoryItemNotFound = errors.New("Item de estoque não encontrado")
)

type InventoryItemRepository interface {
	Create(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) (int64, error)
	UpdateQuantity(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) error
	GetById(ctx context.Context, id int64) (domain.InventoryItem, error)
	GetByIdWithTransaction(ctx context.Context, tx *sql.Tx, id int64) (domain.InventoryItem, error)
	GetBySkuIdAndInventoryId(ctx context.Context, skuId int64, inventoryId int64) (domain.InventoryItem, error)
	GetAll(ctx context.Context) ([]output.GetInventoryItemsOutput, error)
}

type inventoryItemRepository struct {
	db *sql.DB
}

func NewInventoryItemRepository(db *sql.DB) InventoryItemRepository {
	return &inventoryItemRepository{db}
}

func (r *inventoryItemRepository) Create(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedID int64

	query := `INSERT INTO inventory_items (inventory_id, sku_id, quantity, tenant_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := tx.QueryRowContext(ctx, query, inventoryItem.InventoryId, inventoryItem.SkuId, inventoryItem.Quantity, tenantId).Scan(&insertedID)
	return insertedID, err
}

func (r *inventoryItemRepository) UpdateQuantity(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `UPDATE inventory_items SET quantity = $1 WHERE id = $2 AND tenant_id = $3 AND deleted_at IS NULL`
	_, err := tx.ExecContext(ctx, query, inventoryItem.Quantity, inventoryItem.Id, tenantId)
	return err
}

func (r *inventoryItemRepository) GetById(ctx context.Context, id int64) (domain.InventoryItem, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventoryItem domain.InventoryItem

	query := `SELECT id, inventory_id, sku_id, quantity, tenant_id FROM inventory_items WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&inventoryItem.Id, &inventoryItem.InventoryId, &inventoryItem.SkuId, &inventoryItem.Quantity, &tenantId)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return inventoryItem, ErrInventoryItemNotFound
		}
		return inventoryItem, err
	}
	return inventoryItem, nil
}

func (r *inventoryItemRepository) GetBySkuIdAndInventoryId(ctx context.Context, skuId int64, inventoryId int64) (domain.InventoryItem, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventoryItem domain.InventoryItem

	query := `SELECT id, inventory_id, sku_id, quantity, tenant_id FROM inventory_items WHERE sku_id = $1 AND inventory_id = $2 AND tenant_id = $3 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, skuId, inventoryId, tenantId).Scan(&inventoryItem.Id, &inventoryItem.InventoryId, &inventoryItem.SkuId, &inventoryItem.Quantity, &tenantId)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return inventoryItem, ErrInventoryItemNotFound
		}
		return inventoryItem, err
	}
	return inventoryItem, nil
}

func (r *inventoryItemRepository) GetByIdWithTransaction(ctx context.Context, tx *sql.Tx, id int64) (domain.InventoryItem, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventoryItem domain.InventoryItem

	query := `SELECT id, inventory_id, sku_id, quantity, tenant_id FROM inventory_items WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	err := tx.QueryRowContext(ctx, query, id, tenantId).Scan(&inventoryItem.Id, &inventoryItem.InventoryId, &inventoryItem.SkuId, &inventoryItem.Quantity, &tenantId)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return inventoryItem, ErrInventoryItemNotFound
		}
		return inventoryItem, err
	}
	return inventoryItem, nil
}

func (r *inventoryItemRepository) GetAll(ctx context.Context) ([]output.GetInventoryItemsOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventoryItems []output.GetInventoryItemsOutput

	query := ` SELECT inv_items.id, sku.code, sku.color, sku.size, p.name, inv.type, u.name, inv_items.quantity
	FROM inventory_items inv_items 
	INNER JOIN inventories inv ON inv.id = inv_items.inventory_id
	INNER JOIN skus sku ON sku.id = inv_items.sku_id
	INNER JOIN products p ON p.id = sku.product_id
	LEFT JOIN users u ON u.id = inv.user_id 
	WHERE inv_items.tenant_id = $1 AND inv_items.deleted_at IS NULL ORDER BY inv_items.id ASC`
	rows, err := r.db.QueryContext(ctx, query, tenantId)
	if err != nil {
		return inventoryItems, err
	}
	defer rows.Close()

	for rows.Next() {
		var inventoryItem output.GetInventoryItemsOutput

		err = rows.Scan(&inventoryItem.InventoryItemId, &inventoryItem.SkuCode, &inventoryItem.SkuColor, &inventoryItem.SkuSize, &inventoryItem.ProductName, &inventoryItem.InventoryType, &inventoryItem.UserName, &inventoryItem.Quantity)
		if err != nil {
			return inventoryItems, err
		}
		inventoryItems = append(inventoryItems, inventoryItem)
	}
	return inventoryItems, err
}
