package domain

import "time"

type PaymentType string

type PaymentStatus string

const (
	PaymentTypeCash        PaymentType = "CASH"
	PaymentTypeCreditCard  PaymentType = "CREDIT_CARD"
	PaymentTypeDebitCard   PaymentType = "DEBIT_CARD"
	PaymentTypePix         PaymentType = "PIX"
	PaymentTypeCreditStore PaymentType = "CREDIT_STORE"
)

const (
	PaymentStatusPaid    PaymentStatus = "PAID"
	PaymentStatusCancel  PaymentStatus = "CANCEL"
	PaymentStatusPending PaymentStatus = "PENDING"
)

type Sales struct {
	Id       int64
	Date     time.Time
	User     User
	Customer Customer
	Items    []SalesItem
	Payments []SalesPayment
}

func NewSales(date time.Time, user User, customer Customer, items []SalesItem, payments []SalesPayment) Sales {
	return Sales{
		Date:     date,
		User:     user,
		Customer: customer,
		Items:    items,
		Payments: payments,
	}
}

func (s *Sales) IsPaymentValuesMatchTotalValue() bool {
	totalPrice := s.GetTotal()
	for _, payment := range s.Payments {
		for _, date := range payment.Dates {
			totalPrice -= date.InstallmentValue * float64(date.InstallmentNumber)
		}
	}
	return totalPrice == 0
}

func (s *Sales) GetTotal() float64 {
	var total float64
	for _, item := range s.Items {
		total += item.UnitPrice * item.Quantity
	}
	return total
}

type SalesItem struct {
	Sku       Sku
	UnitPrice float64
	Quantity  float64
}

func NewSalesItem(sku Sku, unitPrice float64, quantity float64) SalesItem {
	return SalesItem{
		Sku:       sku,
		UnitPrice: unitPrice,
		Quantity:  quantity,
	}
}

func (s *SalesItem) IsQuantityValid() bool {
	return s.Quantity >= s.Sku.Quantity-s.Quantity
}

type SalesPayment struct {
	Id          int64
	PaymentType PaymentType
	Dates       []SalesPaymentDates
}

func NewSalesPayment(paymentType PaymentType, dates []SalesPaymentDates) SalesPayment {
	return SalesPayment{
		PaymentType: paymentType,
		Dates:       dates,
	}
}

type SalesPaymentDates struct {
	DueDate           time.Time
	PaidDate          time.Time
	InstallmentNumber int
	InstallmentValue  float64
	Status            PaymentStatus
}

func NewSalesPaymentDates(dueDate time.Time, paidDate time.Time, installmentNumber int, installmentValue float64, status PaymentStatus) SalesPaymentDates {
	return SalesPaymentDates{
		DueDate:           dueDate,
		PaidDate:          paidDate,
		InstallmentNumber: installmentNumber,
		InstallmentValue:  installmentValue,
		Status:            status,
	}
}
