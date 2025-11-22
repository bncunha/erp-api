package domain

import (
	"context"
	"database/sql"
	"time"
)

type GetSalesInput struct {
	InitialDate   *time.Time
	FinalDate     *time.Time
	UserId        []int64
	CustomerId    []int64
	PaymentStatus *PaymentStatus
}

type GetSalesItemOutput struct {
	Id            int
	Date          string
	SellerName    string
	CustomerName  string
	TotalValue    float64
	ReceivedValue float64
	FutureRevenue float64
	TotalItems    float64
	Status        PaymentStatus
}

type GetSaleByIdOutput struct {
	Id            int
	Code          string
	Date          time.Time
	TotalValue    float64
	SellerName    string
	CustomerName  string
	ReceivedValue float64
	FutureRevenue float64
	PaymentStatus PaymentStatus
}

type GetSalesPaymentOutput struct {
	Id                int64
	InstallmentNumber int64
	InstallmentValue  float64
	DueDate           time.Time
	PaidDate          *time.Time
	PaymentStatus     PaymentStatus
	PaymentType       PaymentType
}

type GetItemsOutput struct {
	Sku        Sku
	Quantity   float64
	UnitPrice  float64
	TotalValue float64
}

type SalesRepository interface {
	CreateSale(ctx context.Context, tx *sql.Tx, sale Sales) (int64, error)
	CreateManySaleItem(ctx context.Context, tx *sql.Tx, sale Sales, saleItems []SalesItem) ([]int64, error)
	CreatePayment(ctx context.Context, tx *sql.Tx, sale Sales, payment SalesPayment) (int64, error)
	CreateManyPaymentDates(ctx context.Context, tx *sql.Tx, payment SalesPayment, paymentDates []SalesPaymentDates) ([]int64, error)
	GetSales(ctx context.Context, input GetSalesInput) ([]GetSalesItemOutput, error)
	GetSaleById(ctx context.Context, id int64) (GetSaleByIdOutput, error)
	GetPaymentsBySaleId(ctx context.Context, id int64) ([]GetSalesPaymentOutput, error)
	GetItemsBySaleId(ctx context.Context, id int64) ([]GetItemsOutput, error)
	ChangePaymentStatus(ctx context.Context, id int64, status PaymentStatus) (int64, error)
	ChangePaymentDate(ctx context.Context, id int64, date *time.Time) (int64, error)
	GetPaymentDatesBySaleIdAndPaymentDateId(ctx context.Context, id int64, paymentDateId int64) (SalesPaymentDates, error)
}
