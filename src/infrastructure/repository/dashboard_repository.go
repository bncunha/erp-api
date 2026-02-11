package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
)

type dashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) domain.DashboardRepository {
	return &dashboardRepository{db}
}

func (r *dashboardRepository) GetRevenue(ctx context.Context, input domain.DashboardQueryInput) (float64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var total float64

	query := `
	SELECT COALESCE(SUM(si.quantity * si.unit_price), 0)
	FROM sales s
	JOIN sales_versions sv ON sv.sales_id = s.id AND sv.version = s.last_version AND sv.tenant_id = s.tenant_id
	JOIN sales_items si ON si.sales_version_id = sv.id AND si.tenant_id = s.tenant_id
	JOIN skus sk ON sk.id = si.sku_id AND sk.tenant_id = si.tenant_id
	WHERE si.tenant_id = $1
	  AND ($2::bigint IS NULL OR s.user_id = $2)
	  AND ($3::bigint IS NULL OR sk.product_id = $3)
	  AND s.date >= $4
	  AND s.date <= $5`

	err := r.db.QueryRowContext(ctx, query, tenantId, input.ResellerId, input.ProductId, input.From, input.To).Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetSalesCount(ctx context.Context, input domain.DashboardQueryInput) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var total int64

	query := `
	SELECT COALESCE(COUNT(DISTINCT s.id), 0)
	FROM sales s
	JOIN sales_versions sv ON sv.sales_id = s.id AND sv.version = s.last_version AND sv.tenant_id = s.tenant_id
	JOIN sales_items si ON si.sales_version_id = sv.id AND si.tenant_id = s.tenant_id
	JOIN skus sk ON sk.id = si.sku_id AND sk.tenant_id = s.tenant_id
	WHERE s.tenant_id = $1
	  AND ($2::bigint IS NULL OR s.user_id = $2)
	  AND ($3::bigint IS NULL OR sk.product_id = $3)
	  AND s.date >= $4
	  AND s.date <= $5`

	err := r.db.QueryRowContext(ctx, query, tenantId, input.ResellerId, input.ProductId, input.From, input.To).Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetRevenueByDay(ctx context.Context, input domain.DashboardQueryInput) ([]domain.DashboardTimeSeriesItem, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	items := make([]domain.DashboardTimeSeriesItem, 0)

	query := `
	SELECT date_trunc('day', s.date) AS day, COALESCE(SUM(si.quantity * si.unit_price), 0)
	FROM sales s
	JOIN sales_versions sv ON sv.sales_id = s.id AND sv.version = s.last_version AND sv.tenant_id = s.tenant_id
	JOIN sales_items si ON si.sales_version_id = sv.id AND si.tenant_id = s.tenant_id
	JOIN skus sk ON sk.id = si.sku_id AND sk.tenant_id = s.tenant_id
	WHERE s.tenant_id = $1
	  AND ($2::bigint IS NULL OR s.user_id = $2)
	  AND ($3::bigint IS NULL OR sk.product_id = $3)
	  AND s.date >= $4
	  AND s.date <= $5
	GROUP BY day
	ORDER BY day ASC`

	rows, err := r.db.QueryContext(ctx, query, tenantId, input.ResellerId, input.ProductId, input.From, input.To)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.DashboardTimeSeriesItem
		if err := rows.Scan(&item.Date, &item.Value); err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *dashboardRepository) GetSalesCountByDay(ctx context.Context, input domain.DashboardQueryInput) ([]domain.DashboardTimeSeriesItem, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	items := make([]domain.DashboardTimeSeriesItem, 0)

	query := `
	SELECT date_trunc('day', s.date) AS day, COALESCE(COUNT(DISTINCT s.id), 0)
	FROM sales s
	JOIN sales_versions sv ON sv.sales_id = s.id AND sv.version = s.last_version AND sv.tenant_id = s.tenant_id
	JOIN sales_items si ON si.sales_version_id = sv.id AND si.tenant_id = s.tenant_id
	JOIN skus sk ON sk.id = si.sku_id AND sk.tenant_id = s.tenant_id
	WHERE s.tenant_id = $1
	  AND ($2::bigint IS NULL OR s.user_id = $2)
	  AND ($3::bigint IS NULL OR sk.product_id = $3)
	  AND s.date >= $4
	  AND s.date <= $5
	GROUP BY day
	ORDER BY day ASC`

	rows, err := r.db.QueryContext(ctx, query, tenantId, input.ResellerId, input.ProductId, input.From, input.To)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.DashboardTimeSeriesItem
		if err := rows.Scan(&item.Date, &item.Value); err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *dashboardRepository) GetStockTotal(ctx context.Context, input domain.DashboardStockQueryInput) (float64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var total float64

	query := `
	SELECT COALESCE(SUM(ii.quantity), 0)
	FROM inventory_items ii
	JOIN inventories inv ON inv.id = ii.inventory_id AND inv.tenant_id = ii.tenant_id
	WHERE ii.tenant_id = $1
	  AND ii.deleted_at IS NULL
	  AND ($2::bigint IS NULL OR inv.user_id = $2)`

	err := r.db.QueryRowContext(ctx, query, tenantId, input.ResellerId).Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetLowStockProducts(ctx context.Context, input domain.DashboardStockQueryInput) ([]domain.DashboardLowStockItem, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	items := make([]domain.DashboardLowStockItem, 0)

	query := `
	SELECT p.id, p.name, COALESCE(SUM(ii.quantity), 0) AS qty
	FROM products p
	JOIN skus s ON s.product_id = p.id AND s.tenant_id = p.tenant_id AND s.deleted_at IS NULL
	LEFT JOIN inventory_items ii ON ii.sku_id = s.id AND ii.tenant_id = p.tenant_id AND ii.deleted_at IS NULL
	LEFT JOIN inventories inv ON inv.id = ii.inventory_id AND inv.tenant_id = p.tenant_id
	WHERE p.tenant_id = $1 AND p.deleted_at IS NULL
	  AND ($2::bigint IS NULL OR inv.user_id = $2)
	GROUP BY p.id, p.name
	HAVING COALESCE(SUM(ii.quantity), 0) <= $3
	ORDER BY qty ASC, p.name ASC`

	rows, err := r.db.QueryContext(ctx, query, tenantId, input.ResellerId, input.Threshold)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.DashboardLowStockItem
		if err := rows.Scan(&item.ProductId, &item.ProductName, &item.Quantity); err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *dashboardRepository) GetRevenueByReseller(ctx context.Context, input domain.DashboardQueryInput) ([]domain.DashboardResellerSalesItem, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	items := make([]domain.DashboardResellerSalesItem, 0)

	query := `
	SELECT s.user_id, u.name, COALESCE(SUM(si.quantity * si.unit_price), 0) AS revenue
	FROM sales s
	JOIN users u ON u.id = s.user_id AND u.tenant_id = s.tenant_id
	JOIN sales_versions sv ON sv.sales_id = s.id AND sv.version = s.last_version AND sv.tenant_id = s.tenant_id
	JOIN sales_items si ON si.sales_version_id = sv.id AND si.tenant_id = s.tenant_id
	JOIN skus sk ON sk.id = si.sku_id AND sk.tenant_id = s.tenant_id
	WHERE s.tenant_id = $1
	  AND ($2::bigint IS NULL OR s.user_id = $2)
	  AND ($3::bigint IS NULL OR sk.product_id = $3)
	  AND s.date >= $4
	  AND s.date <= $5
	GROUP BY s.user_id, u.name
	ORDER BY revenue DESC, u.name ASC`

	rows, err := r.db.QueryContext(ctx, query, tenantId, input.ResellerId, input.ProductId, input.From, input.To)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.DashboardResellerSalesItem
		if err := rows.Scan(&item.ResellerId, &item.ResellerName, &item.Value); err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *dashboardRepository) GetTopProductsByReseller(ctx context.Context, input domain.DashboardQueryInput, limit int) ([]domain.DashboardProductSalesItem, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	items := make([]domain.DashboardProductSalesItem, 0)

	query := `
	SELECT p.id, p.name, COALESCE(SUM(si.quantity), 0) AS qty
	FROM sales s
	JOIN sales_versions sv ON sv.sales_id = s.id AND sv.version = s.last_version AND sv.tenant_id = s.tenant_id
	JOIN sales_items si ON si.sales_version_id = sv.id AND si.tenant_id = s.tenant_id
	JOIN skus sk ON sk.id = si.sku_id AND sk.tenant_id = s.tenant_id
	JOIN products p ON p.id = sk.product_id AND p.tenant_id = s.tenant_id
	WHERE s.tenant_id = $1
	  AND ($2::bigint IS NULL OR s.user_id = $2)
	  AND ($3::bigint IS NULL OR p.id = $3)
	  AND s.date >= $4
	  AND s.date <= $5
	GROUP BY p.id, p.name
	ORDER BY qty DESC, p.name ASC
	LIMIT $6`

	rows, err := r.db.QueryContext(ctx, query, tenantId, input.ResellerId, input.ProductId, input.From, input.To, limit)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.DashboardProductSalesItem
		if err := rows.Scan(&item.ProductId, &item.ProductName, &item.Quantity); err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}
