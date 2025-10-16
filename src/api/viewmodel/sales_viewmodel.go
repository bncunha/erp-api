package viewmodel

import (
	"time"

	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

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

type SaleByIdViewModel struct {
	Id            int                     `json:"id"`
	Code          string                  `json:"code"`
	Date          time.Time               `json:"date"`
	TotalValue    float64                 `json:"total_value"`
	SellerName    string                  `json:"seller_name"`
	CustomerName  string                  `json:"customer_name"`
	ReceivedValue float64                 `json:"received_value"`
	FutureRevenue float64                 `json:"future_revenue"`
	PaymentStatus domain.PaymentStatus    `json:"payment_status"`
	Payments      []SalePaymentsViewModel `json:"payments"`
	Items         []SaleItemsViewModel    `json:"items"`
}

type SalePaymentsViewModel struct {
	PaymentType  domain.PaymentType          `json:"payment_type"`
	Installments []SalePaymentsItemViewModel `json:"installments"`
}

type SalePaymentsItemViewModel struct {
	Id                int64                `json:"id"`
	InstallmentNumber int64                `json:"installment_number"`
	InstallmentValue  float64              `json:"installment_value"`
	DueDate           string               `json:"due_date"`
	PaidDate          *string              `json:"paid_date"`
	PaymentStatus     domain.PaymentStatus `json:"payment_status"`
	PaymentType       domain.PaymentType   `json:"payment_type"`
}

type SaleItemsViewModel struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalValue  float64 `json:"total_value"`
}

func ToSaleByIdViewModel(sale output.GetSaleByIdOutput, paymentGroupOutput []output.GetSalesPaymentGroupOutput, itemsOutput []output.GetItemsOutput) SaleByIdViewModel {
	return SaleByIdViewModel{
		Id:            sale.Id,
		Code:          sale.Code,
		Date:          sale.Date,
		TotalValue:    sale.TotalValue,
		SellerName:    sale.SellerName,
		CustomerName:  sale.CustomerName,
		ReceivedValue: sale.ReceivedValue,
		FutureRevenue: sale.FutureRevenue,
		PaymentStatus: sale.PaymentStatus,
		Payments:      toSalePaymentsViewModel(paymentGroupOutput),
		Items:         toSaleItemsViewModel(itemsOutput),
	}
}

func toSalePaymentsViewModel(paymentGroupOutput []output.GetSalesPaymentGroupOutput) []SalePaymentsViewModel {
	paymentsViewModel := make([]SalePaymentsViewModel, len(paymentGroupOutput))
	for i, paymentGroup := range paymentGroupOutput {
		paymentsViewModel[i] = SalePaymentsViewModel{
			PaymentType:  paymentGroup.PaymentType,
			Installments: toSalePaymentsItemViewModel(paymentGroup.Installments),
		}
	}
	return paymentsViewModel
}

func toSalePaymentsItemViewModel(payments []output.GetSalesPaymentOutput) []SalePaymentsItemViewModel {
	paymentsViewModel := make([]SalePaymentsItemViewModel, len(payments))

	for i, payment := range payments {
		var paidDate string
		if payment.PaidDate != nil {
			paidDate = payment.PaidDate.Format(time.DateOnly)
		}
		paymentsViewModel[i] = SalePaymentsItemViewModel{
			Id:                payment.Id,
			InstallmentNumber: payment.InstallmentNumber,
			InstallmentValue:  payment.InstallmentValue,
			DueDate:           payment.DueDate.Format(time.DateOnly),
			PaidDate:          &paidDate,
			PaymentStatus:     payment.PaymentStatus,
			PaymentType:       payment.PaymentType,
		}
	}
	return paymentsViewModel
}

func toSaleItemsViewModel(itemsOutput []output.GetItemsOutput) []SaleItemsViewModel {
	itemsViewModel := make([]SaleItemsViewModel, len(itemsOutput))
	for i, item := range itemsOutput {
		itemsViewModel[i] = SaleItemsViewModel{
			Code:        item.Sku.Code,
			Description: item.Sku.GetName(),
			Quantity:    item.Quantity,
			UnitPrice:   item.Sku.Price,
			TotalValue:  item.Sku.Price * item.Quantity,
		}
	}
	return itemsViewModel
}
