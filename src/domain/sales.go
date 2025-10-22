package domain

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/bncunha/erp-api/src/infrastructure/ksuid"
)

type PaymentType string

type PaymentStatus string

var (
	ErrPaymentValueIsMissing       = errors.New("Valor de pagamento insuficiente para o valor da compra")
	ErrQuantityNotValid            = errors.New("Não há quantidade suficiente no estoque")
	ErrPaymentValueIsOverTotal     = errors.New("Valor de pagamento excede o valor total da compra")
	ErrPaymentDatesPast            = errors.New("As datas de pagamento devem ser maiores que a data atual")
	ErrPaymentDatesOrderInvalid    = errors.New("As datas de pagamento devem ser ordenadas")
	ErrPaymentDatesQuantityInvalid = errors.New("Pagamento em dinheiro, PIX ou débito deve ter apenas uma data de pagamento")
	ErrPaymentDatesCashAndPixRange = errors.New("Para pagamentos em dinheiro ou PIX, a data deve ser a partir de hoje e no máximo em 30 dias")
	ErrSkusDuplicated              = errors.New("SKUs duplicados encontrados")
	ErrPaymentTypesDuplicated      = errors.New("Tipos de pagamento duplicados encontrados. Remova os tipos de pagamento duplicados")
	ErrPaymentDatesDuplicated      = errors.New("As datas de pagamento devem ser diferentes")
)

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
	PaymentStatusDelayed PaymentStatus = "DELAYED"
)

type Sales struct {
	Id       int64
	Code     string
	Date     time.Time
	User     User
	Customer Customer
	Items    []SalesItem
	Payments []SalesPayment
}

func NewSales(date time.Time, user User, customer Customer, items []SalesItem, payments []SalesPayment) Sales {
	return Sales{
		Code:     "V-" + ksuid.New().String(),
		Date:     date,
		User:     user,
		Customer: customer,
		Items:    items,
		Payments: payments,
	}
}

func (s *Sales) ValidateSale() error {
	missingValue := s.getMissingValue()
	if missingValue > 0 {
		return errors.New(ErrPaymentValueIsMissing.Error() + fmt.Sprintf(": R$ %.2f", missingValue))
	} else if missingValue < 0 {
		return errors.New(ErrPaymentValueIsOverTotal.Error() + fmt.Sprintf(": R$ %.2f", missingValue))
	}
	if s.isPaymentTypesDuplicated() {
		return ErrPaymentTypesDuplicated
	}
	for _, item := range s.Items {
		if !item.isQuantityValid() {
			return errors.New(ErrQuantityNotValid.Error() + fmt.Sprintf(": (%d) %s", item.Sku.Id, item.Sku.GetName()))
		}
	}
	for _, payment := range s.Payments {
		if !payment.isPaymentDatesOrderValid() {
			return ErrPaymentDatesOrderInvalid
		}
		if payment.isPaymentDatesDuplicated() {
			return errors.New(ErrPaymentDatesDuplicated.Error() + fmt.Sprintf(": (%s)", payment.PaymentType))
		}
		if err := payment.validatePaymentDates(); err != nil {
			return err
		}
		if !payment.isPaymentDatesQuantityValid() {
			return ErrPaymentDatesQuantityInvalid
		}
	}
	return nil
}

func (s *Sales) getMissingValue() float64 {
	missingValue := s.GetTotal()
	for _, payment := range s.Payments {
		for _, date := range payment.Dates {
			missingValue -= date.InstallmentValue
		}
	}
	return math.Round(missingValue*100) / 100
}

func (s *Sales) GetTotal() float64 {
	var total float64
	for _, item := range s.Items {
		total += item.Sku.Price * item.Quantity
	}
	return total
}

func (s *Sales) isPaymentTypesDuplicated() bool {
	paymentTypes := make(map[PaymentType]bool)
	for _, payment := range s.Payments {
		if paymentTypes[payment.PaymentType] {
			return true
		}
		paymentTypes[payment.PaymentType] = true
	}
	return false
}

type SalesItem struct {
	Sku       Sku
	Quantity  float64
	UnitPrice float64
}

func NewSalesItem(sku Sku, quantity float64) SalesItem {
	return SalesItem{
		Sku:      sku,
		Quantity: quantity,
	}
}

func (s *SalesItem) isQuantityValid() bool {
	return s.Sku.Quantity-s.Quantity >= 0
}

type SalesPayment struct {
	Id          int64
	PaymentType PaymentType
	Dates       []SalesPaymentDates
}

func NewSalesPayment(paymentType PaymentType) SalesPayment {
	return SalesPayment{
		PaymentType: paymentType,
		Dates:       make([]SalesPaymentDates, 0),
	}
}

func (s *SalesPayment) shouldConfirmPayment() bool {
	return s.PaymentType == PaymentTypeCreditStore
}

func (s *SalesPayment) isOnCashPayment() bool {
	return s.PaymentType == PaymentTypeCash || s.PaymentType == PaymentTypePix || s.PaymentType == PaymentTypeDebitCard
}

func (s *SalesPayment) isPaymentDatesQuantityValid() bool {
	datesQuantity := len(s.Dates)
	if s.isOnCashPayment() && datesQuantity != 1 {
		return false
	}
	return true
}

func (s *SalesPayment) validatePaymentDates() error {
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	maxDate := todayStart.AddDate(0, 0, 30)
	for _, date := range s.Dates {
		dueDate := time.Date(date.DueDate.Year(), date.DueDate.Month(), date.DueDate.Day(), 0, 0, 0, 0, today.Location())
		if dueDate.Before(todayStart) {
			return ErrPaymentDatesPast
		}
		if s.PaymentType == PaymentTypeCash || s.PaymentType == PaymentTypePix {
			if dueDate.After(maxDate) {
				return ErrPaymentDatesCashAndPixRange
			}
		}
	}
	return nil
}

func isSameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func (s *SalesPayment) isPaymentDatesOrderValid() bool {
	sort.Slice(s.Dates, func(i, j int) bool {
		return s.Dates[i].DueDate.Before(s.Dates[j].DueDate)
	})

	for i := 1; i < len(s.Dates); i++ {
		if !s.Dates[i].DueDate.After(s.Dates[i-1].DueDate) {
			return false
		}
	}
	return true
}

func (s *SalesPayment) isPaymentDatesDuplicated() bool {
	paymentDates := make(map[SalesPaymentDates]bool)
	for _, date := range s.Dates {
		if paymentDates[date] {
			return true
		}
		paymentDates[date] = true
	}
	return false
}

func (s *SalesPayment) AppendNewSalesDate(dueDate time.Time, installmentNumber int, installmentValue float64) {
	newDate := NewSalesPaymentDates(dueDate, nil, installmentNumber, installmentValue, "")
	newDate.PaymentType = s.PaymentType
	if s.shouldConfirmPayment() {
		newDate.Status = PaymentStatusPending
	} else if s.PaymentType == PaymentTypeCash || s.PaymentType == PaymentTypePix {
		if isSameDay(dueDate, time.Now()) {
			paidDate := dueDate
			newDate.PaidDate = &paidDate
			newDate.Status = PaymentStatusPaid
		} else {
			newDate.Status = PaymentStatusPending
		}
	} else {
		paidDate := time.Now()
		newDate.DueDate = paidDate
		newDate.PaidDate = &paidDate
		newDate.Status = PaymentStatusPaid
	}
	s.Dates = append(s.Dates, newDate)

}

type SalesPaymentDates struct {
	Id                int64
	DueDate           time.Time
	PaidDate          *time.Time
	InstallmentNumber int
	InstallmentValue  float64
	Status            PaymentStatus
	PaymentType       PaymentType
}

func NewSalesPaymentDates(dueDate time.Time, paidDate *time.Time, installmentNumber int, installmentValue float64, status PaymentStatus) SalesPaymentDates {
	return SalesPaymentDates{
		DueDate:           dueDate,
		PaidDate:          paidDate,
		InstallmentNumber: installmentNumber,
		InstallmentValue:  installmentValue,
		Status:            status,
	}
}
