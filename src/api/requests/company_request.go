package request

import (
	"github.com/bncunha/erp-api/src/application/errors"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/application/validator"
)

type CreateCompanyRequest struct {
	Name            string                   `json:"name" validate:"required,max=255"`
	LegalName       string                   `json:"legal_name" validate:"required,max=255"`
	Cnpj            string                   `json:"cnpj" validate:"omitempty,max=20"`
	Cpf             string                   `json:"cpf" validate:"omitempty,max=14"`
	Cellphone       string                   `json:"cellphone" validate:"required,max=20"`
	AcceptedTerms   bool                     `json:"accepted_terms" validate:"required"`
	AcceptedPrivacy bool                     `json:"accepted_privacy" validate:"required"`
	Address         CreateCompanyAddress     `json:"address" validate:"required"`
	User            CreateCompanyUserRequest `json:"user" validate:"required"`
}

type CreateCompanyAddress struct {
	Street       string `json:"street" validate:"required,max=255"`
	Neighborhood string `json:"neighborhood" validate:"required,max=255"`
	Number       string `json:"number" validate:"required,max=50"`
	City         string `json:"city" validate:"required,max=255"`
	UF           string `json:"uf" validate:"required,len=2"`
	Cep          string `json:"cep" validate:"required,max=20"`
}

type CreateCompanyUserRequest struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Username    string  `json:"username" validate:"required,max=30"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,max=20"`
	Email       string  `json:"email" validate:"required,email,max=250"`
	Password    string  `json:"password" validate:"required,min=6"`
}

func (r *CreateCompanyRequest) Validate() error {
	if err := validator.Validate(r); err != nil {
		return err
	}

	hasCnpj := r.Cnpj != ""
	hasCpf := r.Cpf != ""
	if hasCnpj == hasCpf {
		return errors.New("Envie apenas CPF ou CNPJ.")
	}

	if hasCnpj {
		if !helper.IsValidCNPJ(r.Cnpj) {
			return errors.New("CNPJ inválido")
		}
		r.Cnpj = helper.SanitizeDocument(r.Cnpj)
	}

	if hasCpf {
		if !helper.IsValidCPF(r.Cpf) {
			return errors.New("CPF inválido")
		}
		r.Cpf = helper.SanitizeDocument(r.Cpf)
	}

	if !r.AcceptedTerms || !r.AcceptedPrivacy {
		return errors.New("Os termos e a politica de privacidade são obrigatórios.")
	}

	return nil
}
