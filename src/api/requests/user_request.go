package request

import (
	"github.com/bncunha/erp-api/src/application/validator"
	"github.com/bncunha/erp-api/src/domain"
)

type CreateUserRequest struct {
	Username    string `json:"username" validate:"required,max=30"`
	Name        string `json:"name" validate:"required,max=100"`
	PhoneNumber string `json:"phone_number" validate:"max=20"`
	Password    string `json:"password" validate:"required,max=20"`
	Role        string `json:"role" validate:"required,max=100"`
}

func (r *CreateUserRequest) Validate() error {
	err := validator.Validate(r)
	if err != nil {
		return err
	}
	return nil
}

type EditUserRequest struct {
	Username    string `json:"username" validate:"required,max=30"`
	Name        string `json:"name" validate:"required,max=100"`
	PhoneNumber string `json:"phone_number" validate:"max=20"`
	Password    string `json:"password" validate:"max=20"`
	Role        string `json:"role" validate:"required,max=100"`
}

func (r *EditUserRequest) Validate() error {
	err := validator.Validate(r)
	if err != nil {
		return err
	}
	return nil
}

type GetAllUserRequest struct {
	Role domain.Role `json:"role"`
}
