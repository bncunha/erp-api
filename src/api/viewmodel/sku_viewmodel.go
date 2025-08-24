package viewmodel

import (
	"github.com/bncunha/erp-api/src/domain"
)

type SkuViewModel struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	Code     string   `json:"code"`
	Color    string   `json:"color"`
	Size     string   `json:"size"`
	Cost     *float64 `json:"cost"`
	Price    *float64 `json:"price"`
	Quantity float64  `json:"quantity"`
}

func ToSkuViewModel(sku domain.Sku) SkuViewModel {
	return SkuViewModel{
		Id:       sku.Id,
		Name:     sku.GetName(),
		Code:     sku.Code,
		Color:    sku.Color,
		Size:     sku.Size,
		Cost:     sku.Cost,
		Price:    sku.Price,
		Quantity: sku.Quantity,
	}
}
