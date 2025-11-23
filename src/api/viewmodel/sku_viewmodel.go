package viewmodel

import (
	"time"

	"github.com/bncunha/erp-api/src/domain"
)

type SkuViewModel struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	ProductName string   `json:"product_name"`
	Code        string   `json:"code"`
	Color       string   `json:"color"`
	Size        string   `json:"size"`
	Cost        *float64 `json:"cost"`
	Price       *float64 `json:"price"`
	Quantity    float64  `json:"quantity"`
}

func ToSkuViewModel(sku domain.Sku) SkuViewModel {
	return SkuViewModel{
		Id:          sku.Id,
		Name:        sku.GetName(),
		ProductName: sku.Product.Name,
		Code:        sku.Code,
		Color:       sku.Color,
		Size:        sku.Size,
		Cost:        sku.Cost,
		Price:       &sku.Price,
		Quantity:    sku.Quantity,
	}
}

type SkuInventoryViewModel struct {
	InventoryName string  `json:"inventory_name"`
	Quantity      float64 `json:"quantity"`
}

func ToSkuInventoryViewModel(inventory domain.GetSkuInventoryOutput) SkuInventoryViewModel {
	return SkuInventoryViewModel{
		InventoryName: formatInventoryName(inventory.InventoryType, inventory.UserName),
		Quantity:      inventory.Quantity,
	}
}

type SkuTransactionViewModel struct {
	Date          string  `json:"date"`
	Type          string  `json:"type"`
	Quantity      float64 `json:"quantity"`
	Origin        *string `json:"origin"`
	Destination   *string `json:"destination"`
	Justification *string `json:"justification"`
}

func ToSkuTransactionViewModel(transaction domain.GetInventoryTransactionsOutput) SkuTransactionViewModel {
	transactionType := "Outro"

	var origin *string
	var destination *string

	originName := formatInventoryName(transaction.InventoryOriginType, transaction.UserOriginName)
	if originName != "" {
		origin = new(string)
		*origin = originName
	}

	destinationName := formatInventoryName(transaction.InventoryDestinationType, transaction.UserDestinationName)
	if destinationName != "" {
		destination = new(string)
		*destination = destinationName
	}

	if inventoryTransactionTypeMap[transaction.Type] != "" {
		transactionType = inventoryTransactionTypeMap[transaction.Type]
	}

	loc, _ := time.LoadLocation("America/Sao_Paulo")

	return SkuTransactionViewModel{
		Date:          transaction.Date.In(loc).Format("02/01/2006 15:04"),
		Type:          transactionType,
		Quantity:      transaction.Quantity,
		Origin:        origin,
		Destination:   destination,
		Justification: transaction.Justification,
	}
}

func formatInventoryName(inventoryType *domain.InventoryType, userName *string) string {
	inventoryName := ""

	if inventoryType != nil {
		inventoryName = inventoryTypeMap[*inventoryType]
		if inventoryName == "" {
			inventoryName = "Outro"
		}
	}

	if userName != nil {
		if inventoryName != "" {
			inventoryName = *userName + " - " + inventoryName
		} else {
			inventoryName = *userName
		}
	}

	return inventoryName
}
