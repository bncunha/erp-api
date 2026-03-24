package domain

import (
	"context"
	"database/sql"
	"time"
)

type GetQuotesInput struct {
	Page            int
	PageSize        int
	SortBy          string
	SortOrder       string
	Statuses        []QuoteStatus
	CustomerID      *int64
	ValidUntilStart *time.Time
	ValidUntilEnd   *time.Time
	CreatedAtStart  *time.Time
	CreatedAtEnd    *time.Time
	Search          *string
	OnlyExpired     bool
}

type QuoteListItemOutput struct {
	ID                    int64
	QuoteNumber           string
	Customer              QuoteCustomerSummary
	Status                QuoteStatus
	ValidUntil            time.Time
	TotalAmount           float64
	DownPaymentPercentage float64
	DownPaymentAmount     float64
	ShippingDescription   string
	CreatedAt             time.Time
}

type GetQuotesOutput struct {
	Items      []QuoteListItemOutput
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type QuoteRepository interface {
	Create(ctx context.Context, tx *sql.Tx, quote *Quote) error
	Update(ctx context.Context, tx *sql.Tx, quote *Quote) error
	UpdateStatus(ctx context.Context, tx *sql.Tx, id int64, status QuoteStatus) error
	GetByID(ctx context.Context, id int64) (Quote, error)
	List(ctx context.Context, input GetQuotesInput) (GetQuotesOutput, error)
	GetCompanyProfile(ctx context.Context) (QuoteCompanyProfile, error)
}
