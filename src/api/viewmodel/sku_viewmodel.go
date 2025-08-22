package viewmodel

import (
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

type SkuViewModel struct {
	Id       int64    `json:"id"`
	Code     string   `json:"code"`
	Color    string   `json:"color"`
	Size     string   `json:"size"`
	Cost     *float64 `json:"cost"`
	Price    *float64 `json:"price"`
	Quantity float64  `json:"quantity"`
}

func ToGetSkusByProductViewModel(output output.GetAllSkusByProductOutput) SkuViewModel {
	sku := ToSkuViewModel(output.Sku)
	return SkuViewModel{
		Id:       sku.Id,
		Code:     sku.Code,
		Color:    sku.Color,
		Size:     sku.Size,
		Cost:     sku.Cost,
		Price:    sku.Price,
		Quantity: output.Quantity,
	}
}

func ToSkuViewModel(sku domain.Sku) SkuViewModel {
	return SkuViewModel{
		Id:    sku.Id,
		Code:  sku.Code,
		Color: sku.Color,
		Size:  sku.Size,
		Cost:  sku.Cost,
		Price: sku.Price,
	}
}
