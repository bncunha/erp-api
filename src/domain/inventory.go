package domain

import "time"

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
	UserId   int64
	TenantId int64
	Type     InventoryType
	Items    []InventoryItem
}

type InventoryItem struct {
	Id          int64
	InventoryId int64
	SkuId       int64
	Quantity    float64
	TenantId    int64
}

type InventoryTransaction struct {
	Id            int64
	Quantity      float64
	Type          InventoryTransactionType
	Date          time.Time
	InventoryIn   Inventory
	InventoryOut  Inventory
	InventoryItem InventoryItem
	TenantId      int64
	Justification string
}
