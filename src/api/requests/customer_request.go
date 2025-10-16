package request

import "github.com/bncunha/erp-api/src/application/validator"

type CreateCustomerRequest struct {
    Name      string `json:"name" validate:"required,max=200"`
    Cellphone string `json:"cellphone" validate:"required,max=20"`
}

func (r *CreateCustomerRequest) Validate() error {
    return validator.Validate(r)
}

type EditCustomerRequest struct {
    Id        int64  `json:"id" validate:"required"`
    Name      string `json:"name" validate:"required,max=200"`
    Cellphone string `json:"cellphone" validate:"required,max=20"`
}

func (r *EditCustomerRequest) Validate() error {
    return validator.Validate(r)
}
