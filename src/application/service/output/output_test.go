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
