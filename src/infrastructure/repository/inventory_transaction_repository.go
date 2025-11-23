package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
)

type inventoryTransactionRepository struct {
	db                      *sql.DB
	intentoryItemRepository domain.InventoryItemRepository
}

func NewInventoryTransactionRepository(db *sql.DB, inventoryItemRepository domain.InventoryItemRepository) domain.InventoryTransactionRepository {
	return &inventoryTransactionRepository{db, inventoryItemRepository}
}

func (r *inventoryTransactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction domain.InventoryTransaction) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedID int64
	var nullableInventoryInId *int64
	var nullableInventoryOutId *int64
	var nullableSaleId *int64

	if transaction.InventoryIn.Id != 0 {
		nullableInventoryInId = &transaction.InventoryIn.Id
	}

	if transaction.InventoryOut.Id != 0 {
		nullableInventoryOutId = &transaction.InventoryOut.Id
	}

	if transaction.Sale.Id != 0 {
		nullableSaleId = &transaction.Sale.Id
	}

	query := `INSERT INTO inventory_transactions (quantity, type, date, inventory_out_id, inventory_in_id, inventory_item_id, tenant_id, justification, sales_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	err := tx.QueryRowContext(ctx, query, transaction.Quantity, transaction.Type, transaction.Date, nullableInventoryOutId, nullableInventoryInId, transaction.InventoryItem.Id, tenantId, transaction.Justification, nullableSaleId).Scan(&insertedID)
	return insertedID, err
}

func (r *inventoryTransactionRepository) GetAll(ctx context.Context) ([]domain.GetInventoryTransactionsOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventoryTransactions []domain.GetInventoryTransactionsOutput

	query := `SELECT
	inv_transactions.id,
	inv_transactions.date, 
	inv_transactions.type, 
	inv_transactions.quantity, 
	sku.code, 
	sku.color, 
	sku.size, 
	p.name, 
	inventory_origin.type, 
	inventory_destination.type, 
	user_origin.name, 
	user_destination.name, 
	inv_transactions.justification,
	s.id,
	s.date
	FROM inventory_transactions inv_transactions
	INNER JOIN inventory_items inv_items ON inv_transactions.inventory_item_id = inv_items.id
	INNER JOIN skus sku ON sku.id = inv_items.sku_id
	INNER JOIN products p ON p.id = sku.product_id
	LEFT JOIN inventories inventory_origin ON inventory_origin.id = inv_transactions.inventory_out_id
	LEFT JOIN inventories inventory_destination ON inventory_destination.id = inv_transactions.inventory_in_id
	LEFT JOIN users user_origin ON user_origin.id = inventory_origin.user_id
	LEFT JOIN users user_destination ON user_destination.id = inventory_destination.user_id
	LEFT JOIN sales s ON s.id = inv_transactions.sales_id
        WHERE inv_transactions.tenant_id = $1 AND inv_transactions.deleted_at IS NULL ORDER BY inv_transactions.date DESC`
	rows, err := r.db.QueryContext(ctx, query, tenantId)
	if err != nil {
		return inventoryTransactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var inventoryTransaction domain.GetInventoryTransactionsOutput
		saleId := sql.NullInt64{}
		saleDate := sql.NullTime{}

		err = rows.Scan(&inventoryTransaction.Id, &inventoryTransaction.Date, &inventoryTransaction.Type, &inventoryTransaction.Quantity, &inventoryTransaction.SkuCode, &inventoryTransaction.SkuColor, &inventoryTransaction.SkuSize, &inventoryTransaction.ProductName, &inventoryTransaction.InventoryOriginType, &inventoryTransaction.InventoryDestinationType, &inventoryTransaction.UserOriginName, &inventoryTransaction.UserDestinationName, &inventoryTransaction.Justification, &saleId, &saleDate)
		if err != nil {
			return inventoryTransactions, err
		}
		if saleId.Valid {
			inventoryTransaction.Sale = &domain.Sales{Id: saleId.Int64, Date: saleDate.Time}
		}
		inventoryTransactions = append(inventoryTransactions, inventoryTransaction)
	}
	return inventoryTransactions, err
}

func (r *inventoryTransactionRepository) GetByInventoryId(ctx context.Context, inventoryId int64) ([]domain.GetInventoryTransactionsOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventoryTransactions []domain.GetInventoryTransactionsOutput

	query := `SELECT
        inv_transactions.id,
        inv_transactions.date,
        inv_transactions.type,
        inv_transactions.quantity,
        sku.code,
        sku.color,
        sku.size,
        p.name,
        inventory_origin.type,
        inventory_destination.type,
        user_origin.name,
        user_destination.name,
        inv_transactions.justification,
        s.id,
        s.date
        FROM inventory_transactions inv_transactions
        INNER JOIN inventory_items inv_items ON inv_transactions.inventory_item_id = inv_items.id
        INNER JOIN skus sku ON sku.id = inv_items.sku_id
        INNER JOIN products p ON p.id = sku.product_id
        LEFT JOIN inventories inventory_origin ON inventory_origin.id = inv_transactions.inventory_out_id
        LEFT JOIN inventories inventory_destination ON inventory_destination.id = inv_transactions.inventory_in_id
        LEFT JOIN users user_origin ON user_origin.id = inventory_origin.user_id
        LEFT JOIN users user_destination ON user_destination.id = inventory_destination.user_id
        LEFT JOIN sales s ON s.id = inv_transactions.sales_id
        WHERE inv_transactions.tenant_id = $1
        AND inv_transactions.deleted_at IS NULL
        AND (
                inv_transactions.inventory_in_id = $2 OR
                inv_transactions.inventory_out_id = $2 OR
                inv_items.inventory_id = $2
        )
        ORDER BY inv_transactions.date DESC`
	rows, err := r.db.QueryContext(ctx, query, tenantId, inventoryId)
	if err != nil {
		return inventoryTransactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var inventoryTransaction domain.GetInventoryTransactionsOutput
		saleId := sql.NullInt64{}
		saleDate := sql.NullTime{}

		err = rows.Scan(&inventoryTransaction.Id, &inventoryTransaction.Date, &inventoryTransaction.Type, &inventoryTransaction.Quantity, &inventoryTransaction.SkuCode, &inventoryTransaction.SkuColor, &inventoryTransaction.SkuSize, &inventoryTransaction.ProductName, &inventoryTransaction.InventoryOriginType, &inventoryTransaction.InventoryDestinationType, &inventoryTransaction.UserOriginName, &inventoryTransaction.UserDestinationName, &inventoryTransaction.Justification, &saleId, &saleDate)
		if err != nil {
			return inventoryTransactions, err
		}
		if saleId.Valid {
			inventoryTransaction.Sale = &domain.Sales{Id: saleId.Int64, Date: saleDate.Time}
		}
		inventoryTransactions = append(inventoryTransactions, inventoryTransaction)
	}
	return inventoryTransactions, err
}
