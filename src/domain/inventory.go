package domain

import (
	"time"
)

type InventoryType string

type InventoryTransactionType string

const (
	InventoryTypePrimary  InventoryType = "PRIMARY"
	InventoryTypeReseller InventoryType = "RESELLER"
)

const (
	InventoryTransactionTypeTransfer InventoryTransactionType = "TRANSFER"
	InventoryTransactionTypeIn       InventoryTransactionType = "IN"
	InventoryTransactionTypeOut      InventoryTransactionType = "OUT"
)

type Inventory struct {
	Id       int64
	TenantId int64
	User     User
	Type     InventoryType
	Items    []InventoryItem
}

type InventoryItem struct {
	Id          int64
	InventoryId int64
	Sku         Sku
	Quantity    float64
	TenantId    int64
}

func NewInventoryItem(inventoryId int64, sku Sku, quantity float64) InventoryItem {
	return InventoryItem{
		InventoryId: inventoryId,
		Sku:         sku,
		Quantity:    quantity,
	}
}

type InventoryTransaction struct {
	Id            int64
	Quantity      float64
	Type          InventoryTransactionType
	Date          time.Time
	InventoryIn   Inventory
	InventoryOut  Inventory
	InventoryItem InventoryItem
	Sale          Sales
	TenantId      int64
	Justification string
}
