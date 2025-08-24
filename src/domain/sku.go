package domain

type Sku struct {
	Id       int64
	Code     string
	Color    string
	Size     string
	Cost     *float64
	Price    *float64
	Quantity float64
	Product  Product
}

func (s *Sku) GetName() string {
	return s.Product.Name + " - " + s.Color + " - " + s.Size
}
