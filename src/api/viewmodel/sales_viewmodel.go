package viewmodel

import "github.com/bncunha/erp-api/src/application/service/output"

type SalesViewModel struct {
	Summary SalesSummaryViewModel `json:"summary"`
	Sales   []SalesItemViewModel  `json:"sales"`
}

type SalesSummaryViewModel struct {
	TotalSales    float64 `json:"total_sales"`
	TotalItems    float64 `json:"total_items"`
	ReceivedValue float64 `json:"received_value"`
	FutureRevenue float64 `json:"future_revenue"`
	AverageTicket float64 `json:"average_ticket"`
}

type SalesItemViewModel struct {
	Id           int     `json:"id"`
	Date         string  `json:"date"`
	SellerName   string  `json:"seller_name"`
	CustomerName string  `json:"customer_name"`
	TotalValue   float64 `json:"total_value"`
	TotalItems   float64 `json:"total_items"`
	Status       string  `json:"status"`
}

func ToSalesViewModel(output output.GetSalesOutput) SalesViewModel {
	salesViewModel := SalesViewModel{
		Summary: toSalesSummaryViewModel(output.GetSummary()),
		Sales:   make([]SalesItemViewModel, len(output.Sales)),
	}
	for i, sale := range output.Sales {
		salesViewModel.Sales[i] = toSalesItemViewModel(sale)
	}
	return salesViewModel
}

func toSalesSummaryViewModel(summary output.GetSalesSummaryOutput) SalesSummaryViewModel {
	return SalesSummaryViewModel{
		TotalSales:    summary.TotalSales,
		TotalItems:    summary.TotalItems,
		ReceivedValue: summary.ReceivedValue,
		FutureRevenue: summary.FutureRevenue,
		AverageTicket: summary.AverageTicket,
	}
}

func toSalesItemViewModel(sale output.GetSalesItemOutput) SalesItemViewModel {
	return SalesItemViewModel{
		Id:           sale.Id,
		Date:         sale.Date,
		SellerName:   sale.SellerName,
		CustomerName: sale.CustomerName,
		TotalValue:   sale.TotalValue,
		TotalItems:   sale.TotalItems,
		Status:       string(sale.Status),
	}
}
