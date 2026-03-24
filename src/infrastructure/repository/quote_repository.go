package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type quoteRepository struct {
	db *sql.DB
}

func NewQuoteRepository(db *sql.DB) domain.QuoteRepository {
	return &quoteRepository{db: db}
}

func (r *quoteRepository) Create(ctx context.Context, tx *sql.Tx, quote *domain.Quote) error {
	tenantID := ctx.Value(constants.TENANT_KEY)
	quoteNumber, err := r.nextQuoteNumber(ctx, tx)
	if err != nil {
		return err
	}
	quote.QuoteNumber = quoteNumber

	query := `INSERT INTO quotes (
		tenant_id, customer_id, quote_number, status, valid_until, down_payment_percentage,
		discount_percentage, notes, shipping_type, shipping_cost, shipping_region, shipping_min_value, shipping_description,
		subtotal_amount, total_amount, created_by_user_id, production_order_id
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	RETURNING id, created_at, updated_at`
	if err := queryRowContextOrDB(ctx, tx, r.db, query,
		tenantID,
		quote.CustomerID,
		quote.QuoteNumber,
		quote.Status,
		quote.ValidUntil,
		quote.DownPaymentPercentage,
		quote.DiscountPercentage,
		nullableString(quote.Notes),
		quote.ShippingType,
		quote.ShippingCost,
		nullableString(quote.ShippingRegion),
		nullableFloat64(quote.ShippingMinValue),
		quote.ShippingDescription,
		quote.SubtotalAmount,
		quote.TotalAmount,
		quote.CreatedByUserID,
		nullableInt64(quote.ProductionOrderID),
	).Scan(&quote.ID, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
		return err
	}

	return r.replaceItems(ctx, tx, quote.ID, quote.Items)
}

func (r *quoteRepository) Update(ctx context.Context, tx *sql.Tx, quote *domain.Quote) error {
	tenantID := ctx.Value(constants.TENANT_KEY)
	query := `UPDATE quotes
	SET customer_id=$1, valid_until=$2, down_payment_percentage=$3, discount_percentage=$4, notes=$5,
		shipping_type=$6, shipping_cost=$7, shipping_region=$8, shipping_min_value=$9, shipping_description=$10,
		subtotal_amount=$11, total_amount=$12, updated_at=NOW()
	WHERE id=$13 AND tenant_id=$14`
	result, err := execContextOrDB(ctx, tx, r.db, query,
		quote.CustomerID,
		quote.ValidUntil,
		quote.DownPaymentPercentage,
		quote.DiscountPercentage,
		nullableString(quote.Notes),
		quote.ShippingType,
		quote.ShippingCost,
		nullableString(quote.ShippingRegion),
		nullableFloat64(quote.ShippingMinValue),
		quote.ShippingDescription,
		quote.SubtotalAmount,
		quote.TotalAmount,
		quote.ID,
		tenantID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("orÃ§amento nÃ£o encontrado")
	}
	return r.replaceItems(ctx, tx, quote.ID, quote.Items)
}

func (r *quoteRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, id int64, status domain.QuoteStatus) error {
	tenantID := ctx.Value(constants.TENANT_KEY)
	query := `UPDATE quotes SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3`
	result, err := execContextOrDB(ctx, tx, r.db, query, status, id, tenantID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("orÃ§amento nÃ£o encontrado")
	}
	return nil
}

func (r *quoteRepository) GetByID(ctx context.Context, id int64) (domain.Quote, error) {
	tenantID := ctx.Value(constants.TENANT_KEY)
	var quote domain.Quote
	query := `SELECT q.id, q.tenant_id, q.customer_id, q.quote_number, q.status, q.valid_until,
	q.down_payment_percentage, q.discount_percentage, q.notes, q.shipping_type, q.shipping_cost, q.shipping_region, q.shipping_min_value,
	q.shipping_description, q.subtotal_amount, q.total_amount, q.created_by_user_id,
	q.production_order_id, q.created_at, q.updated_at,
	c.id, c.name, c.phone_number
	FROM quotes q
	JOIN customers c ON c.id=q.customer_id AND c.tenant_id=q.tenant_id
	WHERE q.id=$1 AND q.tenant_id=$2`
	var notes, shippingRegion sql.NullString
	var shippingMinValue sql.NullFloat64
	var productionOrderID sql.NullInt64
	if err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&quote.ID,
		&quote.TenantID,
		&quote.CustomerID,
		&quote.QuoteNumber,
		&quote.Status,
		&quote.ValidUntil,
		&quote.DownPaymentPercentage,
		&quote.DiscountPercentage,
		&notes,
		&quote.ShippingType,
		&quote.ShippingCost,
		&shippingRegion,
		&shippingMinValue,
		&quote.ShippingDescription,
		&quote.SubtotalAmount,
		&quote.TotalAmount,
		&quote.CreatedByUserID,
		&productionOrderID,
		&quote.CreatedAt,
		&quote.UpdatedAt,
		&quote.Customer.ID,
		&quote.Customer.Name,
		&quote.Customer.Phone,
	); err != nil {
		if errors.IsNoRowsFinded(err) {
			return quote, errors.New("orÃ§amento nÃ£o encontrado")
		}
		return quote, err
	}
	if notes.Valid {
		quote.Notes = &notes.String
	}
	if shippingRegion.Valid {
		quote.ShippingRegion = &shippingRegion.String
	}
	if shippingMinValue.Valid {
		quote.ShippingMinValue = &shippingMinValue.Float64
	}
	if productionOrderID.Valid {
		quote.ProductionOrderID = &productionOrderID.Int64
	}

	items, err := r.getItemsByQuoteID(ctx, id)
	if err != nil {
		return quote, err
	}
	quote.Items = items
	return quote, nil
}

func (r *quoteRepository) List(ctx context.Context, input domain.GetQuotesInput) (domain.GetQuotesOutput, error) {
	tenantID := ctx.Value(constants.TENANT_KEY)
	output := domain.GetQuotesOutput{Page: input.Page, PageSize: input.PageSize}
	if output.Page < 1 {
		output.Page = 1
	}
	if output.PageSize < 1 {
		output.PageSize = 10
	}
	if output.PageSize > 100 {
		output.PageSize = 100
	}
	offset := (output.Page - 1) * output.PageSize

	where := []string{"q.tenant_id=$1"}
	args := []any{tenantID}
	argPos := 2

	if len(input.Statuses) > 0 {
		statusPlaceholders := make([]string, 0, len(input.Statuses))
		for _, status := range input.Statuses {
			statusPlaceholders = append(statusPlaceholders, fmt.Sprintf("$%d", argPos))
			args = append(args, status)
			argPos++
		}
		where = append(where, fmt.Sprintf("q.status IN (%s)", strings.Join(statusPlaceholders, ",")))
	}
	if input.CustomerID != nil {
		where = append(where, fmt.Sprintf("q.customer_id=$%d", argPos))
		args = append(args, *input.CustomerID)
		argPos++
	}
	if input.ValidUntilStart != nil {
		where = append(where, fmt.Sprintf("q.valid_until >= $%d", argPos))
		args = append(args, *input.ValidUntilStart)
		argPos++
	}
	if input.ValidUntilEnd != nil {
		where = append(where, fmt.Sprintf("q.valid_until <= $%d", argPos))
		args = append(args, *input.ValidUntilEnd)
		argPos++
	}
	if input.CreatedAtStart != nil {
		where = append(where, fmt.Sprintf("q.created_at >= $%d", argPos))
		args = append(args, *input.CreatedAtStart)
		argPos++
	}
	if input.CreatedAtEnd != nil {
		where = append(where, fmt.Sprintf("q.created_at <= $%d", argPos))
		args = append(args, *input.CreatedAtEnd)
		argPos++
	}
	if input.Search != nil && strings.TrimSpace(*input.Search) != "" {
		where = append(where, fmt.Sprintf("(q.quote_number ILIKE $%d OR c.name ILIKE $%d OR COALESCE(q.notes,'') ILIKE $%d OR EXISTS (SELECT 1 FROM quote_items qi WHERE qi.quote_id=q.id AND qi.product_description_snapshot ILIKE $%d))", argPos, argPos, argPos, argPos))
		args = append(args, "%"+strings.TrimSpace(*input.Search)+"%")
		argPos++
	}
	if input.OnlyExpired {
		where = append(where, "q.valid_until < CURRENT_DATE AND q.status NOT IN ('APPROVED','REJECTED','CANCELED')")
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM quotes q JOIN customers c ON c.id=q.customer_id AND c.tenant_id=q.tenant_id WHERE %s`, strings.Join(where, " AND "))
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&output.Total); err != nil {
		return output, err
	}
	output.TotalPages = int(math.Ceil(float64(output.Total) / float64(output.PageSize)))

	sortBy := normalizeQuoteSortBy(input.SortBy)
	sortOrder := normalizeSortOrder(input.SortOrder)
	listQuery := fmt.Sprintf(`SELECT q.id, q.quote_number, c.id, c.name, c.phone_number, q.status, q.valid_until,
	q.total_amount, q.down_payment_percentage, q.shipping_description, q.created_at
	FROM quotes q
	JOIN customers c ON c.id=q.customer_id AND c.tenant_id=q.tenant_id
	WHERE %s
	ORDER BY %s %s
	LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), sortBy, sortOrder, argPos, argPos+1)
	args = append(args, output.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	items := make([]domain.QuoteListItemOutput, 0, output.PageSize)
	for rows.Next() {
		var item domain.QuoteListItemOutput
		if err := rows.Scan(
			&item.ID,
			&item.QuoteNumber,
			&item.Customer.ID,
			&item.Customer.Name,
			&item.Customer.Phone,
			&item.Status,
			&item.ValidUntil,
			&item.TotalAmount,
			&item.DownPaymentPercentage,
			&item.ShippingDescription,
			&item.CreatedAt,
		); err != nil {
			return output, err
		}
		item.DownPaymentAmount = item.TotalAmount * item.DownPaymentPercentage / 100
		items = append(items, item)
	}
	output.Items = items
	return output, nil
}

func (r *quoteRepository) GetCompanyProfile(ctx context.Context) (domain.QuoteCompanyProfile, error) {
	tenantID := ctx.Value(constants.TENANT_KEY)
	profile := domain.QuoteCompanyProfile{}
	query := `SELECT c.name, c.legal_name, COALESCE(c.cnpj,''), COALESCE(c.email,''), COALESCE(c.cellphone,''), c.logo_url,
	COALESCE(a.street,'') || CASE WHEN a.number IS NOT NULL AND a.number <> '' THEN ', ' || a.number ELSE '' END ||
	CASE WHEN a.neighborhood IS NOT NULL AND a.neighborhood <> '' THEN ' - ' || a.neighborhood ELSE '' END ||
	CASE WHEN a.city IS NOT NULL AND a.city <> '' THEN ' - ' || a.city ELSE '' END ||
	CASE WHEN a.uf IS NOT NULL AND a.uf <> '' THEN '/' || a.uf ELSE '' END ||
	CASE WHEN a.cep IS NOT NULL AND a.cep <> '' THEN ' - CEP ' || a.cep ELSE '' END AS full_address
	FROM companies c
	LEFT JOIN LATERAL (
		SELECT street, number, neighborhood, city, uf, cep
		FROM addresses
		WHERE tenant_id=c.id
		ORDER BY id DESC
		LIMIT 1
	) a ON TRUE
	WHERE c.id=$1`
	var logoURL sql.NullString
	if err := r.db.QueryRowContext(ctx, query, tenantID).Scan(
		&profile.Name,
		&profile.LegalName,
		&profile.CNPJ,
		&profile.Email,
		&profile.Phone,
		&logoURL,
		&profile.Address,
	); err != nil {
		if errors.IsNoRowsFinded(err) {
			return profile, errors.New("empresa não encontrada")
		}
		return profile, err
	}
	if logoURL.Valid {
		profile.LogoURL = &logoURL.String
	}
	return profile, nil
}

func (r *quoteRepository) getItemsByQuoteID(ctx context.Context, quoteID int64) ([]domain.QuoteItem, error) {
	tenantID := ctx.Value(constants.TENANT_KEY)
	query := `SELECT qi.id, qi.quote_id, qi.sku_id, COALESCE(s.product_id,0), COALESCE(s.code,''),
	qi.product_description_snapshot, qi.quantity, qi.unit_price, qi.total_price, qi.created_at, qi.updated_at
	FROM quote_items qi
	JOIN quotes q ON q.id=qi.quote_id AND q.tenant_id=$2
	LEFT JOIN skus s ON s.id=qi.sku_id
	WHERE qi.quote_id=$1
	ORDER BY qi.id ASC`
	rows, err := r.db.QueryContext(ctx, query, quoteID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.QuoteItem, 0)
	for rows.Next() {
		var item domain.QuoteItem
		if err := rows.Scan(
			&item.ID,
			&item.QuoteID,
			&item.SkuID,
			&item.ProductID,
			&item.SkuCode,
			&item.ProductDescriptionSnapshot,
			&item.Quantity,
			&item.UnitPrice,
			&item.TotalPrice,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *quoteRepository) replaceItems(ctx context.Context, tx *sql.Tx, quoteID int64, items []domain.QuoteItem) error {
	if _, err := execContextOrDB(ctx, tx, r.db, `DELETE FROM quote_items WHERE quote_id=$1`, quoteID); err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	query := `INSERT INTO quote_items (quote_id, sku_id, product_description_snapshot, quantity, unit_price, total_price)
	VALUES %s`
	valueStrings := make([]string, 0, len(items))
	valueArgs := make([]interface{}, 0, len(items)*6)
	for i, item := range items {
		n := i * 6
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d)", n+1, n+2, n+3, n+4, n+5, n+6))
		valueArgs = append(valueArgs,
			quoteID,
			item.SkuID,
			item.ProductDescriptionSnapshot,
			item.Quantity,
			item.UnitPrice,
			item.TotalPrice,
		)
	}
	_, err := execContextOrDB(ctx, tx, r.db, fmt.Sprintf(query, strings.Join(valueStrings, ",")), valueArgs...)
	return err
}

func (r *quoteRepository) nextQuoteNumber(ctx context.Context, tx *sql.Tx) (string, error) {
	tenantID := ctx.Value(constants.TENANT_KEY)
	query := `INSERT INTO quote_number_sequences (tenant_id, last_number, updated_at)
	VALUES ($1, 1, NOW())
	ON CONFLICT (tenant_id)
	DO UPDATE SET last_number = quote_number_sequences.last_number + 1, updated_at=NOW()
	RETURNING last_number`
	var number int64
	if err := queryRowContextOrDB(ctx, tx, r.db, query, tenantID).Scan(&number); err != nil {
		return "", err
	}
	return fmt.Sprintf("O-%06d", number), nil
}

func normalizeQuoteSortBy(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "valid_until":
		return "q.valid_until"
	case "total_amount":
		return "q.total_amount"
	default:
		return "q.created_at"
	}
}

func normalizeSortOrder(order string) string {
	if strings.EqualFold(strings.TrimSpace(order), "asc") {
		return "ASC"
	}
	return "DESC"
}

func nullableFloat64(value *float64) any {
	if value == nil {
		return nil
	}
	return *value
}

func nullableInt64(value *int64) any {
	if value == nil {
		return nil
	}
	return *value
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func queryRowContextOrDB(ctx context.Context, tx *sql.Tx, db *sql.DB, query string, args ...any) *sql.Row {
	if tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return db.QueryRowContext(ctx, query, args...)
}

func execContextOrDB(ctx context.Context, tx *sql.Tx, db *sql.DB, query string, args ...any) (sql.Result, error) {
	if tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return db.ExecContext(ctx, query, args...)
}
