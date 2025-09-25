package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/service/input"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/lib/pq"
)

type SalesRepository interface {
	CreateSale(ctx context.Context, tx *sql.Tx, sale domain.Sales) (int64, error)
	CreateManySaleItem(ctx context.Context, tx *sql.Tx, sale domain.Sales, saleItems []domain.SalesItem) ([]int64, error)
	CreatePayment(ctx context.Context, tx *sql.Tx, sale domain.Sales, payment domain.SalesPayment) (int64, error)
	CreateManyPaymentDates(ctx context.Context, tx *sql.Tx, payment domain.SalesPayment, paymentDates []domain.SalesPaymentDates) ([]int64, error)
	GetSales(ctx context.Context, input input.GetSalesInput) ([]output.GetSalesItemOutput, error)
	GetSaleById(ctx context.Context, id int64) (output.GetSaleByIdOutput, error)
	GetPaymentsBySaleId(ctx context.Context, id int64) ([]output.GetSalesPaymentOutput, error)
	GetItemsBySaleId(ctx context.Context, id int64) ([]output.GetItemsOutput, error)
	ChangePaymentStatus(ctx context.Context, id int64, status domain.PaymentStatus) (int64, error)
	ChangePaymentDate(ctx context.Context, id int64, date *time.Time) (int64, error)
	GetPaymentDatesBySaleIdAndPaymentDateId(ctx context.Context, id int64, paymentDateId int64) (domain.SalesPaymentDates, error)
}

type salesRepository struct {
	db *sql.DB
}

func NewSalesRepository(db *sql.DB) SalesRepository {
	return &salesRepository{db}
}

func (r *salesRepository) CreateSale(ctx context.Context, tx *sql.Tx, sale domain.Sales) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedId int64
	query := `INSERT INTO sales (date, user_id, customer_id, tenant_id, code) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := tx.QueryRowContext(ctx, query, sale.Date, sale.User.Id, sale.Customer.Id, tenantId, sale.Code).Scan(&insertedId)
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

func (r *salesRepository) GetSales(ctx context.Context, input input.GetSalesInput) ([]output.GetSalesItemOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var sales []output.GetSalesItemOutput

	query := `
SELECT
  s.id,
  s.date,
  c.name AS customer,
  u.name AS seller,

  COALESCE(SUM(pd.installment_value), 0) AS total_value,

  (
    SELECT COALESCE(SUM(si.quantity), 0)
    FROM sales_items si
    WHERE si.sales_id = s.id
      AND si.tenant_id = s.tenant_id
  ) AS total_items,

  CASE
    WHEN COUNT(*) > 0
         AND COUNT(*) = COUNT(*) FILTER (WHERE pd.status = 'PAID')
      THEN 'PAID'
    WHEN BOOL_OR(pd.status = 'PENDING' AND pd.due_date < CURRENT_DATE)
      THEN 'DELAYED'
    ELSE 'IN_DAY'
  END AS summary_status

FROM sales s
JOIN users u
  ON u.id = s.user_id AND u.tenant_id = s.tenant_id
JOIN customers c
  ON c.id = s.customer_id AND c.tenant_id = s.tenant_id
LEFT JOIN payments p
  ON p.sales_id  = s.id AND p.tenant_id = s.tenant_id
LEFT JOIN payment_dates pd
  ON pd.payment_id = p.id AND pd.tenant_id = s.tenant_id

WHERE s.tenant_id = $1
  AND ($2::bigint[]      IS NULL OR s.user_id     = ANY($2))
  AND ($3::bigint[]      IS NULL OR s.customer_id = ANY($3))
  AND ($4::timestamptz IS NULL OR s.date >= $4)
  AND ($5::timestamptz IS NULL OR s.date <= $5)

GROUP BY s.id, s.date, c.name, u.name

HAVING
  $6::text IS NULL
  OR (
       -- $6 = 'PAID'
       ($6 = 'PAID' AND COUNT(*) > 0
                    AND COUNT(*) = COUNT(*) FILTER (WHERE pd.status = 'PAID'))
       -- $6 = 'DELAYED'
    OR ($6 = 'DELAYED' AND BOOL_OR(pd.status = 'PENDING' AND pd.due_date < CURRENT_DATE))
       -- $6 = 'IN_DAY'
    OR ($6 = 'IN_DAY' AND NOT (
          (COUNT(*) > 0 AND COUNT(*) = COUNT(*) FILTER (WHERE pd.status = 'PAID'))
          OR BOOL_OR(pd.status = 'PENDING' AND pd.due_date < CURRENT_DATE)
       ))
  )

ORDER BY s.date DESC, s.id DESC;
	`
	valueArgs := []interface{}{tenantId, pq.Array(input.UserId), pq.Array(input.CustomerId), input.InitialDate, input.FinalDate, input.PaymentStatus}

	rows, err := r.db.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sale output.GetSalesItemOutput
		if err := rows.Scan(&sale.Id, &sale.Date, &sale.CustomerName, &sale.SellerName, &sale.TotalValue, &sale.TotalItems, &sale.Status); err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}
	return sales, nil
}

func (r *salesRepository) GetSaleById(ctx context.Context, id int64) (output.GetSaleByIdOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var output output.GetSaleByIdOutput

	query := `
	SELECT
  s.id,
  s.code,
  s.date,
  u.name AS seller_name,
  c.name AS customer_name,

  /* totais */
  COALESCE(SUM(pd.installment_value), 0) AS total_value,
  COALESCE(SUM(pd.installment_value) FILTER (WHERE pd.status = 'PAID'), 0) AS received_value,
  COALESCE(SUM(pd.installment_value) FILTER (WHERE pd.status IN ('PENDING','DELAYED')), 0) AS future_revenue,

  /* status agregado da venda */
  CASE
    WHEN COUNT(pd.id) > 0
         AND COUNT(pd.id) = COUNT(pd.id) FILTER (WHERE pd.status = 'PAID')
      THEN 'PAID'
    WHEN COALESCE(BOOL_OR(pd.status = 'PENDING' AND pd.due_date < CURRENT_DATE), FALSE)
      THEN 'DELAYED'
    ELSE 'IN_DAY'
  END AS payment_status

FROM sales s
JOIN users u
  ON u.id = s.user_id
 AND u.tenant_id = s.tenant_id
JOIN customers c
  ON c.id = s.customer_id
 AND c.tenant_id = s.tenant_id
LEFT JOIN payments p
  ON p.sales_id = s.id
 AND p.tenant_id = s.tenant_id
LEFT JOIN payment_dates pd
  ON pd.payment_id = p.id
 AND pd.tenant_id = s.tenant_id

WHERE s.id = $1
  AND s.tenant_id = $2

GROUP BY
  s.id, s.code, s.date, u.name, c.name;

	`
	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&output.Id, &output.Code, &output.Date, &output.SellerName, &output.CustomerName, &output.TotalValue, &output.ReceivedValue, &output.FutureRevenue, &output.PaymentStatus)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (r *salesRepository) GetPaymentsBySaleId(ctx context.Context, id int64) ([]output.GetSalesPaymentOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var o []output.GetSalesPaymentOutput

	query := `
	SELECT 
		pd.id,
		pd.installment_number,
		pd.installment_value,
		pd.due_date,
		pd.paid_date,
		p.payment_type,
		pd.status
	FROM payments p
	JOIN payment_dates pd ON p.id = pd.payment_id
	WHERE p.sales_id = $1 AND p.tenant_id = $2
	ORDER BY pd.installment_number ASC;
	`
	rows, err := r.db.QueryContext(ctx, query, id, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var payment output.GetSalesPaymentOutput
		if err := rows.Scan(&payment.Id, &payment.InstallmentNumber, &payment.InstallmentValue, &payment.DueDate, &payment.PaidDate, &payment.PaymentType, &payment.PaymentStatus); err != nil {
			return nil, err
		}
		o = append(o, payment)
	}
	return o, nil
}

func (r *salesRepository) GetItemsBySaleId(ctx context.Context, id int64) ([]output.GetItemsOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var items []output.GetItemsOutput

	query := `
	SELECT
		si.quantity,
		si.sku_id,
		s.price,
		s.id,
		s.code,
		s.color,
		s.size,
		p.name,
		p.description
	FROM sales_items si
	JOIN skus s ON si.sku_id = s.id
	JOIN products p ON s.product_id = p.id
	WHERE si.sales_id = $1 AND si.tenant_id = $2
	ORDER BY si.sku_id ASC;
	`
	rows, err := r.db.QueryContext(ctx, query, id, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item output.GetItemsOutput
		if err := rows.Scan(&item.Quantity, &item.Sku.Id, &item.Sku.Price, &item.Sku.Id, &item.Sku.Code, &item.Sku.Color, &item.Sku.Size, &item.Sku.Product.Name, &item.Sku.Product.Description); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *salesRepository) ChangePaymentStatus(ctx context.Context, id int64, status domain.PaymentStatus) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedId int64
	query := `UPDATE payment_dates SET status = $1 WHERE id = $2 AND tenant_id = $3 RETURNING id`
	err := r.db.QueryRowContext(ctx, query, status, id, tenantId).Scan(&insertedId)
	if err != nil {
		return insertedId, err
	}
	return insertedId, nil
}

func (r *salesRepository) GetPaymentDatesBySaleIdAndPaymentDateId(ctx context.Context, id int64, paymentDateId int64) (domain.SalesPaymentDates, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var output domain.SalesPaymentDates

	query := `
	SELECT
		pd.id,
		pd.installment_number,
		pd.installment_value,
		pd.due_date,
		pd.paid_date,
		pd.status
	FROM payment_dates pd
	JOIN payments p ON pd.payment_id = p.id
	WHERE pd.id = $1 AND pd.tenant_id = $2
	ORDER BY pd.installment_number ASC;
	`
	err := r.db.QueryRowContext(ctx, query, paymentDateId, tenantId).Scan(&output.Id, &output.InstallmentNumber, &output.InstallmentValue, &output.DueDate, &output.PaidDate, &output.Status)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return output, errors.New("Pagamento não encontrado")
		}
		return output, err
	}
	return output, nil
}

func (r *salesRepository) ChangePaymentDate(ctx context.Context, id int64, date *time.Time) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedId int64
	query := `UPDATE payment_dates SET paid_date = $1 WHERE id = $2 AND tenant_id = $3 RETURNING id`
	err := r.db.QueryRowContext(ctx, query, date, id, tenantId).Scan(&insertedId)
	if err != nil {
		return insertedId, err
	}
	return insertedId, nil
}
