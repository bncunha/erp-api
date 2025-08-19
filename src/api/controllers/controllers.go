package controller

import "github.com/bncunha/erp-api/src/application/service"

type Controller struct {
	services           *service.ApplicationService
	ProductController  *ProductController
	SkuController      *SkuController
	CategoryController *CategoryController
	AuthController     *AuthController
	UserController     *UserController
}

func NewController(services *service.ApplicationService) *Controller {
	return &Controller{
		services: services,
	}
}

func (c *Controller) SetupControllers() {
	c.ProductController = NewProductController(c.services.ProductService)
	c.SkuController = NewSkuController(c.services.SkuService)
	c.CategoryController = NewCategoryController(c.services.CategoryService)
	c.AuthController = NewAuthController(c.services.AuthService)
	c.UserController = NewUserController(c.services.UserService)
}
