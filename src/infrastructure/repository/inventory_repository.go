package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type inventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) domain.InventoryRepository {
	return &inventoryRepository{db}
}

func (r *inventoryRepository) Create(ctx context.Context, inventory domain.Inventory) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedID int64

	query := `INSERT INTO inventories (user_id, tenant_id, type) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, inventory.User.Id, tenantId, inventory.Type).Scan(&insertedID)
	return insertedID, err
}

func (r *inventoryRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, inventory domain.Inventory) (int64, error) {
	tenantId := inventory.TenantId
	if tenantId == 0 {
		tenantId = ctx.Value(constants.TENANT_KEY).(int64)
	}

	var insertedID int64
	query := `INSERT INTO inventories (user_id, tenant_id, type) VALUES ($1, $2, $3) RETURNING id`
	err := tx.QueryRowContext(ctx, query, inventory.User.Id, tenantId, inventory.Type).Scan(&insertedID)
	return insertedID, err
}

func (r *inventoryRepository) GetById(ctx context.Context, id int64) (domain.Inventory, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventory domain.Inventory

	query := `SELECT id, tenant_id, type FROM inventories WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&inventory.Id, &inventory.TenantId, &inventory.Type)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return inventory, domain.ErrInventoryNotFound
		}
		return inventory, err
	}

	return inventory, nil
}

func (r *inventoryRepository) GetAll(ctx context.Context) ([]domain.Inventory, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventories []domain.Inventory

	query := `SELECT i.id, u.id, u.name, i.tenant_id, i.type
	FROM inventories i
	LEFT JOIN users u ON u.id = i.user_id
	WHERE i.tenant_id = $1 AND i.deleted_at IS NULL ORDER BY i.id ASC`
	rows, err := r.db.QueryContext(ctx, query, tenantId)
	if err != nil {
		return inventories, err
	}
	defer rows.Close()

	for rows.Next() {
		var inventory domain.Inventory
		var nullableUserId sql.NullInt64
		var userName sql.NullString

		err = rows.Scan(&inventory.Id, &nullableUserId, &userName, &inventory.TenantId, &inventory.Type)
		if err != nil {
			return inventories, err
		}
		if nullableUserId.Valid {
			inventory.User.Id = nullableUserId.Int64
			inventory.User.Name = userName.String
		}
		inventories = append(inventories, inventory)
	}
	return inventories, err
}

func (r *inventoryRepository) GetByUserId(ctx context.Context, userId int64) (domain.Inventory, error) {
	var inventory domain.Inventory

	query := `SELECT i.id, u.id, u.name, i.tenant_id, i.type
        FROM inventories i
        LEFT JOIN users u ON u.id = i.user_id
        WHERE i.user_id = $1 AND i.deleted_at IS NULL AND i.tenant_id = $2 ORDER BY i.id ASC`
	err := r.db.QueryRowContext(ctx, query, userId, ctx.Value(constants.TENANT_KEY)).Scan(&inventory.Id, &inventory.User.Id, &inventory.User.Name, &inventory.TenantId, &inventory.Type)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return inventory, domain.ErrInventoryNotFound
		}
		return inventory, err
	}
	return inventory, nil
}

func (r *inventoryRepository) GetPrimaryInventory(ctx context.Context) (domain.Inventory, error) {
	var inventory domain.Inventory

	query := `SELECT i.id, i.tenant_id, i.type
        FROM inventories i
        WHERE i.deleted_at IS NULL AND i.tenant_id = $1 AND i.type = $2 ORDER BY i.id ASC`
	err := r.db.QueryRowContext(ctx, query, ctx.Value(constants.TENANT_KEY), domain.InventoryTypePrimary).Scan(&inventory.Id, &inventory.TenantId, &inventory.Type)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return inventory, domain.ErrInventoryNotFound
		}
		return inventory, err
	}
	return inventory, nil
}

func (r *inventoryRepository) GetSummary(ctx context.Context) ([]domain.GetInventorySummaryOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	summaries := make([]domain.GetInventorySummaryOutput, 0)

	query := `SELECT i.id,
       i.type,
       u.name,
       COUNT(DISTINCT ii.sku_id) AS total_skus,
       COALESCE(SUM(ii.quantity), 0) AS total_quantity,
       COALESCE(SUM(CASE WHEN ii.quantity = 0 THEN 1 ELSE 0 END), 0) AS zero_quantity_items
FROM inventories i
LEFT JOIN users u ON u.id = i.user_id
LEFT JOIN inventory_items ii ON ii.inventory_id = i.id AND ii.deleted_at IS NULL AND ii.tenant_id = i.tenant_id
WHERE i.tenant_id = $1 AND i.deleted_at IS NULL
GROUP BY i.id, i.type, u.name
ORDER BY i.id ASC`

	rows, err := r.db.QueryContext(ctx, query, tenantId)
	if err != nil {
		return summaries, err
	}
	defer rows.Close()

	for rows.Next() {
		var summary domain.GetInventorySummaryOutput
		var inventoryType string
		var userName sql.NullString

		if err := rows.Scan(&summary.InventoryId, &inventoryType, &userName, &summary.TotalSkus, &summary.TotalQuantity, &summary.ZeroQuantityItems); err != nil {
			return summaries, err
		}

		summary.InventoryType = domain.InventoryType(inventoryType)
		if userName.Valid {
			name := userName.String
			summary.InventoryUserName = &name
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (r *inventoryRepository) GetSummaryById(ctx context.Context, id int64) (domain.GetInventorySummaryByIdOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var summary domain.GetInventorySummaryByIdOutput
	var inventoryType string
	var userName sql.NullString

	query := `SELECT i.id,
   i.type,
   u.name,
   COUNT(DISTINCT ii.sku_id) AS total_skus,
   COALESCE(SUM(ii.quantity), 0) AS total_quantity,
   COALESCE(SUM(CASE WHEN ii.quantity = 0 THEN 1 ELSE 0 END), 0) AS zero_quantity_items,
   (
           SELECT DATE_PART('day', NOW() - MAX(it.date))::bigint
           FROM inventory_transactions it
           WHERE it.tenant_id = i.tenant_id
             AND it.deleted_at IS NULL
             AND (it.inventory_in_id = i.id OR it.inventory_out_id = i.id)
   ) AS last_transaction_days
FROM inventories i
LEFT JOIN users u ON u.id = i.user_id
LEFT JOIN inventory_items ii ON ii.inventory_id = i.id AND ii.deleted_at IS NULL AND ii.tenant_id = i.tenant_id
WHERE i.id = $1 AND i.tenant_id = $2 AND i.deleted_at IS NULL
GROUP BY i.id, i.type, u.name`

	var lastTransactionDays sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&summary.InventoryId, &inventoryType, &userName, &summary.TotalSkus, &summary.TotalQuantity, &summary.ZeroQuantityItems, &lastTransactionDays)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return summary, domain.ErrInventoryNotFound
		}
		return summary, err
	}

	summary.InventoryType = domain.InventoryType(inventoryType)
	if userName.Valid {
		name := userName.String
		summary.InventoryUserName = &name
	}
	if lastTransactionDays.Valid {
		days := lastTransactionDays.Int64
		summary.LastTransactionDays = &days
	}

	return summary, nil
}
