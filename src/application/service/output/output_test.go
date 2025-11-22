package output

import (
	"testing"

	"github.com/bncunha/erp-api/src/domain"
)

func TestGetAllProductsOutput(t *testing.T) {
	out := GetAllProductsOutput{Product: domain.Product{Name: "Test"}, Quantity: 2}
	if out.Product.Name != "Test" || out.Quantity != 2 {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestGetInventoryItemsOutput(t *testing.T) {
	out := GetInventoryItemsOutput{InventoryItemId: 1, Quantity: 2}
	if out.InventoryItemId != 1 || out.Quantity != 2 {
		t.Fatalf("unexpected inventory item output: %+v", out)
	}
}

func TestGetInventoryTransactionsOutput(t *testing.T) {
	out := GetInventoryTransactionsOutput{Id: 1, Quantity: 2}
	if out.Id != 1 || out.Quantity != 2 {
		t.Fatalf("unexpected inventory transaction output: %+v", out)
	}
}

func TestLoginOutput(t *testing.T) {
	out := LoginOutput{Name: "User", Token: "token"}
	if out.Name != "User" || out.Token != "token" {
		t.Fatalf("unexpected login output: %+v", out)
	}
}

func TestGetSalesOutputGetSummary(t *testing.T) {
	sales := []GetSalesItemOutput{
		{ReceivedValue: 10, FutureRevenue: 5, TotalItems: 2, TotalValue: 20},
		{ReceivedValue: 5, FutureRevenue: 10, TotalItems: 1, TotalValue: 10},
	}
	out := GetSalesOutput{Sales: sales}
	summary := out.GetSummary()

	if summary.TotalSales != 2 || summary.AverageTicket != 15 {
		t.Fatalf("unexpected summary calculation: %+v", summary)
	}
	if summary.ReceivedValue != 15 || summary.FutureRevenue != 15 || summary.TotalItems != 3 {
		t.Fatalf("unexpected aggregate values: %+v", summary)
	}
}

func TestGetSalesOutputGetSummaryEmpty(t *testing.T) {
	out := GetSalesOutput{}
	summary := out.GetSummary()
	if summary.TotalSales != 0 || summary.AverageTicket != 0 || summary.TotalItems != 0 {
		t.Fatalf("expected zeroed summary, got %+v", summary)
	}
}
