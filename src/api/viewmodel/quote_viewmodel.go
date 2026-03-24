package viewmodel

import (
	"math"
	"time"

	"github.com/bncunha/erp-api/src/domain"
)

type QuoteCustomerViewModel struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type QuoteCompanyViewModel struct {
	Name      string  `json:"name"`
	LegalName string  `json:"legal_name"`
	CNPJ      string  `json:"cnpj"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Address   string  `json:"address"`
	LogoURL   *string `json:"logo_url"`
}

type QuoteItemViewModel struct {
	ID                 int64   `json:"id"`
	SkuID              int64   `json:"sku_id"`
	ProductID          int64   `json:"product_id"`
	SkuCode            string  `json:"sku_code"`
	ProductDescription string  `json:"product_description"`
	Quantity           float64 `json:"quantity"`
	UnitPrice          float64 `json:"unit_price"`
	TotalPrice         float64 `json:"total_price"`
}

type QuoteDetailedViewModel struct {
	ID                    int64                  `json:"id"`
	QuoteNumber           string                 `json:"quote_number"`
	Company               QuoteCompanyViewModel  `json:"company"`
	Customer              QuoteCustomerViewModel `json:"customer"`
	Status                string                 `json:"status"`
	ValidUntil            string                 `json:"valid_until"`
	DownPaymentPercentage float64                `json:"down_payment_percentage"`
	DiscountPercentage    float64                `json:"discount_percentage"`
	DiscountAmount        float64                `json:"discount_amount"`
	ShippingCost          float64                `json:"shipping_cost"`
	DownPaymentAmount     float64                `json:"down_payment_amount"`
	RemainingAmount       float64                `json:"remaining_amount"`
	ShippingType          string                 `json:"shipping_type"`
	ShippingRegion        *string                `json:"shipping_region"`
	ShippingMinValue      *float64               `json:"shipping_min_value"`
	ShippingDescription   string                 `json:"shipping_description"`
	Notes                 *string                `json:"notes"`
	Items                 []QuoteItemViewModel   `json:"items"`
	SubtotalAmount        float64                `json:"subtotal_amount"`
	TotalAmount           float64                `json:"total_amount"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
}

type QuoteListItemViewModel struct {
	ID                    int64                  `json:"id"`
	QuoteNumber           string                 `json:"quote_number"`
	Customer              QuoteCustomerViewModel `json:"customer"`
	Status                string                 `json:"status"`
	ValidUntil            string                 `json:"valid_until"`
	TotalAmount           float64                `json:"total_amount"`
	DownPaymentPercentage float64                `json:"down_payment_percentage"`
	DownPaymentAmount     float64                `json:"down_payment_amount"`
	ShippingDescription   string                 `json:"shipping_description"`
	CreatedAt             time.Time              `json:"created_at"`
}

type QuoteListViewModel struct {
	Items      []QuoteListItemViewModel `json:"items"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalPages int                      `json:"total_pages"`
}

type QuoteDuplicateViewModel struct {
	ID int64 `json:"id"`
}

func ToQuoteDetailedViewModel(quote domain.Quote) QuoteDetailedViewModel {
	items := make([]QuoteItemViewModel, 0, len(quote.Items))
	for _, item := range quote.Items {
		items = append(items, QuoteItemViewModel{
			ID:                 item.ID,
			SkuID:              item.SkuID,
			ProductID:          item.ProductID,
			SkuCode:            item.SkuCode,
			ProductDescription: item.ProductDescriptionSnapshot,
			Quantity:           item.Quantity,
			UnitPrice:          item.UnitPrice,
			TotalPrice:         item.TotalPrice,
		})
	}
	discountAmount := roundMoney(quote.SubtotalAmount * quote.DiscountPercentage / 100)
	downPaymentAmount := roundMoney(quote.TotalAmount * quote.DownPaymentPercentage / 100)
	remainingAmount := roundMoney(quote.TotalAmount - downPaymentAmount)
	return QuoteDetailedViewModel{
		ID:          quote.ID,
		QuoteNumber: quote.QuoteNumber,
		Company: QuoteCompanyViewModel{
			Name:      quote.Company.Name,
			LegalName: quote.Company.LegalName,
			CNPJ:      quote.Company.CNPJ,
			Email:     quote.Company.Email,
			Phone:     quote.Company.Phone,
			Address:   quote.Company.Address,
			LogoURL:   quote.Company.LogoURL,
		},
		Customer: QuoteCustomerViewModel{
			ID:    quote.Customer.ID,
			Name:  quote.Customer.Name,
			Phone: quote.Customer.Phone,
		},
		Status:                string(quote.Status),
		ValidUntil:            quote.ValidUntil.Format(time.DateOnly),
		DownPaymentPercentage: quote.DownPaymentPercentage,
		DiscountPercentage:    quote.DiscountPercentage,
		DiscountAmount:        discountAmount,
		ShippingCost:          quote.ShippingCost,
		DownPaymentAmount:     downPaymentAmount,
		RemainingAmount:       remainingAmount,
		ShippingType:          string(quote.ShippingType),
		ShippingRegion:        quote.ShippingRegion,
		ShippingMinValue:      quote.ShippingMinValue,
		ShippingDescription:   quote.ShippingDescription,
		Notes:                 quote.Notes,
		Items:                 items,
		SubtotalAmount:        quote.SubtotalAmount,
		TotalAmount:           quote.TotalAmount,
		CreatedAt:             quote.CreatedAt,
		UpdatedAt:             quote.UpdatedAt,
	}
}

func ToQuoteListViewModel(output domain.GetQuotesOutput) QuoteListViewModel {
	items := make([]QuoteListItemViewModel, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, QuoteListItemViewModel{
			ID:          item.ID,
			QuoteNumber: item.QuoteNumber,
			Customer: QuoteCustomerViewModel{
				ID:    item.Customer.ID,
				Name:  item.Customer.Name,
				Phone: item.Customer.Phone,
			},
			Status:                string(item.Status),
			ValidUntil:            item.ValidUntil.Format(time.DateOnly),
			TotalAmount:           item.TotalAmount,
			DownPaymentPercentage: item.DownPaymentPercentage,
			DownPaymentAmount:     roundMoney(item.TotalAmount * item.DownPaymentPercentage / 100),
			ShippingDescription:   item.ShippingDescription,
			CreatedAt:             item.CreatedAt,
		})
	}
	return QuoteListViewModel{
		Items:      items,
		Total:      output.Total,
		Page:       output.Page,
		PageSize:   output.PageSize,
		TotalPages: output.TotalPages,
	}
}

func ToQuoteDuplicateViewModel(id int64) QuoteDuplicateViewModel {
	return QuoteDuplicateViewModel{ID: id}
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}
