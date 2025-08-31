package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
)

type SalesRepository interface {
	CreateSale(ctx context.Context, tx *sql.Tx, sale domain.Sales) (int64, error)
	CreateManySaleItem(ctx context.Context, tx *sql.Tx, sale domain.Sales, saleItem []domain.SalesItem) (int64, error)
	CreatePayment(ctx context.Context, tx *sql.Tx, sale domain.Sales, payment domain.SalesPayment) (int64, error)
	CreateManyPaymentDates(ctx context.Context, tx *sql.Tx, payment domain.SalesPayment, paymentDates []domain.SalesPaymentDates) (int64, error)
}

type salesRepository struct {
}

func NewSalesRepository() SalesRepository {
	return &salesRepository{}
}

func (r *salesRepository) CreateSale(ctx context.Context, tx *sql.Tx, sale domain.Sales) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedId int64
	query := `INSERT INTO sales (date, user_id, customer_id, tenant_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := tx.QueryRowContext(ctx, query, sale.Date, sale.User.Id, sale.Customer.Id, tenantId).Scan(&insertedId)
	if err != nil {
		return insertedId, err
	}
	return insertedId, nil
}

func (r *salesRepository) CreateManySaleItem(ctx context.Context, tx *sql.Tx, sale domain.Sales, saleItem []domain.SalesItem) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedId int64
	query := `INSERT INTO sales_items (quantity, unit_price, sku_id, sales_id, tenant_id) VALUES %s RETURNING id`
	values := ""
	for _, item := range saleItem {
		values += fmt.Sprintf("(%f, %f, %d, %d, %d)", item.Quantity, item.UnitPrice, item.Sku.Id, sale.Id, tenantId)
	}
	query = fmt.Sprintf("%s %s", query, values)
	err := tx.QueryRowContext(ctx, query).Scan(&insertedId)
	if err != nil {
		return insertedId, err
	}
	return insertedId, nil
}

func (r *salesRepository) CreatePayment(ctx context.Context, tx *sql.Tx, sale domain.Sales, payment domain.SalesPayment) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedId int64
	query := `INSERT INTO payments (payment_type, sales_id, tenant_id) VALUES ($1, $2, $3) RETURNING id`
	err := tx.QueryRowContext(ctx, query, payment.PaymentType, sale.Id, tenantId).Scan(&insertedId)
	if err != nil {
		return insertedId, err
	}
	return insertedId, nil
}

func (r *salesRepository) CreateManyPaymentDates(ctx context.Context, tx *sql.Tx, payment domain.SalesPayment, paymentDates []domain.SalesPaymentDates) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedId int64
	query := `INSERT INTO payment_dates (due_date, paid_date, installment_number, installment_value, status, payment_id, tenant_id) VALUES %s RETURNING id`
	values := ""
	for _, date := range paymentDates {
		values += fmt.Sprintf("(%s, %s, %d, %f, %s, %d, %d)", date.DueDate, date.PaidDate, date.InstallmentNumber, date.InstallmentValue, date.Status, payment.Id, tenantId)
	}
	query = fmt.Sprintf("%s %s", query, values)
	err := tx.QueryRowContext(ctx, query).Scan(&insertedId)
	if err != nil {
		return insertedId, err
	}
	return insertedId, nil
}
