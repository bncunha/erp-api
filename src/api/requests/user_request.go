package request

import (
	"github.com/bncunha/erp-api/src/application/validator"
	"github.com/bncunha/erp-api/src/domain"
)

type CreateUserRequest struct {
	Username    string `json:"username" validate:"required,max=30"`
	Name        string `json:"name" validate:"required,max=100"`
	PhoneNumber string `json:"phone_number" validate:"max=20"`
	Role        string `json:"role" validate:"required,max=100"`
	Email       string `json:"email" validate:"required,email,max=250"`
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
	Role        string `json:"role" validate:"required,max=100"`
	Email       string `json:"email" validate:"required,email,max=250"`
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

type ResetPasswordRequest struct {
	Code     string `json:"code" validate:"required"`
	Uuid     string `json:"uuid" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *ResetPasswordRequest) Validate() error {
	err := validator.Validate(r)
	if err != nil {
		return err
	}
	return nil
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (r *ForgotPasswordRequest) Validate() error {
	err := validator.Validate(r)
	if err != nil {
		return err
	}
	return nil
}
