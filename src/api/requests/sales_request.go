package request

import (
	"time"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/validator"
	"github.com/bncunha/erp-api/src/domain"
)

type ListSalesRequest struct {
	CustomerId    []int64               `json:"customer_id"`
	MinDate       *time.Time            `json:"min_date"`
	MaxDate       *time.Time            `json:"max_date"`
	UserId        []int64               `json:"user_id"`
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
	PaymentType          domain.PaymentType `json:"payment_type" validate:"required,oneof=CASH CREDIT_CARD DEBIT_CARD PIX CREDIT_STORE"`
	Value                float64            `json:"value" validate:"required,gt=0"`
	InstallmentsQuantity *int               `json:"installments_quantity" validate:"omitempty,gt=1"`
	FirstInstallmentDate *time.Time         `json:"first_installment_date"`
}

func (p *CreateSaleRequestPayments) Validate() error {
	err := validator.Validate(p)
	if (p.PaymentType == domain.PaymentTypeCreditStore || p.PaymentType == domain.PaymentTypeCreditCard) && p.InstallmentsQuantity == nil {
		return errors.New("Quantidade de parcelas é obrigatória para Cartão de Crédito ou Notinha")
	}
	if p.PaymentType == domain.PaymentTypeCreditStore && p.FirstInstallmentDate == nil {
		return errors.New("Data de primeira parcela é obrigatória para Notinha")
	}
	if err != nil {
		return err
	}
	return nil
}

type ChangePaymentStatusRequest struct {
	Status string    `json:"status" validate:"required,oneof=PAID CANCEL PENDING"`
	Date   time.Time `json:"date"`
}

func (c *ChangePaymentStatusRequest) Validate() error {
	err := validator.Validate(c)
	if err != nil {
		return err
	}

	if c.Status == string(domain.PaymentStatusPaid) && c.Date.IsZero() {
		return errors.New("Data de pagamento é obrigatória")
	}
	return nil
}
