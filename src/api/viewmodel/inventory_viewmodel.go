package viewmodel

import (
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

type GetInventoryItemsViewModel struct {
	InventoryItemId int64   `json:"inventory_item_id"`
	SkuCode         string  `json:"sku_code"`
	ProductName     string  `json:"product_name"`
	InventoryType   string  `json:"inventory_type"`
	UserName        *string `json:"user_name"`
	Quantity        float64 `json:"quantity"`
}

var inventoryTypeMap = map[domain.InventoryType]string{
	domain.InventoryTypePrimary:  "Central",
	domain.InventoryTypeReseller: "Revendedor",
}

func ToGetInventoryItemsViewModel(inventoryItem output.GetInventoryItemsOutput) GetInventoryItemsViewModel {
	productName := ""
	inventoryType := ""

	if inventoryItem.ProductName != nil {
		productName = *inventoryItem.ProductName
	}

	if inventoryItem.SkuColor != nil {
		productName = productName + " - " + *inventoryItem.SkuColor
	}

	if inventoryItem.SkuSize != nil {
		productName = productName + " - " + *inventoryItem.SkuSize
	}

	if inventoryItem.InventoryType != nil {
		inventoryType = inventoryTypeMap[*inventoryItem.InventoryType]
		if inventoryType == "" {
			inventoryType = "Outro"
		}
	}

	return GetInventoryItemsViewModel{
		InventoryItemId: inventoryItem.InventoryItemId,
		SkuCode:         *inventoryItem.SkuCode,
		ProductName:     productName,
		InventoryType:   inventoryType,
		UserName:        inventoryItem.UserName,
		Quantity:        inventoryItem.Quantity,
	}
}
