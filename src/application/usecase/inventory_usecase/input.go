package inventory_usecase

import "github.com/bncunha/erp-api/src/domain"

type DoTransactionInput struct {
	Type                   domain.InventoryTransactionType
	Skus                   []DoTransactionSkusInput
	InventoryOriginId      int64
	InventoryDestinationId int64
	Justification          string
	Sale                   domain.Sales
}

type DoTransactionSkusInput struct {
	SkuId    int64
	Quantity float64
}
