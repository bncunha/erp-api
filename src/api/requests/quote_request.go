package request

import (
	"strings"
	"time"

	"github.com/bncunha/erp-api/src/application/validator"
	"github.com/bncunha/erp-api/src/domain"
)

type UpsertQuoteItemRequest struct {
	SkuID     int64   `json:"sku_id" validate:"required,gt=0"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" validate:"required,gte=0"`
}

type CreateQuoteRequest struct {
	CustomerID            int64                    `json:"customer_id" validate:"required,gt=0"`
	ValidUntil            string                   `json:"valid_until" validate:"required"`
	DownPaymentPercentage float64                  `json:"down_payment_percentage" validate:"required,gte=0,lte=100"`
	DiscountPercentage    float64                  `json:"discount_percentage" validate:"gte=0,lte=100"`
	ShippingType          domain.QuoteShippingType `json:"shipping_type" validate:"required,oneof=FREE FREE_REGION FREE_MIN_VALUE TO_CALCULATE"`
	ShippingCost          float64                  `json:"shipping_cost" validate:"gte=0"`
	ShippingRegion        *string                  `json:"shipping_region"`
	ShippingMinValue      *float64                 `json:"shipping_min_value"`
	Notes                 *string                  `json:"notes"`
	Items                 []UpsertQuoteItemRequest `json:"items" validate:"required,min=1"`
}

type UpdateQuoteRequest struct {
	CustomerID            int64                    `json:"customer_id" validate:"required,gt=0"`
	ValidUntil            string                   `json:"valid_until" validate:"required"`
	DownPaymentPercentage float64                  `json:"down_payment_percentage" validate:"required,gte=0,lte=100"`
	DiscountPercentage    float64                  `json:"discount_percentage" validate:"gte=0,lte=100"`
	ShippingType          domain.QuoteShippingType `json:"shipping_type" validate:"required,oneof=FREE FREE_REGION FREE_MIN_VALUE TO_CALCULATE"`
	ShippingCost          float64                  `json:"shipping_cost" validate:"gte=0"`
	ShippingRegion        *string                  `json:"shipping_region"`
	ShippingMinValue      *float64                 `json:"shipping_min_value"`
	Notes                 *string                  `json:"notes"`
	Items                 []UpsertQuoteItemRequest `json:"items" validate:"required,min=1"`
}

type PatchQuoteStatusRequest struct {
	Status domain.QuoteStatus `json:"status" validate:"required,oneof=DRAFT SENT APPROVED REJECTED EXPIRED CANCELED"`
}

func (r *CreateQuoteRequest) Validate() error {
	if err := validator.Validate(r); err != nil {
		return err
	}
	return validateQuoteShipping(r.ShippingType, r.ShippingRegion, r.ShippingMinValue)
}

func (r *UpdateQuoteRequest) Validate() error {
	if err := validator.Validate(r); err != nil {
		return err
	}
	return validateQuoteShipping(r.ShippingType, r.ShippingRegion, r.ShippingMinValue)
}

func (r *PatchQuoteStatusRequest) Validate() error {
	return validator.Validate(r)
}

func (r *CreateQuoteRequest) ParseValidUntil() (time.Time, error) {
	return parseDateOnlyOrRFC3339(r.ValidUntil)
}

func (r *UpdateQuoteRequest) ParseValidUntil() (time.Time, error) {
	return parseDateOnlyOrRFC3339(r.ValidUntil)
}

func validateQuoteShipping(shippingType domain.QuoteShippingType, shippingRegion *string, shippingMinValue *float64) error {
	if shippingType == domain.QuoteShippingTypeFreeRegion {
		if shippingRegion == nil || strings.TrimSpace(*shippingRegion) == "" {
			return domain.ErrQuoteShippingRegionRequired
		}
	}
	if shippingType == domain.QuoteShippingTypeFreeMin {
		if shippingMinValue == nil || *shippingMinValue <= 0 {
			return domain.ErrQuoteShippingMinValueInvalid
		}
	}
	return nil
}

func parseDateOnlyOrRFC3339(value string) (time.Time, error) {
	parsed, err := time.Parse(time.DateOnly, strings.TrimSpace(value))
	if err == nil {
		return parsed, nil
	}
	return time.Parse(time.RFC3339, strings.TrimSpace(value))
}
