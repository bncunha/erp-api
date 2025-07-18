package viewmodel

import "github.com/bncunha/erp-api/src/domain"

type GetProductViewModel struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryId  int64 `json:"categoryId,omitempty"`
	Skus        []SkuViewModel `json:"skus,omitempty"`
}

func ToGetProductViewModel(product domain.Product) GetProductViewModel {
	var categoryId int64;
	skuViewModel := []SkuViewModel{}
	for _, sku := range product.Skus {
		skuViewModel = append(skuViewModel, ToSkuViewModel(sku))
	}

	if product.Category.Id != 0 {
		categoryId = product.Category.Id
	}

	return GetProductViewModel{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		CategoryId:  categoryId,
		Skus:        skuViewModel,
	}	
}