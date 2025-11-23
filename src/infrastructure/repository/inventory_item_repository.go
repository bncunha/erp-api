package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/lib/pq"
)

type inventoryItemRepository struct {
	db *sql.DB
}

func NewInventoryItemRepository(db *sql.DB) domain.InventoryItemRepository {
	return &inventoryItemRepository{db}
}

func (r *inventoryItemRepository) Create(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedID int64

	query := `INSERT INTO inventory_items (inventory_id, sku_id, quantity, tenant_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := tx.QueryRowContext(ctx, query, inventoryItem.InventoryId, inventoryItem.Sku.Id, inventoryItem.Quantity, tenantId).Scan(&insertedID)
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

	query := `SELECT ii.id, ii.inventory_id, ii.quantity, ii.tenant_id, s.id, s.code, s.color, s.size, s.cost, s.price, p.name
	FROM inventory_items ii
	INNER JOIN skus s ON s.id = ii.sku_id
	INNER JOIN products p ON p.id = s.product_id
	WHERE ii.id = $1 AND ii.tenant_id = $2 AND ii.deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&inventoryItem.Id, &inventoryItem.InventoryId, &inventoryItem.Quantity, &tenantId, &inventoryItem.Sku.Id, &inventoryItem.Sku.Code, &inventoryItem.Sku.Color, &inventoryItem.Sku.Size, &inventoryItem.Sku.Cost, &inventoryItem.Sku.Price, &inventoryItem.Sku.Product.Name)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return inventoryItem, domain.ErrInventoryItemNotFound
		}
		return inventoryItem, err
	}
	return inventoryItem, nil
}

func (r *inventoryItemRepository) GetByManySkuIdsAndInventoryId(ctx context.Context, skuIds []int64, inventoryId int64) ([]domain.InventoryItem, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventoryItems []domain.InventoryItem

	query := `SELECT ii.id, ii.inventory_id, ii.quantity, ii.tenant_id, s.id, s.code, s.color, s.size, s.cost, s.price, p.name
	FROM inventory_items ii
	INNER JOIN skus s ON s.id = ii.sku_id
	INNER JOIN products p ON p.id = s.product_id
	WHERE ii.sku_id = ANY($1) AND ii.inventory_id = $2 AND ii.tenant_id = $3 AND ii.deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query, pq.Array(skuIds), inventoryId, tenantId)
	if err != nil {
		return inventoryItems, err
	}
	defer rows.Close()

	for rows.Next() {
		var inventoryItem domain.InventoryItem
		err = rows.Scan(&inventoryItem.Id, &inventoryItem.InventoryId, &inventoryItem.Quantity, &tenantId, &inventoryItem.Sku.Id, &inventoryItem.Sku.Code, &inventoryItem.Sku.Color, &inventoryItem.Sku.Size, &inventoryItem.Sku.Cost, &inventoryItem.Sku.Price, &inventoryItem.Sku.Product.Name)
		if err != nil {
			return inventoryItems, err
		}
		inventoryItem.Sku.Quantity = inventoryItem.Quantity
		inventoryItems = append(inventoryItems, inventoryItem)
	}
	return inventoryItems, err
}

func (r *inventoryItemRepository) GetByIdWithTransaction(ctx context.Context, tx *sql.Tx, id int64) (domain.InventoryItem, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventoryItem domain.InventoryItem

	query := `SELECT ii.id, ii.inventory_id, ii.quantity, ii.tenant_id, s.id, s.code, s.color, s.size, s.cost, s.price, p.name
	FROM inventory_items ii
	INNER JOIN skus s ON s.id = ii.sku_id
	INNER JOIN products p ON p.id = s.product_id
	WHERE ii.id = $1 AND ii.tenant_id = $2 AND ii.deleted_at IS NULL`
	err := tx.QueryRowContext(ctx, query, id, tenantId).Scan(&inventoryItem.Id, &inventoryItem.InventoryId, &inventoryItem.Quantity, &tenantId, &inventoryItem.Sku.Id, &inventoryItem.Sku.Code, &inventoryItem.Sku.Color, &inventoryItem.Sku.Size, &inventoryItem.Sku.Cost, &inventoryItem.Sku.Price, &inventoryItem.Sku.Product.Name)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return inventoryItem, domain.ErrInventoryItemNotFound
		}
		return inventoryItem, err
	}
	return inventoryItem, nil
}

func (r *inventoryItemRepository) GetAll(ctx context.Context) ([]domain.GetInventoryItemsOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventoryItems []domain.GetInventoryItemsOutput

	query := ` SELECT inv_items.id, sku.code, sku.color, sku.size, p.name, inv.type, u.name, inv_items.quantity, sku.id
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
		var inventoryItem domain.GetInventoryItemsOutput

		err = rows.Scan(&inventoryItem.InventoryItemId, &inventoryItem.SkuCode, &inventoryItem.SkuColor, &inventoryItem.SkuSize, &inventoryItem.ProductName, &inventoryItem.InventoryType, &inventoryItem.UserName, &inventoryItem.Quantity, &inventoryItem.SkuId)
		if err != nil {
			return inventoryItems, err
		}
		inventoryItems = append(inventoryItems, inventoryItem)
	}
	return inventoryItems, err
}

func (r *inventoryItemRepository) GetByInventoryId(ctx context.Context, id int64) ([]domain.GetInventoryItemsOutput, error) {
	inventoryItems := make([]domain.GetInventoryItemsOutput, 0)

	query := `SELECT inv_items.id, sku.code, sku.color, sku.size, p.name, inv.type, u.name, inv_items.quantity, sku.id
	FROM inventory_items inv_items 
	INNER JOIN inventories inv ON inv.id = inv_items.inventory_id
	INNER JOIN skus sku ON sku.id = inv_items.sku_id
	INNER JOIN products p ON p.id = sku.product_id
	LEFT JOIN users u ON u.id = inv.user_id 
	WHERE inv_items.inventory_id = $1 AND inv_items.tenant_id = $2 AND inv_items.deleted_at IS NULL ORDER BY inv_items.id ASC`
	rows, err := r.db.QueryContext(ctx, query, id, ctx.Value(constants.TENANT_KEY))
	if err != nil {
		return inventoryItems, err
	}
	defer rows.Close()

	for rows.Next() {
		var inventoryItem domain.GetInventoryItemsOutput

		err = rows.Scan(&inventoryItem.InventoryItemId, &inventoryItem.SkuCode, &inventoryItem.SkuColor, &inventoryItem.SkuSize, &inventoryItem.ProductName, &inventoryItem.InventoryType, &inventoryItem.UserName, &inventoryItem.Quantity, &inventoryItem.SkuId)
		if err != nil {
			return inventoryItems, err
		}
		inventoryItems = append(inventoryItems, inventoryItem)
	}
	return inventoryItems, err
}

func (r *inventoryItemRepository) GetBySkuId(ctx context.Context, skuId int64) ([]domain.GetSkuInventoryOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	inventories := make([]domain.GetSkuInventoryOutput, 0)

	query := `SELECT inv.id, inv.type, u.name, inv_items.quantity
        FROM inventory_items inv_items
        INNER JOIN inventories inv ON inv.id = inv_items.inventory_id
        LEFT JOIN users u ON u.id = inv.user_id
        WHERE inv_items.sku_id = $1 AND inv_items.tenant_id = $2 AND inv_items.deleted_at IS NULL
        ORDER BY inv.id ASC`

	rows, err := r.db.QueryContext(ctx, query, skuId, tenantId)
	if err != nil {
		return inventories, err
	}
	defer rows.Close()

	for rows.Next() {
		var inventory domain.GetSkuInventoryOutput
		err = rows.Scan(&inventory.InventoryId, &inventory.InventoryType, &inventory.UserName, &inventory.Quantity)
		if err != nil {
			return inventories, err
		}
		inventories = append(inventories, inventory)
	}

	return inventories, err
}
