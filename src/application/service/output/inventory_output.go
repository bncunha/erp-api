package output

import (
	"time"

	"github.com/bncunha/erp-api/src/domain"
)

type GetInventoryItemsOutput struct {
	InventoryItemId int64
	SkuId           int64
	SkuCode         *string
	SkuColor        *string
	SkuSize         *string
	ProductName     *string
	InventoryType   *domain.InventoryType
	UserName        *string
	Quantity        float64
}

type GetInventoryTransactionsOutput struct {
	Id                       int64
	Date                     time.Time
	Type                     domain.InventoryTransactionType
	Quantity                 float64
	SkuCode                  string
	SkuColor                 *string
	SkuSize                  *string
	ProductName              string
	InventoryOriginType      *domain.InventoryType
	InventoryDestinationType *domain.InventoryType
	Sale                     *domain.Sales
	UserOriginName           *string
	UserDestinationName      *string
	Justification            *string
}
