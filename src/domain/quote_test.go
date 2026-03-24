package domain

import (
	"testing"
	"time"
)

func TestQuoteRecalculateTotals(t *testing.T) {
	quote := Quote{
		ShippingType:          QuoteShippingTypeFree,
		DownPaymentPercentage: 50,
		DiscountPercentage:    10,
		ShippingCost:          4.5,
		Items: []QuoteItem{
			{Quantity: 2, UnitPrice: 10},
			{Quantity: 1, UnitPrice: 5.5},
		},
	}

	if err := quote.RecalculateTotals(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if quote.SubtotalAmount != 25.5 {
		t.Fatalf("expected subtotal 25.5, got %.2f", quote.SubtotalAmount)
	}
	if quote.TotalAmount != 27.45 {
		t.Fatalf("expected total 27.45, got %.2f", quote.TotalAmount)
	}
	if quote.Items[0].TotalPrice != 20 {
		t.Fatalf("expected first item total 20, got %.2f", quote.Items[0].TotalPrice)
	}
}

func TestQuoteRecalculateTotalsWithoutDiscountAndShipping(t *testing.T) {
	quote := Quote{
		ShippingType: QuoteShippingTypeFree,
		Items: []QuoteItem{
			{Quantity: 3, UnitPrice: 7},
		},
	}

	if err := quote.RecalculateTotals(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if quote.SubtotalAmount != 21 {
		t.Fatalf("expected subtotal 21, got %.2f", quote.SubtotalAmount)
	}
	if quote.TotalAmount != 21 {
		t.Fatalf("expected total 21, got %.2f", quote.TotalAmount)
	}
}

func TestQuoteStatusTransition(t *testing.T) {
	if err := ValidateQuoteStatusTransition(QuoteStatusDraft, QuoteStatusSent); err != nil {
		t.Fatalf("expected transition DRAFT->SENT to be valid")
	}
	if err := ValidateQuoteStatusTransition(QuoteStatusSent, QuoteStatusCanceled); err != nil {
		t.Fatalf("expected transition SENT->CANCELED to be valid")
	}
	if err := ValidateQuoteStatusTransition(QuoteStatusApproved, QuoteStatusDraft); err == nil {
		t.Fatalf("expected transition APPROVED->DRAFT to be invalid")
	}
}

func TestEffectiveQuoteStatus(t *testing.T) {
	now := time.Date(2026, 3, 16, 10, 0, 0, 0, time.UTC)
	validUntil := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	status := EffectiveQuoteStatus(QuoteStatusSent, validUntil, now)
	if status != QuoteStatusExpired {
		t.Fatalf("expected expired status, got %s", status)
	}
}
