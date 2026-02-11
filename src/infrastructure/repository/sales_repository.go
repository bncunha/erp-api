package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/lib/pq"
)

type salesRepository struct {
	db *sql.DB
}

func NewSalesRepository(db *sql.DB) domain.SalesRepository {
	return &salesRepository{db}
}

func (r *salesRepository) CreateSale(ctx context.Context, tx *sql.Tx, sale domain.Sales) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedId int64
	query := `INSERT INTO sales (date, user_id, customer_id, tenant_id, code, last_version) VALUES ($1, $2, $3, $4, $5, 1) RETURNING id`
	err := tx.QueryRowContext(ctx, query, sale.Date, sale.User.Id, sale.Customer.Id, tenantId, sale.Code).Scan(&insertedId)
	if err != nil {
		return insertedId, err
	}
	return insertedId, nil
}

func (r *salesRepository) CreateSaleVersion(ctx context.Context, tx *sql.Tx, saleId int64, version int, date time.Time) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedId int64
	query := `INSERT INTO sales_versions (sales_id, version, date, tenant_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := tx.QueryRowContext(ctx, query, saleId, version, date, tenantId).Scan(&insertedId)
	return insertedId, err
}

func (r *salesRepository) CreateManySaleItem(ctx context.Context, tx *sql.Tx, sale domain.Sales, saleItems []domain.SalesItem) ([]int64, error) {
	if len(saleItems) == 0 {
		return []int64{}, nil
	}

	tenantId := ctx.Value(constants.TENANT_KEY)
	var ids []int64

	query := `INSERT INTO sales_items (quantity, unit_price, sku_id, sales_id, sales_version_id, tenant_id) VALUES %s RETURNING id`
	valueStrings := make([]string, 0, len(saleItems))
	valueArgs := make([]interface{}, 0, len(saleItems)*6)

	for i, item := range saleItems {
		unitPrice := item.UnitPrice
		if unitPrice == 0 {
			unitPrice = item.Sku.Price
		}
		n := i * 6
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d)", n+1, n+2, n+3, n+4, n+5, n+6))
		valueArgs = append(valueArgs,
			item.Quantity,
			unitPrice,
			item.Sku.Id,
			sale.Id,
			sale.SalesVersionId,
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
	query := `INSERT INTO payments (payment_type, sales_id, sales_version_id, tenant_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := tx.QueryRowContext(ctx, query, payment.PaymentType, sale.Id, sale.SalesVersionId, tenantId).Scan(&insertedId)
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

func (r *salesRepository) CreateSalesReturn(ctx context.Context, tx *sql.Tx, saleId int64, fromSalesVersionId int64, toSalesVersionId int64, salesReturn domain.SalesReturn, createdByUserId int64) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedId int64
	query := `INSERT INTO sales_returns (sales_id, from_sales_version_id, to_sales_version_id, return_date, returner_name, reason, created_by_user_id, tenant_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err := tx.QueryRowContext(ctx, query, saleId, fromSalesVersionId, toSalesVersionId, salesReturn.ReturnDate, salesReturn.Returner, salesReturn.Reason, createdByUserId, tenantId).Scan(&insertedId)
	return insertedId, err
}

func (r *salesRepository) CreateSalesReturnItems(ctx context.Context, tx *sql.Tx, salesReturnId int64, items []domain.SalesReturnItem) ([]int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var ids []int64
	valueStrings := make([]string, 0, len(items))
	valueArgs := make([]interface{}, 0, len(items)*5)
	query := `INSERT INTO sales_return_items (sales_return_id, sku_id, quantity, unit_price, tenant_id) VALUES %s RETURNING id`

	for i, item := range items {
		n := i * 5
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", n+1, n+2, n+3, n+4, n+5))
		valueArgs = append(valueArgs, salesReturnId, item.Sku.Id, item.Quantity, item.UnitPrice, tenantId)
	}
	rows, err := tx.QueryContext(ctx, fmt.Sprintf(query, strings.Join(valueStrings, ",")), valueArgs...)
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

func (r *salesRepository) UpdateSaleLastVersion(ctx context.Context, tx *sql.Tx, saleId int64, version int) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	_, err := tx.ExecContext(ctx, `UPDATE sales SET last_version = $1 WHERE id = $2 AND tenant_id = $3`, version, saleId, tenantId)
	return err
}

func (r *salesRepository) CancelPaymentDatesBySaleVersionId(ctx context.Context, tx *sql.Tx, saleVersionId int64) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `
	UPDATE payment_dates pd
	SET status = $1
	FROM payments p
	WHERE pd.payment_id = p.id
	  AND p.sales_version_id = $2
	  AND pd.tenant_id = $3`
	_, err := tx.ExecContext(ctx, query, domain.PaymentStatusCancel, saleVersionId, tenantId)
	return err
}

func (r *salesRepository) GetSaleByIdForUpdate(ctx context.Context, tx *sql.Tx, id int64) (domain.SaleWithVersionOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var output domain.SaleWithVersionOutput
	query := `
	SELECT s.id, s.code, s.date, s.user_id, s.customer_id, c.name, s.last_version, sv.id
	FROM sales s
	JOIN customers c ON c.id = s.customer_id AND c.tenant_id = s.tenant_id
	JOIN sales_versions sv ON sv.sales_id = s.id AND sv.version = s.last_version AND sv.tenant_id = s.tenant_id
	WHERE s.id = $1 AND s.tenant_id = $2
	FOR UPDATE`
	err := tx.QueryRowContext(ctx, query, id, tenantId).Scan(
		&output.Id,
		&output.Code,
		&output.Date,
		&output.UserId,
		&output.CustomerId,
		&output.CustomerName,
		&output.LastVersion,
		&output.SalesVersionId,
	)
	return output, err
}

func (r *salesRepository) GetSaleVersionIdBySaleIdAndVersion(ctx context.Context, tx *sql.Tx, saleId int64, version int) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var id int64
	err := tx.QueryRowContext(ctx, `SELECT id FROM sales_versions WHERE sales_id = $1 AND version = $2 AND tenant_id = $3`, saleId, version, tenantId).Scan(&id)
	return id, err
}

func (r *salesRepository) GetSales(ctx context.Context, input domain.GetSalesInput) ([]domain.GetSalesItemOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var sales []domain.GetSalesItemOutput

	query := `
SELECT
  s.id,
  s.date,
  c.name AS customer,
  u.name AS seller,
  COALESCE(SUM(pd.installment_value), 0) AS total_value,
  COALESCE(SUM(pd.installment_value) FILTER (WHERE pd.status IN ('PAID','REVERSAL')), 0) AS received_value,
  COALESCE(SUM(pd.installment_value) FILTER (WHERE pd.status IN ('PENDING','DELAYED')), 0) AS future_revenue,
  (
    SELECT COALESCE(SUM(si.quantity), 0)
    FROM sales_items si
    WHERE si.sales_version_id = sv.id
      AND si.tenant_id = s.tenant_id
  ) AS total_items,
  CASE
    WHEN (
      SELECT COALESCE(SUM(si.quantity), 0)
      FROM sales_items si
      WHERE si.sales_version_id = sv.id
        AND si.tenant_id = s.tenant_id
    ) = 0
      THEN 'CANCEL'
    WHEN COUNT(pd.id) FILTER (WHERE pd.status IN ('PAID','PENDING','DELAYED')) = 0
      THEN 'CANCEL'
    WHEN COUNT(pd.id) FILTER (WHERE pd.status IN ('PAID','PENDING','DELAYED')) > 0
      AND COUNT(pd.id) FILTER (WHERE pd.status IN ('PAID','PENDING','DELAYED')) = COUNT(pd.id) FILTER (WHERE pd.status = 'PAID')
      THEN 'PAID'
    WHEN COALESCE(BOOL_OR(pd.status = 'PENDING' AND pd.due_date < CURRENT_DATE), FALSE)
      THEN 'DELAYED'
    ELSE 'PENDING'
  END AS summary_status
FROM sales s
JOIN sales_versions sv ON sv.sales_id = s.id AND sv.version = s.last_version AND sv.tenant_id = s.tenant_id
JOIN users u ON u.id = s.user_id AND u.tenant_id = s.tenant_id
JOIN customers c ON c.id = s.customer_id AND c.tenant_id = s.tenant_id
LEFT JOIN payments p ON p.sales_version_id = sv.id AND p.tenant_id = s.tenant_id
LEFT JOIN payment_dates pd ON pd.payment_id = p.id AND pd.tenant_id = s.tenant_id
WHERE s.tenant_id = $1
  AND ($2::bigint[] IS NULL OR s.user_id = ANY($2))
  AND ($3::bigint[] IS NULL OR s.customer_id = ANY($3))
  AND ($4::timestamptz IS NULL OR s.date >= $4)
  AND ($5::timestamptz IS NULL OR s.date <= $5)
GROUP BY s.id, s.date, c.name, u.name, sv.id
HAVING
  $6::text IS NULL
  OR (
      ($6 = 'CANCEL' AND (
        (
          SELECT COALESCE(SUM(si.quantity), 0)
          FROM sales_items si
          WHERE si.sales_version_id = sv.id
            AND si.tenant_id = s.tenant_id
        ) = 0
        OR COUNT(pd.id) FILTER (WHERE pd.status IN ('PAID','PENDING','DELAYED')) = 0
      ))
      OR ($6 = 'PAID'
        AND (
          SELECT COALESCE(SUM(si.quantity), 0)
          FROM sales_items si
          WHERE si.sales_version_id = sv.id
            AND si.tenant_id = s.tenant_id
        ) > 0
        AND COUNT(pd.id) FILTER (WHERE pd.status IN ('PAID','PENDING','DELAYED')) > 0
        AND COUNT(pd.id) FILTER (WHERE pd.status IN ('PAID','PENDING','DELAYED')) = COUNT(pd.id) FILTER (WHERE pd.status = 'PAID')
      )
      OR ($6 = 'DELAYED'
        AND (
          SELECT COALESCE(SUM(si.quantity), 0)
          FROM sales_items si
          WHERE si.sales_version_id = sv.id
            AND si.tenant_id = s.tenant_id
        ) > 0
        AND COALESCE(BOOL_OR(pd.status = 'PENDING' AND pd.due_date < CURRENT_DATE), FALSE)
      )
      OR ($6 = 'PENDING'
        AND (
          SELECT COALESCE(SUM(si.quantity), 0)
          FROM sales_items si
          WHERE si.sales_version_id = sv.id
            AND si.tenant_id = s.tenant_id
        ) > 0
        AND COUNT(pd.id) FILTER (WHERE pd.status IN ('PAID','PENDING','DELAYED')) > 0
        AND NOT (
          (COUNT(pd.id) FILTER (WHERE pd.status IN ('PAID','PENDING','DELAYED')) = COUNT(pd.id) FILTER (WHERE pd.status = 'PAID'))
          OR COALESCE(BOOL_OR(pd.status = 'PENDING' AND pd.due_date < CURRENT_DATE), FALSE)
        )
      )
  )
ORDER BY s.date DESC, s.id DESC`
	valueArgs := []interface{}{tenantId, pq.Array(input.UserId), pq.Array(input.CustomerId), input.InitialDate, input.FinalDate, input.PaymentStatus}

	rows, err := r.db.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sale domain.GetSalesItemOutput
		if err := rows.Scan(&sale.Id, &sale.Date, &sale.CustomerName, &sale.SellerName, &sale.TotalValue, &sale.ReceivedValue, &sale.FutureRevenue, &sale.TotalItems, &sale.Status); err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}
	return sales, nil
}

func (r *salesRepository) GetSaleById(ctx context.Context, id int64) (domain.GetSaleByIdOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var output domain.GetSaleByIdOutput

	query := `
	SELECT
		s.id,
		s.code,
		s.date,
		u.name AS seller_name,
		c.name AS customer_name,
		COALESCE(SUM(pd.installment_value), 0) AS total_value,
		COALESCE(SUM(pd.installment_value) FILTER (WHERE pd.status IN ('PAID','REVERSAL')), 0) AS received_value,
		COALESCE(SUM(pd.installment_value) FILTER (WHERE pd.status IN ('PENDING','DELAYED')), 0) AS future_revenue,
		CASE
			WHEN (
				SELECT COALESCE(SUM(si.quantity), 0)
				FROM sales_items si
				WHERE si.sales_version_id = sv.id
				  AND si.tenant_id = s.tenant_id
			) = 0 THEN 'CANCEL'
			WHEN COUNT(pd.id) FILTER (WHERE pd.status IN ('PAID','PENDING','DELAYED')) = 0 THEN 'CANCEL'
			WHEN COUNT(pd.id) FILTER (WHERE pd.status IN ('PAID','PENDING','DELAYED')) > 0
			  AND COUNT(pd.id) FILTER (WHERE pd.status IN ('PAID','PENDING','DELAYED')) = COUNT(pd.id) FILTER (WHERE pd.status = 'PAID') THEN 'PAID'
			WHEN COALESCE(BOOL_OR(pd.status = 'PENDING' AND pd.due_date < CURRENT_DATE), FALSE) THEN 'DELAYED'
			ELSE 'PENDING'
		END AS payment_status
	FROM sales s
	JOIN sales_versions sv ON sv.sales_id = s.id AND sv.version = s.last_version AND sv.tenant_id = s.tenant_id
	JOIN users u ON u.id = s.user_id AND u.tenant_id = s.tenant_id
	JOIN customers c ON c.id = s.customer_id AND c.tenant_id = s.tenant_id
	LEFT JOIN payments p ON p.sales_version_id = sv.id AND p.tenant_id = s.tenant_id
	LEFT JOIN payment_dates pd ON pd.payment_id = p.id AND pd.tenant_id = s.tenant_id
	WHERE s.id = $1
	  AND s.tenant_id = $2
	GROUP BY s.id, s.code, s.date, u.name, c.name`
	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&output.Id, &output.Code, &output.Date, &output.SellerName, &output.CustomerName, &output.TotalValue, &output.ReceivedValue, &output.FutureRevenue, &output.PaymentStatus)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (r *salesRepository) GetPaymentsBySaleId(ctx context.Context, id int64) ([]domain.GetSalesPaymentOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var o []domain.GetSalesPaymentOutput

	query := `
	SELECT
		pd.id,
		pd.installment_number,
		pd.installment_value,
		pd.due_date,
		pd.paid_date,
		p.payment_type,
		CASE
			WHEN pd.status = 'PENDING' AND pd.due_date < CURRENT_DATE
				THEN 'DELAYED'
				ELSE pd.status
		END AS status
	FROM sales s
	JOIN sales_versions sv ON sv.sales_id = s.id AND sv.version = s.last_version AND sv.tenant_id = s.tenant_id
	JOIN payments p ON p.sales_version_id = sv.id AND p.tenant_id = s.tenant_id
	JOIN payment_dates pd ON p.id = pd.payment_id AND pd.tenant_id = s.tenant_id
	WHERE s.id = $1 AND s.tenant_id = $2
	ORDER BY p.payment_type, pd.installment_number ASC`
	rows, err := r.db.QueryContext(ctx, query, id, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var payment domain.GetSalesPaymentOutput
		if err := rows.Scan(&payment.Id, &payment.InstallmentNumber, &payment.InstallmentValue, &payment.DueDate, &payment.PaidDate, &payment.PaymentType, &payment.PaymentStatus); err != nil {
			return nil, err
		}
		o = append(o, payment)
	}
	return o, nil
}

func (r *salesRepository) GetPaymentsBySaleVersionId(ctx context.Context, saleVersionId int64) ([]domain.GetSalesPaymentOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var o []domain.GetSalesPaymentOutput
	query := `
	SELECT
		pd.id,
		pd.installment_number,
		pd.installment_value,
		pd.due_date,
		pd.paid_date,
		p.payment_type,
		CASE
			WHEN pd.status = 'PENDING' AND pd.due_date < CURRENT_DATE
				THEN 'DELAYED'
				ELSE pd.status
		END AS status
	FROM payments p
	JOIN payment_dates pd ON p.id = pd.payment_id
	WHERE p.sales_version_id = $1 AND p.tenant_id = $2
	ORDER BY p.payment_type, pd.installment_number ASC`
	rows, err := r.db.QueryContext(ctx, query, saleVersionId, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var payment domain.GetSalesPaymentOutput
		if err := rows.Scan(&payment.Id, &payment.InstallmentNumber, &payment.InstallmentValue, &payment.DueDate, &payment.PaidDate, &payment.PaymentType, &payment.PaymentStatus); err != nil {
			return nil, err
		}
		o = append(o, payment)
	}
	return o, nil
}

func (r *salesRepository) GetItemsBySaleId(ctx context.Context, id int64) ([]domain.GetItemsOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var items []domain.GetItemsOutput

	query := `
	SELECT
		si.quantity,
		si.sku_id,
		sku.price,
		sku.id,
		sku.code,
		sku.color,
		sku.size,
		p.name,
		p.description
	FROM sales s
	JOIN sales_versions sv ON sv.sales_id = s.id AND sv.version = s.last_version AND sv.tenant_id = s.tenant_id
	JOIN sales_items si ON si.sales_version_id = sv.id AND si.tenant_id = s.tenant_id
	JOIN skus sku ON si.sku_id = sku.id
	JOIN products p ON sku.product_id = p.id
	WHERE s.id = $1 AND s.tenant_id = $2
	ORDER BY si.sku_id ASC`
	rows, err := r.db.QueryContext(ctx, query, id, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item domain.GetItemsOutput
		if err := rows.Scan(&item.Quantity, &item.Sku.Id, &item.Sku.Price, &item.Sku.Id, &item.Sku.Code, &item.Sku.Color, &item.Sku.Size, &item.Sku.Product.Name, &item.Sku.Product.Description); err != nil {
			return nil, err
		}
		item.UnitPrice = item.Sku.Price
		items = append(items, item)
	}
	return items, nil
}

func (r *salesRepository) GetItemsBySaleVersionId(ctx context.Context, saleVersionId int64) ([]domain.GetItemsOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var items []domain.GetItemsOutput
	query := `
	SELECT
		si.quantity,
		si.unit_price,
		si.sku_id,
		s.id,
		s.code,
		s.color,
		s.size,
		p.name,
		p.description
	FROM sales_items si
	JOIN skus s ON si.sku_id = s.id
	JOIN products p ON s.product_id = p.id
	WHERE si.sales_version_id = $1 AND si.tenant_id = $2
	ORDER BY si.sku_id ASC`
	rows, err := r.db.QueryContext(ctx, query, saleVersionId, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item domain.GetItemsOutput
		if err := rows.Scan(&item.Quantity, &item.UnitPrice, &item.Sku.Id, &item.Sku.Id, &item.Sku.Code, &item.Sku.Color, &item.Sku.Size, &item.Sku.Product.Name, &item.Sku.Product.Description); err != nil {
			return nil, err
		}
		item.Sku.Price = item.UnitPrice
		items = append(items, item)
	}
	return items, nil
}

func (r *salesRepository) GetReturnsBySaleId(ctx context.Context, id int64) ([]domain.GetSalesReturnOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	output := make([]domain.GetSalesReturnOutput, 0)
	query := `
	SELECT id, return_date, returner_name, reason
	FROM sales_returns
	WHERE sales_id = $1 AND tenant_id = $2
	ORDER BY return_date DESC, id DESC`
	rows, err := r.db.QueryContext(ctx, query, id, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ret domain.GetSalesReturnOutput
		if err := rows.Scan(&ret.Id, &ret.ReturnDate, &ret.Returner, &ret.Reason); err != nil {
			return nil, err
		}
		items, err := r.getSalesReturnItems(ctx, ret.Id)
		if err != nil {
			return nil, err
		}
		ret.Items = items
		output = append(output, ret)
	}
	return output, nil
}

func (r *salesRepository) getSalesReturnItems(ctx context.Context, salesReturnId int64) ([]domain.GetSalesReturnItemOutput, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	output := make([]domain.GetSalesReturnItemOutput, 0)
	query := `
	SELECT sri.sku_id, sri.quantity, sri.unit_price, s.code, s.color, s.size, p.name, p.description
	FROM sales_return_items sri
	JOIN skus s ON s.id = sri.sku_id
	JOIN products p ON p.id = s.product_id
	WHERE sri.sales_return_id = $1 AND sri.tenant_id = $2
	ORDER BY sri.sku_id ASC`
	rows, err := r.db.QueryContext(ctx, query, salesReturnId, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item domain.GetSalesReturnItemOutput
		if err := rows.Scan(&item.Sku.Id, &item.Quantity, &item.UnitPrice, &item.Sku.Code, &item.Sku.Color, &item.Sku.Size, &item.Sku.Product.Name, &item.Sku.Product.Description); err != nil {
			return nil, err
		}
		item.Sku.Price = item.UnitPrice
		output = append(output, item)
	}
	return output, nil
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
			pd.status,
			p.payment_type
		FROM payment_dates pd
		JOIN payments p ON pd.payment_id = p.id
		JOIN sales s ON s.id = p.sales_id AND s.tenant_id = p.tenant_id
		WHERE pd.id = $1 AND pd.tenant_id = $2 AND s.id = $3
		ORDER BY pd.installment_number ASC`
	err := r.db.QueryRowContext(ctx, query, paymentDateId, tenantId, id).Scan(&output.Id, &output.InstallmentNumber, &output.InstallmentValue, &output.DueDate, &output.PaidDate, &output.Status, &output.PaymentType)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return output, errors.New("Pagamento nÃ£o encontrado")
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
