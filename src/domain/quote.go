package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type QuoteStatus string

type QuoteShippingType string

const (
	QuoteStatusDraft    QuoteStatus = "DRAFT"
	QuoteStatusSent     QuoteStatus = "SENT"
	QuoteStatusApproved QuoteStatus = "APPROVED"
	QuoteStatusRejected QuoteStatus = "REJECTED"
	QuoteStatusExpired  QuoteStatus = "EXPIRED"
	QuoteStatusCanceled QuoteStatus = "CANCELED"
)

const (
	QuoteShippingTypeFree       QuoteShippingType = "FREE"
	QuoteShippingTypeFreeRegion QuoteShippingType = "FREE_REGION"
	QuoteShippingTypeFreeMin    QuoteShippingType = "FREE_MIN_VALUE"
	QuoteShippingTypeCalculate  QuoteShippingType = "TO_CALCULATE"
)

var (
	ErrQuoteItemsRequired           = errors.New("o orçamento deve possuir ao menos um item")
	ErrQuoteInvalidTransition       = errors.New("transição de status inválida")
	ErrQuoteShippingRegionRequired  = errors.New("região do frete é obrigatória para FREE_REGION")
	ErrQuoteShippingMinValueInvalid = errors.New("valor mínimo do frete deve ser maior que zero para FREE_MIN_VALUE")
)

type Quote struct {
	ID                    int64
	TenantID              int64
	CustomerID            int64
	QuoteNumber           string
	Status                QuoteStatus
	ValidUntil            time.Time
	DownPaymentPercentage float64
	DiscountPercentage    float64
	Notes                 *string
	ShippingType          QuoteShippingType
	ShippingCost          float64
	ShippingRegion        *string
	ShippingMinValue      *float64
	ShippingDescription   string
	SubtotalAmount        float64
	TotalAmount           float64
	CreatedByUserID       int64
	ProductionOrderID     *int64
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Items                 []QuoteItem
	Customer              QuoteCustomerSummary
	Company               QuoteCompanyProfile
}

type QuoteItem struct {
	ID                         int64
	QuoteID                    int64
	SkuID                      int64
	ProductID                  int64
	SkuCode                    string
	ProductDescriptionSnapshot string
	Quantity                   float64
	UnitPrice                  float64
	TotalPrice                 float64
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

type QuoteCustomerSummary struct {
	ID    int64
	Name  string
	Phone string
}

type QuoteCompanyProfile struct {
	Name      string
	LegalName string
	CNPJ      string
	Email     string
	Phone     string
	Address   string
	LogoURL   *string
}

func (q *Quote) RecalculateTotals() error {
	if len(q.Items) == 0 {
		return ErrQuoteItemsRequired
	}
	q.SubtotalAmount = 0
	for i := range q.Items {
		q.Items[i].TotalPrice = q.Items[i].Quantity * q.Items[i].UnitPrice
		q.SubtotalAmount += q.Items[i].TotalPrice
	}
	discountAmount := q.SubtotalAmount * q.DiscountPercentage / 100
	q.TotalAmount = q.SubtotalAmount - discountAmount + q.ShippingCost
	q.ShippingDescription = BuildQuoteShippingDescription(q.ShippingType, q.ShippingRegion, q.ShippingMinValue)
	return nil
}

func BuildQuoteShippingDescription(shippingType QuoteShippingType, shippingRegion *string, shippingMinValue *float64) string {
	switch shippingType {
	case QuoteShippingTypeFree:
		return "Frete grátis"
	case QuoteShippingTypeFreeRegion:
		if shippingRegion == nil || strings.TrimSpace(*shippingRegion) == "" {
			return "Frete grátis para a região"
		}
		return fmt.Sprintf("Frete grátis para a região %s", strings.TrimSpace(*shippingRegion))
	case QuoteShippingTypeFreeMin:
		if shippingMinValue == nil {
			return "Frete grátis para pedidos a partir de R$ 0,00"
		}
		return fmt.Sprintf("Frete grátis para pedidos a partir de R$ %.2f", *shippingMinValue)
	default:
		return "Frete a calcular"
	}
}

func IsQuoteStatusFinal(status QuoteStatus) bool {
	return status == QuoteStatusApproved || status == QuoteStatusRejected || status == QuoteStatusCanceled
}

func EffectiveQuoteStatus(status QuoteStatus, validUntil time.Time, now time.Time) QuoteStatus {
	if IsQuoteStatusFinal(status) {
		return status
	}
	validDate := time.Date(validUntil.Year(), validUntil.Month(), validUntil.Day(), 0, 0, 0, 0, now.Location())
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if validDate.Before(today) {
		return QuoteStatusExpired
	}
	return status
}

func ValidateQuoteStatusTransition(from QuoteStatus, to QuoteStatus) error {
	if from == to {
		return nil
	}
	allowed := map[QuoteStatus]map[QuoteStatus]bool{
		QuoteStatusDraft: {
			QuoteStatusSent:     true,
			QuoteStatusCanceled: true,
		},
		QuoteStatusSent: {
			QuoteStatusApproved: true,
			QuoteStatusRejected: true,
			QuoteStatusCanceled: true,
		},
	}
	if nextMap, ok := allowed[from]; ok && nextMap[to] {
		return nil
	}
	return ErrQuoteInvalidTransition
}
