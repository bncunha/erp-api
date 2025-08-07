package request

import "github.com/bncunha/erp-api/src/application/validator"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *LoginRequest) Validate() error {
	return validator.Validate(r)
}