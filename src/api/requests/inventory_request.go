package request

import (
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/validator"
	"github.com/bncunha/erp-api/src/domain"
)

type CreateInventoryTransactionRequest struct {
	Type                   domain.InventoryTransactionType         `json:"type" validate:"required,oneof=TRANSFER IN OUT"`
	Skus                   []CreateInventoryTransactionSkusRequest `json:"skus" validate:"required,gt=0"`
	InventoryOriginId      int64                                   `json:"inventory_origin_id"`
	InventoryDestinationId int64                                   `json:"inventory_destination_id"`
	Justification          string                                  `json:"justification" validate:"max=200"`
}

type CreateInventoryTransactionSkusRequest struct {
	SkuId    int64   `json:"sku_id" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required,gt=0"`
}

func (r *CreateInventoryTransactionRequest) Validate() error {
	switch r.Type {
	case domain.InventoryTransactionTypeTransfer:
		if r.InventoryOriginId == 0 || r.InventoryDestinationId == 0 {
			return errors.New("Origem e Destino são obrigatórios")
		}
	case domain.InventoryTransactionTypeIn:
		if r.InventoryDestinationId == 0 {
			return errors.New("Destino é obrigatório")
		}
	case domain.InventoryTransactionTypeOut:
		if r.InventoryOriginId == 0 {
			return errors.New("Origem é obrigatório")
		}
	default:
		return errors.New("Tipo de transação inválida")
	}

	err := validator.Validate(r)
	if err != nil {
		return err
	}
	return nil
}
