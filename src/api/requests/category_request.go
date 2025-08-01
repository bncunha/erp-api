package request

import "github.com/bncunha/erp-api/src/application/validator"

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,max=200"`
}

func (r *CreateCategoryRequest) Validate() error {
	return validator.Validate(r)
}

type EditCategoryRequest struct {
	Id   int64  `json:"id" validate:"required"`
	CreateCategoryRequest
}

func (r *EditCategoryRequest) Validate() error {
	return validator.Validate(r)
}