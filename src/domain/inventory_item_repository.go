package domain

import (
	"context"
	"database/sql"
	"errors"
)

var ErrInventoryItemNotFound = errors.New("Item de estoque não encontrado")

type GetInventoryItemsOutput struct {
	InventoryItemId int64
	SkuId           int64
	SkuCode         *string
	SkuColor        *string
	SkuSize         *string
	ProductName     *string
	InventoryType   *InventoryType
	UserName        *string
	Quantity        float64
}

type InventoryItemRepository interface {
	Create(ctx context.Context, tx *sql.Tx, inventoryItem InventoryItem) (int64, error)
	UpdateQuantity(ctx context.Context, tx *sql.Tx, inventoryItem InventoryItem) error
	GetById(ctx context.Context, id int64) (InventoryItem, error)
	GetByIdWithTransaction(ctx context.Context, tx *sql.Tx, id int64) (InventoryItem, error)
	GetByManySkuIdsAndInventoryId(ctx context.Context, skuIds []int64, inventoryId int64) ([]InventoryItem, error)
	GetAll(ctx context.Context) ([]GetInventoryItemsOutput, error)
	GetByInventoryId(ctx context.Context, id int64) ([]GetInventoryItemsOutput, error)
}
