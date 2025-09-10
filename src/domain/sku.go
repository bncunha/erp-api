package domain

type Sku struct {
	Id       int64
	Code     string
	Color    string
	Size     string
	Cost     *float64
	Price    float64
	Quantity float64
	Product  Product
}

func (s *Sku) GetName() string {
	skuName := ""
	if s.Product.Name != "" {
		skuName = s.Product.Name + " - "
	}
	if s.Color != "" {
		skuName = skuName + s.Color + " - "
	}
	if s.Size != "" {
		skuName = skuName + s.Size
	}
	return skuName
}
