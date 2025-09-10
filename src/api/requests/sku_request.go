package request

import (
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/validator"
)

type CreateSkuRequest struct {
	Code          string   `json:"code" validate:"required,max=20"`
	Color         string   `json:"color" validate:"max=200"`
	Size          string   `json:"size" validate:"max=200"`
	Cost          *float64 `json:"cost" validate:"omitempty,gt=0"`
	Price         float64  `json:"price" validate:"omitempty,gt=0"`
	Quantity      *float64 `json:"quantity" validate:"omitempty,gt=0"`
	DestinationId *int64   `json:"destination_id" validate:"omitempty,gt=0"`
}

func (r *CreateSkuRequest) Validate() error {
	if r.Color == "" && r.Size == "" {
		return errors.New("Cor ou Tamanho são obrigatórios")
	}
	err := validator.Validate(r)
	if r.Quantity != nil && r.DestinationId == nil {
		return errors.New("Necessário informar o estoque de destino")
	}
	if r.Quantity == nil && r.DestinationId != nil {
		return errors.New("Necessário informar a quantidade")
	}
	if err != nil {
		return err
	}
	return nil
}

type EditSkuRequest struct {
	CreateSkuRequest
}
