package controller

import (
	"errors"
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/labstack/echo/v4"
)

type ProductController struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (c *ProductController) Create(context echo.Context) error {
	var productRequest request.CreateProductRequest
	if err := context.Bind(&productRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(errors.New("parametros invalidos")))
	}

	err := c.productService.Create(context.Request().Context(), productRequest)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusCreated, nil)
}