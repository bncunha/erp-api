package output

import "github.com/bncunha/erp-api/src/domain"

type GetSalesOutput struct {
	Sales []GetSalesItemOutput
}

func (o *GetSalesOutput) GetSummary() GetSalesSummaryOutput {
	summary := GetSalesSummaryOutput{}
	totalSum := 0.0
	summary.TotalSales = 0
	summary.ReceivedValue = 0
	summary.FutureRevenue = 0
	summary.AverageTicket = 0

	for _, sale := range o.Sales {
		if sale.Status == domain.PaymentStatusPaid {
			summary.ReceivedValue += sale.TotalValue
		} else {
			summary.FutureRevenue += sale.TotalValue
		}
		summary.TotalItems += sale.TotalItems
		totalSum += sale.TotalValue
	}
	if len(o.Sales) > 0 {
		summary.AverageTicket = totalSum / float64(len(o.Sales))
		summary.TotalSales = float64(len(o.Sales))
	}
	return summary
}

type GetSalesSummaryOutput struct {
	TotalItems    float64
	TotalSales    float64
	ReceivedValue float64
	FutureRevenue float64
	AverageTicket float64
}

type GetSalesItemOutput struct {
	Id           int
	Date         string
	SellerName   string
	CustomerName string
	TotalValue   float64
	TotalItems   float64
	Status       domain.PaymentStatus
}
