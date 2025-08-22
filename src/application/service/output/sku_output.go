package output

import "github.com/bncunha/erp-api/src/domain"

type GetAllSkusByProductOutput struct {
	Sku      domain.Sku
	Quantity float64
}
