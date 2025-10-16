package domain

import "testing"

func TestSkuGetName(t *testing.T) {
	name := (&Sku{Product: Product{Name: "Shirt"}, Color: "Blue", Size: "M"}).GetName()
	if name != "Shirt - Blue - M" {
		t.Fatalf("unexpected name: %s", name)
	}

	if name := (&Sku{}).GetName(); name != "" {
		t.Fatalf("expected empty name, got %s", name)
	}
}
