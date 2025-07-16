package controller

import (
	"errors"
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	product_service "github.com/bncunha/erp-api/src/application/service/product"
	"github.com/bncunha/erp-api/src/application/service/product/input"
	"github.com/labstack/echo/v4"
)

type ProductController struct {
	productService product_service.ProductService
}

func NewProductController(productService product_service.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (c *ProductController) Create(context echo.Context) error {
	var productRequest request.CreateProductRequest
	if err := context.Bind(&productRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(errors.New("parametros invalidos")))
	}

	err := c.productService.Create(context.Request().Context(), input.CreateProductInput{
		Name: productRequest.Name,
		Description: productRequest.Description,
	})
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusCreated, nil)
}