package viewmodel

import (
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

type GetInventoryItemsViewModel struct {
	InventoryItemId int64   `json:"inventory_item_id"`
	SkuId           int64   `json:"sku_id"`
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

var inventoryTransactionTypeMap = map[domain.InventoryTransactionType]string{
	domain.InventoryTransactionTypeTransfer: "Transferência",
	domain.InventoryTransactionTypeIn:       "Entrada",
	domain.InventoryTransactionTypeOut:      "Saída",
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
		SkuId:           inventoryItem.SkuId,
		SkuCode:         *inventoryItem.SkuCode,
		ProductName:     productName,
		InventoryType:   inventoryType,
		UserName:        inventoryItem.UserName,
		Quantity:        inventoryItem.Quantity,
	}
}

type GetInventoryTransactionsViewModel struct {
	Id            int64   `json:"id"`
	Date          string  `json:"date"`
	Type          string  `json:"type"`
	Quantity      float64 `json:"quantity"`
	SkuCode       string  `json:"sku_code"`
	ProductName   string  `json:"product_name"`
	Origin        *string `json:"origin"`
	Destination   *string `json:"destination"`
	Justification *string `json:"justification"`
}

func ToGetInventoryTransactionsViewModel(inventoryTransaction output.GetInventoryTransactionsOutput) GetInventoryTransactionsViewModel {
	productName := ""
	transactionType := "Outro"

	var origin *string
	var destionation *string

	if inventoryTransaction.ProductName != "" {
		productName = inventoryTransaction.ProductName
	}

	if inventoryTransaction.SkuSize != nil {
		productName = productName + " - " + *inventoryTransaction.SkuSize
	}

	if inventoryTransaction.SkuColor != nil {
		productName = productName + " - " + *inventoryTransaction.SkuColor
	}

	if inventoryTransaction.InventoryOriginType != nil {
		origin = new(string)
		*origin = inventoryTypeMap[*inventoryTransaction.InventoryOriginType]

		if inventoryTransaction.UserOriginName != nil {
			*origin = *inventoryTransaction.UserOriginName + " - " + *origin
		}
	}

	if inventoryTransaction.InventoryDestinationType != nil {
		destionation = new(string)
		*destionation = inventoryTypeMap[*inventoryTransaction.InventoryDestinationType]

		if inventoryTransaction.UserDestinationName != nil {
			*destionation = *inventoryTransaction.UserDestinationName + " - " + *destionation
		}
	}

	if inventoryTransactionTypeMap[inventoryTransaction.Type] != "" {
		transactionType = inventoryTransactionTypeMap[inventoryTransaction.Type]
	}

	return GetInventoryTransactionsViewModel{
		Id:            inventoryTransaction.Id,
		Date:          inventoryTransaction.Date.Format("02/01/2006 15:04"),
		Type:          transactionType,
		Quantity:      inventoryTransaction.Quantity,
		ProductName:   productName,
		Origin:        origin,
		Destination:   destionation,
		SkuCode:       inventoryTransaction.SkuCode,
		Justification: inventoryTransaction.Justification,
	}
}

type GetInventoriesViewModel struct {
	Id   int64  `json:"id"`
	Type string `json:"type"`
}

func ToGetInventoriesViewModel(inventory domain.Inventory) GetInventoriesViewModel {
	inventoryType := inventoryTypeMap[inventory.Type]
	if inventory.User.Name != "" {
		inventoryType = inventory.User.Name + " - " + inventoryType
	}

	return GetInventoriesViewModel{
		Id:   inventory.Id,
		Type: inventoryType,
	}
}

type GetInventorySummaryViewModel struct {
	InventoryId       int64   `json:"inventory_id"`
	InventoryName     string  `json:"inventory_name"`
	TotalSkus         int64   `json:"total_skus"`
	TotalQuantity     float64 `json:"total_quantity"`
	ZeroQuantityItems int64   `json:"zero_quantity_items"`
}

func ToGetInventorySummaryViewModel(summary output.GetInventorySummaryOutput) GetInventorySummaryViewModel {
	inventoryName := inventoryTypeMap[summary.InventoryType]
	if summary.InventoryUserName != nil && *summary.InventoryUserName != "" {
		inventoryName = *summary.InventoryUserName + " - " + inventoryName
	}

	return GetInventorySummaryViewModel{
		InventoryId:       summary.InventoryId,
		InventoryName:     inventoryName,
		TotalSkus:         summary.TotalSkus,
		TotalQuantity:     summary.TotalQuantity,
		ZeroQuantityItems: summary.ZeroQuantityItems,
	}
}

type GetInventorySummaryByIdViewModel struct {
	InventoryId         int64   `json:"inventory_id"`
	InventoryName       string  `json:"inventory_name"`
	TotalQuantity       float64 `json:"total_quantity"`
	ZeroQuantityItems   int64   `json:"zero_quantity_items"`
	LastTransactionDays *int64  `json:"last_transaction_days"`
}

func ToGetInventorySummaryByIdViewModel(summary output.GetInventorySummaryByIdOutput) GetInventorySummaryByIdViewModel {
	inventoryName := inventoryTypeMap[summary.InventoryType]
	if summary.InventoryUserName != nil && *summary.InventoryUserName != "" {
		inventoryName = *summary.InventoryUserName + " - " + inventoryName
	}

	return GetInventorySummaryByIdViewModel{
		InventoryId:         summary.InventoryId,
		InventoryName:       inventoryName,
		TotalQuantity:       summary.TotalQuantity,
		ZeroQuantityItems:   summary.ZeroQuantityItems,
		LastTransactionDays: summary.LastTransactionDays,
	}
}
