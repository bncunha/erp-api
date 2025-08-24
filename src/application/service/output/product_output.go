package output

import "github.com/bncunha/erp-api/src/domain"

type GetAllProductsOutput struct {
	Product  domain.Product
	Quantity float64
}
