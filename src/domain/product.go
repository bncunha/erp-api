package domain

type Product struct {
	Name        string
	Description string
	Category    Category
	Skus        []Sku
}