package request

import (
	"time"

	"github.com/bncunha/erp-api/src/application/validator"
	"github.com/bncunha/erp-api/src/domain"
)

type ListSalesRequest struct {
	CustomerId    *int64                `json:"customer_id"`
	MinDate       *time.Time            `json:"min_date"`
	MaxDate       *time.Time            `json:"max_date"`
	UserId        *int64                `json:"user_id"`
	PaymentStatus *domain.PaymentStatus `json:"payment_status"`
}

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
	for _, item := range r.Items {
		err = item.Validate()
		if err != nil {
			return err
		}
	}
	for _, payment := range r.Payments {
		err = payment.Validate()
		if err != nil {
			return err
		}
		for _, date := range payment.Dates {
			err = date.Validate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type CreateSaleRequestItems struct {
	SkuId    int64   `json:"sku_id" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required,gt=0"`
}

func (i *CreateSaleRequestItems) Validate() error {
	err := validator.Validate(i)
	if err != nil {
		return err
	}
	return nil
}

type CreateSaleRequestPayments struct {
	PaymentType domain.PaymentType              `json:"payment_type" validate:"required,oneof=CASH CREDIT_CARD DEBIT_CARD PIX CREDIT_STORE"`
	Dates       []CreateSaleRequestPaymentDates `json:"dates" validate:"required"`
}

func (p *CreateSaleRequestPayments) Validate() error {
	err := validator.Validate(p)
	if err != nil {
		return err
	}
	return nil
}

type CreateSaleRequestPaymentDates struct {
	Date             time.Time `json:"date" validate:"required"`
	InstallmentValue float64   `json:"installment_value" validate:"required,gt=0"`
}

func (d *CreateSaleRequestPaymentDates) Validate() error {
	err := validator.Validate(d)
	if err != nil {
		return err
	}
	return nil
}
