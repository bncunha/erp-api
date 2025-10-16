package domain

import "testing"

func TestStructInitialization(t *testing.T) {
	user := User{Id: 1, Username: "john"}
	if user.Id != 1 || user.Username != "john" {
		t.Fatalf("unexpected user: %+v", user)
	}

	company := Company{Name: "Acme"}
	if company.Name != "Acme" {
		t.Fatalf("unexpected company: %+v", company)
	}

	category := Category{Id: 2, Name: "Clothes"}
	if category.Name != "Clothes" {
		t.Fatalf("unexpected category: %+v", category)
	}

	product := Product{Id: 3, Name: "T-Shirt"}
	if product.Name != "T-Shirt" {
		t.Fatalf("unexpected product: %+v", product)
	}
}
