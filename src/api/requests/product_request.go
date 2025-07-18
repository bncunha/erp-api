package request

import (
	"github.com/bncunha/erp-api/src/application/validator"
)

type CreateProductRequest struct {
	Name         string                 `json:"name" validate:"required,max=200"`
	Description  string                 `json:"description" validate:"max=500"`
	CategoryID   int64                 	`json:"categoryId"`
	CategoryName string                 `json:"categoryName" validate:"max=200"`
	Skus         []CreateSkuRequest `json:"skus"`
}

func (r *CreateProductRequest) Validate() error {
	err := validator.Validate(r)
	if err != nil {
		return err
	}
	for _, sku := range r.Skus {
		err = sku.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

type EditProductRequest struct {
	Id int64 `json:"id" validate:"required"`
	Name         string                 `json:"name" validate:"required,max=200"`
	Description  string                 `json:"description" validate:"max=500"`
	CategoryID   int64                 	`json:"categoryId"`
	CategoryName string                 `json:"categoryName" validate:"max=200"`
}