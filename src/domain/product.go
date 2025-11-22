package domain

import "time"

type Product struct {
	Id          int64
	Name        string
	Description string
	Category    Category
	Skus        []Sku
	DeletedAt   *time.Time
}
