package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
)

type SalesRepository interface {
	CreateSale(ctx context.Context, tx *sql.Tx, sale domain.Sales) (int64, error)
	CreateManySaleItem(ctx context.Context, tx *sql.Tx, sale domain.Sales, saleItems []domain.SalesItem) ([]int64, error)
	CreatePayment(ctx context.Context, tx *sql.Tx, sale domain.Sales, payment domain.SalesPayment) (int64, error)
	CreateManyPaymentDates(ctx context.Context, tx *sql.Tx, payment domain.SalesPayment, paymentDates []domain.SalesPaymentDates) ([]int64, error)
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

func (r *salesRepository) CreateManySaleItem(ctx context.Context, tx *sql.Tx, sale domain.Sales, saleItems []domain.SalesItem) ([]int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var ids []int64

	query := `INSERT INTO sales_items (quantity, unit_price, sku_id, sales_id, tenant_id) VALUES %s RETURNING id`
	valueStrings := make([]string, 0, len(saleItems))
	valueArgs := make([]interface{}, 0, len(saleItems)*5)

	for i, item := range saleItems {
		n := i * 5
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", n+1, n+2, n+3, n+4, n+5))
		valueArgs = append(valueArgs,
			item.Quantity,
			item.Sku.Price,
			item.Sku.Id,
			sale.Id,
			tenantId,
		)
	}
	query = fmt.Sprintf(query, strings.Join(valueStrings, ","))
	rows, err := tx.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
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

func (r *salesRepository) CreateManyPaymentDates(ctx context.Context, tx *sql.Tx, payment domain.SalesPayment, paymentDates []domain.SalesPaymentDates) ([]int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var ids []int64
	valueStrings := make([]string, 0, len(paymentDates))
	valueArgs := make([]interface{}, 0, len(paymentDates)*7)

	query := `INSERT INTO payment_dates (due_date, paid_date, installment_number, installment_value, status, payment_id, tenant_id) VALUES %s RETURNING id`

	for i, date := range paymentDates {
		n := i * 7
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)", n+1, n+2, n+3, n+4, n+5, n+6, n+7))
		valueArgs = append(valueArgs,
			date.DueDate,
			date.PaidDate,
			date.InstallmentNumber,
			date.InstallmentValue,
			date.Status,
			payment.Id,
			tenantId,
		)
	}
	query = fmt.Sprintf(query, strings.Join(valueStrings, ","))
	rows, err := tx.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
