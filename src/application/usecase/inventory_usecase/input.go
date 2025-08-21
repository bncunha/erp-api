package inventory_usecase

import "github.com/bncunha/erp-api/src/domain"

type DoTransactionInput struct {
	Type                   domain.InventoryTransactionType
	SkuId                  int64
	InventoryOriginId      int64
	InventoryDestinationId int64
	Quantity               float64
	Justification          string
}
