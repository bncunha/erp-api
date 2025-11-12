package domain

import (
	"context"
	"errors"
)

var ErrInventoryNotFound = errors.New("Inventário não encontrado")

type GetInventorySummaryOutput struct {
	InventoryId       int64
	InventoryType     InventoryType
	InventoryUserName *string
	TotalSkus         int64
	TotalQuantity     float64
	ZeroQuantityItems int64
}

type GetInventorySummaryByIdOutput struct {
	InventoryId         int64
	InventoryType       InventoryType
	InventoryUserName   *string
	TotalSkus           int64
	TotalQuantity       float64
	ZeroQuantityItems   int64
	LastTransactionDays *int64
}

type InventoryRepository interface {
	Create(ctx context.Context, inventory Inventory) (int64, error)
	GetById(ctx context.Context, id int64) (Inventory, error)
	GetAll(ctx context.Context) ([]Inventory, error)
	GetByUserId(ctx context.Context, userId int64) (Inventory, error)
	GetPrimaryInventory(ctx context.Context) (Inventory, error)
	GetSummary(ctx context.Context) ([]GetInventorySummaryOutput, error)
	GetSummaryById(ctx context.Context, id int64) (GetInventorySummaryByIdOutput, error)
}
