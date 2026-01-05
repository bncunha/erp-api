package domain

import (
	"context"
	"time"
)

type DashboardQueryInput struct {
	From       time.Time
	To         time.Time
	ResellerId *int64
	ProductId  *int64
}

type DashboardStockQueryInput struct {
	ResellerId *int64
	Threshold  float64
}

type DashboardTimeSeriesItem struct {
	Date  time.Time
	Value float64
}

type DashboardResellerSalesItem struct {
	ResellerId   int64
	ResellerName string
	Value        float64
}

type DashboardProductSalesItem struct {
	ProductId   int64
	ProductName string
	Quantity    float64
}

type DashboardLowStockItem struct {
	ProductId   int64
	ProductName string
	Quantity    float64
}

type DashboardRepository interface {
	GetRevenue(ctx context.Context, input DashboardQueryInput) (float64, error)
	GetSalesCount(ctx context.Context, input DashboardQueryInput) (int64, error)
	GetRevenueByDay(ctx context.Context, input DashboardQueryInput) ([]DashboardTimeSeriesItem, error)
	GetSalesCountByDay(ctx context.Context, input DashboardQueryInput) ([]DashboardTimeSeriesItem, error)
	GetStockTotal(ctx context.Context, input DashboardStockQueryInput) (float64, error)
	GetLowStockProducts(ctx context.Context, input DashboardStockQueryInput) ([]DashboardLowStockItem, error)
	GetRevenueByReseller(ctx context.Context, input DashboardQueryInput) ([]DashboardResellerSalesItem, error)
	GetTopProductsByReseller(ctx context.Context, input DashboardQueryInput, limit int) ([]DashboardProductSalesItem, error)
}
