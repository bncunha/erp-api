package domain

import "testing"

func TestInventoryConstants(t *testing.T) {
	if InventoryTypePrimary != "PRIMARY" || InventoryTypeReseller != "RESELLER" {
		t.Fatalf("unexpected inventory type constants")
	}
	if InventoryTransactionTypeTransfer != "TRANSFER" || InventoryTransactionTypeIn != "IN" || InventoryTransactionTypeOut != "OUT" {
		t.Fatalf("unexpected transaction type constants")
	}
}

func TestInventoryStructInitialization(t *testing.T) {
	inv := Inventory{Id: 1, TenantId: 2, User: User{Id: 3}, Type: InventoryTypePrimary}
	if inv.Id != 1 || inv.User.Id != 3 {
		t.Fatalf("unexpected inventory values: %+v", inv)
	}
}
