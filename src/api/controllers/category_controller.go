package controller

import (
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/api/viewmodel"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/labstack/echo/v4"
)

type CategoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) *CategoryController {
	return &CategoryController{categoryService}
}

func (c *CategoryController) Create(context echo.Context) error {
	var categoryRequest request.CreateCategoryRequest
	if err := context.Bind(&categoryRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	err := c.categoryService.Create(context.Request().Context(), categoryRequest)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusCreated, nil)
}

func (c *CategoryController) Edit(context echo.Context) error {
	var categoryRequest request.EditCategoryRequest
	if err := context.Bind(&categoryRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	err := c.categoryService.Edit(context.Request().Context(), categoryRequest)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, nil)
}

func (c *CategoryController) GetById(context echo.Context) error {
	id := helper.ParseInt64(context.Param("id"))

	category, err := c.categoryService.GetById(context.Request().Context(), id)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, viewmodel.ToGetCategoryViewModel(category))
}

func (c *CategoryController) GetAll(context echo.Context) error {
	categories, err := c.categoryService.GetAll(context.Request().Context())
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	var categoryViewModels []viewmodel.GetCategoryViewModel
	for _, category := range categories {
		categoryViewModels = append(categoryViewModels, viewmodel.ToGetCategoryViewModel(category))
	}

	return context.JSON(_http.StatusOK, categoryViewModels)
}

func (c *CategoryController) Inactivate(context echo.Context) error {
	id := helper.ParseInt64(context.Param("id"))

	err := c.categoryService.Inactivate(context.Request().Context(), id)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, nil)
}
