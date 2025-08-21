package output

import "github.com/bncunha/erp-api/src/domain"

type GetInventoryItemsOutput struct {
	InventoryItemId int64
	SkuCode         *string
	SkuColor        *string
	SkuSize         *string
	ProductName     *string
	InventoryType   *domain.InventoryType
	UserName        *string
	Quantity        float64
}
