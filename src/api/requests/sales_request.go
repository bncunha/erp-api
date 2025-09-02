package request

import (
	"time"

	"github.com/bncunha/erp-api/src/application/validator"
	"github.com/bncunha/erp-api/src/domain"
)

type CreateSaleRequest struct {
	CustomerId int64                       `json:"customer_id" validate:"required"`
	Items      []CreateSaleRequestItems    `json:"items" validate:"required"`
	Payments   []CreateSaleRequestPayments `json:"payments" validate:"required"`
}

func (r *CreateSaleRequest) Validate() error {
	err := validator.Validate(r)
	if err != nil {
		return err
	}
	return nil
}

type CreateSaleRequestItems struct {
	SkuId    int64   `json:"sku_id" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required,gt=0"`
}

type CreateSaleRequestPayments struct {
	PaymentType domain.PaymentType              `json:"payment_type" validate:"required,oneof=CASH CREDIT_CARD DEBIT_CARD PIX CREDIT_STORE"`
	Dates       []CreateSaleRequestPaymentDates `json:"dates" validate:"required"`
}

type CreateSaleRequestPaymentDates struct {
	Date             time.Time `json:"date" validate:"required"`
	InstallmentValue float64   `json:"installment_value" validate:"required,gt=0"`
}
