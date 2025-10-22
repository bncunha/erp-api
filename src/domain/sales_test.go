package domain

import (
	"strings"
	"testing"
	"time"
)

func TestNewSalesGeneratesCodeAndCopiesFields(t *testing.T) {
	now := time.Now()
	user := User{Id: 1}
	customer := Customer{Id: 2}
	sale := NewSales(now, user, customer, nil, nil)

	if !strings.HasPrefix(sale.Code, "V-") {
		t.Fatalf("expected code to have ksuid prefix, got %s", sale.Code)
	}
	if !sale.Date.Equal(now) || sale.User.Id != user.Id || sale.Customer.Id != customer.Id {
		t.Fatalf("unexpected sale data: %+v", sale)
	}
}

func TestSalesValidateSaleSuccess(t *testing.T) {
	sku := Sku{Id: 1, Price: 10, Quantity: 5}
	item := SalesItem{Sku: sku, Quantity: 2}
	due := time.Now().Add(48 * time.Hour)
	payment := SalesPayment{PaymentType: PaymentTypeCash, Dates: []SalesPaymentDates{{
		DueDate:           due,
		InstallmentNumber: 1,
		InstallmentValue:  20,
		Status:            PaymentStatusPaid,
	}}}
	sale := Sales{Items: []SalesItem{item}, Payments: []SalesPayment{payment}}

	if err := sale.ValidateSale(); err != nil {
		t.Fatalf("expected sale to be valid, got %v", err)
	}
}

func TestSalesValidateSaleMissingValue(t *testing.T) {
	sku := Sku{Price: 10, Quantity: 5}
	item := SalesItem{Sku: sku, Quantity: 2}
	due := time.Now().Add(48 * time.Hour)
	payment := SalesPayment{PaymentType: PaymentTypeCash, Dates: []SalesPaymentDates{{
		DueDate:           due,
		InstallmentNumber: 1,
		InstallmentValue:  10,
		Status:            PaymentStatusPaid,
	}}}
	sale := Sales{Items: []SalesItem{item}, Payments: []SalesPayment{payment}}

	if err := sale.ValidateSale(); err == nil || !strings.Contains(err.Error(), ErrPaymentValueIsMissing.Error()) {
		t.Fatalf("expected missing value error, got %v", err)
	}
}

func TestSalesValidateSaleOverPayment(t *testing.T) {
	sku := Sku{Price: 10, Quantity: 5}
	item := SalesItem{Sku: sku, Quantity: 2}
	due := time.Now().Add(48 * time.Hour)
	payment := SalesPayment{PaymentType: PaymentTypeCash, Dates: []SalesPaymentDates{{
		DueDate:           due,
		InstallmentNumber: 1,
		InstallmentValue:  30,
		Status:            PaymentStatusPaid,
	}}}
	sale := Sales{Items: []SalesItem{item}, Payments: []SalesPayment{payment}}

	if err := sale.ValidateSale(); err == nil || !strings.Contains(err.Error(), ErrPaymentValueIsOverTotal.Error()) {
		t.Fatalf("expected over payment error, got %v", err)
	}
}

func TestSalesValidateSaleDuplicatePaymentTypes(t *testing.T) {
	due := time.Now().Add(48 * time.Hour)
	payment := SalesPayment{PaymentType: PaymentTypeCash, Dates: []SalesPaymentDates{{
		DueDate:           due,
		InstallmentNumber: 1,
		InstallmentValue:  2.5,
		Status:            PaymentStatusPaid,
	}}}
	sale := Sales{Items: []SalesItem{{Sku: Sku{Price: 5, Quantity: 1}, Quantity: 1}}, Payments: []SalesPayment{payment, payment}}

	if err := sale.ValidateSale(); err != ErrPaymentTypesDuplicated {
		t.Fatalf("expected duplicate payment types error, got %v", err)
	}
}

func TestSalesValidateSaleQuantityInvalid(t *testing.T) {
	payment := SalesPayment{PaymentType: PaymentTypeCash, Dates: []SalesPaymentDates{{
		DueDate:           time.Now().Add(48 * time.Hour),
		InstallmentNumber: 1,
		InstallmentValue:  20,
		Status:            PaymentStatusPaid,
	}}}
	sale := Sales{Items: []SalesItem{{Sku: Sku{Id: 1, Price: 10, Quantity: 1}, Quantity: 2}}, Payments: []SalesPayment{payment}}

	if err := sale.ValidateSale(); err == nil || !strings.Contains(err.Error(), ErrQuantityNotValid.Error()) {
		t.Fatalf("expected quantity invalid error, got %v", err)
	}
}

func TestSalesValidateSalePaymentOrderInvalid(t *testing.T) {
	due := time.Now().Add(48 * time.Hour)
	payment := SalesPayment{PaymentType: PaymentTypeCreditStore, Dates: []SalesPaymentDates{
		{DueDate: due, InstallmentNumber: 2, InstallmentValue: 10},
		{DueDate: due, InstallmentNumber: 1, InstallmentValue: 10},
	}}
	sale := Sales{Items: []SalesItem{{Sku: Sku{Price: 10, Quantity: 5}, Quantity: 2}}, Payments: []SalesPayment{payment}}

	if err := sale.ValidateSale(); err != ErrPaymentDatesOrderInvalid {
		t.Fatalf("expected payment date order error, got %v", err)
	}
}

func TestSalesPaymentDatesQuantityValid(t *testing.T) {
	payment := SalesPayment{PaymentType: PaymentTypeCash, Dates: []SalesPaymentDates{{}, {}}}
	if payment.isPaymentDatesQuantityValid() {
		t.Fatalf("expected cash payments with more than one date to be invalid")
	}
}

func TestSalesPaymentValidateDatesPast(t *testing.T) {
	payment := SalesPayment{PaymentType: PaymentTypeCreditStore, Dates: []SalesPaymentDates{{DueDate: time.Now().Add(-48 * time.Hour)}}}
	if err := payment.validatePaymentDates(); err != ErrPaymentDatesPast {
		t.Fatalf("expected past due date error, got %v", err)
	}
}

func TestSalesPaymentValidateDatesCashPixRange(t *testing.T) {
	payment := SalesPayment{PaymentType: PaymentTypeCash, Dates: []SalesPaymentDates{{DueDate: time.Now().AddDate(0, 0, 31)}}}
	if err := payment.validatePaymentDates(); err != ErrPaymentDatesCashAndPixRange {
		t.Fatalf("expected cash/pix range error, got %v", err)
	}
}

func TestSalesPaymentDatesDuplicated(t *testing.T) {
	due := time.Now().Add(48 * time.Hour)
	duplicated := SalesPaymentDates{DueDate: due, InstallmentNumber: 1}
	payment := SalesPayment{PaymentType: PaymentTypeCreditStore, Dates: []SalesPaymentDates{duplicated, duplicated}}
	if !payment.isPaymentDatesDuplicated() {
		t.Fatalf("expected duplicated payment dates")
	}
}

func TestSalesPaymentAppendNewSalesDate(t *testing.T) {
	payment := NewSalesPayment(PaymentTypeCreditStore)
	due := time.Now().Add(48 * time.Hour)
	payment.AppendNewSalesDate(due, 1, 10)

	if len(payment.Dates) != 1 {
		t.Fatalf("expected one payment date")
	}
	if payment.Dates[0].Status != PaymentStatusPending {
		t.Fatalf("expected pending status for credit store: %+v", payment.Dates[0])
	}
	if payment.Dates[0].PaidDate != nil {
		t.Fatalf("expected nil paid date for credit store")
	}
}

func TestSalesPaymentAppendNewSalesDateCashFuture(t *testing.T) {
	payment := NewSalesPayment(PaymentTypeCash)
	due := time.Now().Add(48 * time.Hour)
	payment.AppendNewSalesDate(due, 1, 10)

	if len(payment.Dates) != 1 {
		t.Fatalf("expected one payment date")
	}
	if payment.Dates[0].Status != PaymentStatusPending || payment.Dates[0].PaidDate != nil {
		t.Fatalf("expected pending status without paid date for future cash payment: %+v", payment.Dates[0])
	}
}

func TestSalesPaymentAppendNewSalesDateCashSameDay(t *testing.T) {
	payment := NewSalesPayment(PaymentTypeCash)
	due := time.Now()
	payment.AppendNewSalesDate(due, 1, 10)

	if len(payment.Dates) != 1 {
		t.Fatalf("expected one payment date")
	}
	if payment.Dates[0].Status != PaymentStatusPaid || payment.Dates[0].PaidDate == nil {
		t.Fatalf("expected paid status for cash on same day: %+v", payment.Dates[0])
	}
	if !payment.Dates[0].DueDate.Equal(*payment.Dates[0].PaidDate) {
		t.Fatalf("expected due date to match paid date")
	}
}

func TestSalesHelpers(t *testing.T) {
	sale := Sales{}
	if sale.GetTotal() != 0 {
		t.Fatalf("expected total zero")
	}

	sku := Sku{Price: 15, Quantity: 5}
	sale.Items = []SalesItem{{Sku: sku, Quantity: 2}}
	if sale.GetTotal() != 30 {
		t.Fatalf("expected total 30, got %v", sale.GetTotal())
	}

	sale.Payments = []SalesPayment{{Dates: []SalesPaymentDates{{InstallmentValue: 20}}}}
	missing := sale.getMissingValue()
	if missing != 10 {
		t.Fatalf("expected missing value 10, got %v", missing)
	}

	if sale.isPaymentTypesDuplicated() {
		t.Fatalf("expected no duplicated payment types")
	}
	sale.Payments = append(sale.Payments, SalesPayment{PaymentType: PaymentTypeCash})
	sale.Payments[0].PaymentType = PaymentTypeCash
	if !sale.isPaymentTypesDuplicated() {
		t.Fatalf("expected duplicated payment types")
	}
}

func TestSalesItemHelpers(t *testing.T) {
	sku := Sku{Quantity: 5}
	item := NewSalesItem(sku, 3)
	if !item.isQuantityValid() {
		t.Fatalf("expected quantity to be valid")
	}
	item.Quantity = 6
	if item.isQuantityValid() {
		t.Fatalf("expected quantity to be invalid")
	}
}

func TestSalesPaymentHelpers(t *testing.T) {
	payment := NewSalesPayment(PaymentTypeCash)
	if !payment.isOnCashPayment() {
		t.Fatalf("expected cash payment to be recognized")
	}
	if payment.shouldConfirmPayment() {
		t.Fatalf("expected cash payment to not require confirmation")
	}
	credit := NewSalesPayment(PaymentTypeCreditStore)
	if !credit.shouldConfirmPayment() {
		t.Fatalf("expected credit store to require confirmation")
	}
}

func TestSalesPaymentDatesConstructor(t *testing.T) {
	due := time.Now()
	paid := time.Now()
	date := NewSalesPaymentDates(due, &paid, 1, 5, PaymentStatusPaid)
	if !date.DueDate.Equal(due) || date.PaidDate != &paid || date.InstallmentNumber != 1 || date.InstallmentValue != 5 || date.Status != PaymentStatusPaid {
		t.Fatalf("unexpected payment date struct: %+v", date)
	}
}

func TestSkuGetNameWithCost(t *testing.T) {
	cost := 10.0
	sku := Sku{Code: "S1", Color: "Red", Size: "M", Cost: &cost, Product: Product{Name: "Shirt"}}
	if got := sku.GetName(); got != "Shirt - Red - M" {
		t.Fatalf("unexpected sku name: %s", got)
	}
}

func TestInventoryHelpers(t *testing.T) {
	sku := Sku{Id: 1}
	item := NewInventoryItem(2, sku, 3)
	if item.InventoryId != 2 || item.Sku.Id != sku.Id || item.Quantity != 3 {
		t.Fatalf("unexpected inventory item: %+v", item)
	}
}
