package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
)

type InventoryTransactionRepository interface {
	Create(ctx context.Context, tx *sql.Tx, transaction domain.InventoryTransaction) (int64, error)
}

type inventoryTransactionRepository struct {
	db                      *sql.DB
	intentoryItemRepository InventoryItemRepository
}

func NewInventoryTransactionRepository(db *sql.DB, inventoryItemRepository InventoryItemRepository) InventoryTransactionRepository {
	return &inventoryTransactionRepository{db, inventoryItemRepository}
}

func (r *inventoryTransactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction domain.InventoryTransaction) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedID int64
	var nullableInventoryInId *int64
	var nullableInventoryOutId *int64

	if transaction.InventoryIn.Id != 0 {
		nullableInventoryInId = &transaction.InventoryIn.Id
	}

	if transaction.InventoryOut.Id != 0 {
		nullableInventoryOutId = &transaction.InventoryOut.Id
	}

	query := `INSERT INTO inventory_transactions (quantity, type, date, inventory_out_id, inventory_in_id, inventory_item_id, tenant_id, justification) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err := tx.QueryRowContext(ctx, query, transaction.Quantity, transaction.Type, transaction.Date, nullableInventoryOutId, nullableInventoryInId, transaction.InventoryItem.Id, tenantId, transaction.Justification).Scan(&insertedID)
	return insertedID, err
}
