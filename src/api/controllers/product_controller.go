package controller

import (
	"errors"
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/api/viewmodel"
	helper "github.com/bncunha/erp-api/src/application/helpers"
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

func (c *ProductController) Edit(context echo.Context) error {
	var productRequest request.EditProductRequest
	if err := context.Bind(&productRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(errors.New("parametros invalidos")))
	}

	productRequest.Id =  helper.ParseInt64(context.Param("id"))
	err := c.productService.Edit(context.Request().Context(), productRequest)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, nil)
}

func (c *ProductController) GetById(context echo.Context) error {
	id := helper.ParseInt64(context.Param("id"))

	product, err := c.productService.GetById(context.Request().Context(), id)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, viewmodel.ToGetProductViewModel(product))
}