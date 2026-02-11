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
	SkuId    int64
	Quantity float64
}

type DoSalePaymentsInput struct {
	PaymentType domain.PaymentType
	Value       float64
	Dates       []DoSalePaymentDatesInput
}

type DoSalePaymentDatesInput struct {
	DueDate           time.Time
	PaidDate          time.Time
	InstallmentNumber int
	InstallmentValue  float64
	Status            domain.PaymentStatus
	DateInformed      bool
}

type DoReturnInput struct {
	SaleId                int64
	UserId                int64
	InventoryDestinationId int64
	ReturnerName          string
	Reason                string
	Items                 []DoReturnItemInput
}

type DoReturnItemInput struct {
	SkuId    int64
	Quantity float64
}
