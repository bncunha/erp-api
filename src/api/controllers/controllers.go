package controller

import "github.com/bncunha/erp-api/src/application/service"

type Controller struct {
	services *service.ApplicationService
	ProductController *ProductController
	SkuController *SkuController
}

func NewController(services *service.ApplicationService) *Controller {
	return &Controller{
		services: services,
	}
}

func (c *Controller) SetupControllers() {
	c.ProductController = NewProductController(c.services.ProductService)
	c.SkuController = NewSkuController(c.services.SkuService)
}