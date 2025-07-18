package request

import (
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/validator"
)

type CreateProductRequest struct {
	Name         string                 `json:"name" validate:"required,max=200"`
	Description  string                 `json:"description" validate:"max=500"`
	CategoryID   int64                 	`json:"categoryId"`
	CategoryName string                 `json:"categoryName" validate:"max=200"`
	Skus         []CreateProductSkuRequest `json:"skus"`
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

type CreateProductSkuRequest struct {
	Code  string  `json:"code" validate:"required,max=20"`
	Color string  `json:"color" validate:"max=200"`
	Size  string  `json:"size" validate:"max=200"`
	Cost  *float64 `json:"cost" validate:"omitempty,gt=0"`
	Price *float64 `json:"price" validate:"omitempty,gt=0"`
}

func (r *CreateProductSkuRequest) Validate() error {
	if r.Color == "" && r.Size == "" {
		return errors.New("Cor ou Tamanho são obrigatórios") 
	}
	err := validator.Validate(r)
	if err != nil {
		return err
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