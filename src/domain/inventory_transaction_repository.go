package domain

import (
	"context"
	"database/sql"
	"time"
)

type GetInventoryTransactionsOutput struct {
	Id                       int64
	Date                     time.Time
	Type                     InventoryTransactionType
	Quantity                 float64
	SkuCode                  string
	SkuColor                 *string
	SkuSize                  *string
	ProductName              string
	InventoryOriginType      *InventoryType
	InventoryDestinationType *InventoryType
	Sale                     *Sales
	UserOriginName           *string
	UserDestinationName      *string
	Justification            *string
}

type InventoryTransactionRepository interface {
	Create(ctx context.Context, tx *sql.Tx, transaction InventoryTransaction) (int64, error)
	GetAll(ctx context.Context) ([]GetInventoryTransactionsOutput, error)
	GetByInventoryId(ctx context.Context, inventoryId int64) ([]GetInventoryTransactionsOutput, error)
}
