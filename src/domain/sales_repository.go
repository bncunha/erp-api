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

type SaleWithVersionOutput struct {
	Id            int64
	Code          string
	Date          time.Time
	UserId        int64
	CustomerId    int64
	CustomerName  string
	LastVersion   int
	SalesVersionId int64
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

type GetSalesReturnOutput struct {
	Id         int64
	ReturnDate time.Time
	Returner   string
	Reason     string
	Items      []GetSalesReturnItemOutput
}

type GetSalesReturnItemOutput struct {
	Sku       Sku
	Quantity  float64
	UnitPrice float64
}

type SalesRepository interface {
	CreateSale(ctx context.Context, tx *sql.Tx, sale Sales) (int64, error)
	CreateSaleVersion(ctx context.Context, tx *sql.Tx, saleId int64, version int, date time.Time) (int64, error)
	CreateManySaleItem(ctx context.Context, tx *sql.Tx, sale Sales, saleItems []SalesItem) ([]int64, error)
	CreatePayment(ctx context.Context, tx *sql.Tx, sale Sales, payment SalesPayment) (int64, error)
	CreateManyPaymentDates(ctx context.Context, tx *sql.Tx, payment SalesPayment, paymentDates []SalesPaymentDates) ([]int64, error)
	CreateSalesReturn(ctx context.Context, tx *sql.Tx, saleId int64, fromSalesVersionId int64, toSalesVersionId int64, salesReturn SalesReturn, createdByUserId int64) (int64, error)
	CreateSalesReturnItems(ctx context.Context, tx *sql.Tx, salesReturnId int64, items []SalesReturnItem) ([]int64, error)
	UpdateSaleLastVersion(ctx context.Context, tx *sql.Tx, saleId int64, version int) error
	CancelPaymentDatesBySaleVersionId(ctx context.Context, tx *sql.Tx, saleVersionId int64) error
	GetSaleByIdForUpdate(ctx context.Context, tx *sql.Tx, id int64) (SaleWithVersionOutput, error)
	GetSaleVersionIdBySaleIdAndVersion(ctx context.Context, tx *sql.Tx, saleId int64, version int) (int64, error)
	GetPaymentsBySaleVersionId(ctx context.Context, saleVersionId int64) ([]GetSalesPaymentOutput, error)
	GetItemsBySaleVersionId(ctx context.Context, saleVersionId int64) ([]GetItemsOutput, error)
	GetReturnsBySaleId(ctx context.Context, id int64) ([]GetSalesReturnOutput, error)
	GetSales(ctx context.Context, input GetSalesInput) ([]GetSalesItemOutput, error)
	GetSaleById(ctx context.Context, id int64) (GetSaleByIdOutput, error)
	GetPaymentsBySaleId(ctx context.Context, id int64) ([]GetSalesPaymentOutput, error)
	GetItemsBySaleId(ctx context.Context, id int64) ([]GetItemsOutput, error)
	ChangePaymentStatus(ctx context.Context, id int64, status PaymentStatus) (int64, error)
	ChangePaymentDate(ctx context.Context, id int64, date *time.Time) (int64, error)
	GetPaymentDatesBySaleIdAndPaymentDateId(ctx context.Context, id int64, paymentDateId int64) (SalesPaymentDates, error)
}
