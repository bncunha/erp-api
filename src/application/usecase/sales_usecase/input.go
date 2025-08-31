package sales_usecase

import (
	"time"

	"github.com/bncunha/erp-api/src/domain"
)

type DoSaleInput struct {
	Date       time.Time
	UserId     int64
	CustomerId int64
	Payments   []DoSalePaymentsInput
	Items      []DoSaleItemsInput
}

type DoSaleItemsInput struct {
	SkuId     int64
	UnitPrice float64
	Quantity  float64
}

type DoSalePaymentsInput struct {
	PaymentType domain.PaymentType
	Dates       []DoSalePaymentDatesInput
}

type DoSalePaymentDatesInput struct {
	DueDate           time.Time
	PaidDate          time.Time
	InstallmentNumber int
	InstallmentValue  float64
	Status            string
}
