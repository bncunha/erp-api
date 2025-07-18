package domain

type Sku struct {
	Id    int64
	Code  string
	Color string
	Size  string
	Cost  *float64
	Price *float64
}