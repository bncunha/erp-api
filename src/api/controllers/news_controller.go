package controller

import (
	"errors"
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	"github.com/bncunha/erp-api/src/api/viewmodel"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/labstack/echo/v4"
)

type NewsController struct {
	newsService service.NewsService
}

func NewNewsController(newsService service.NewsService) *NewsController {
	return &NewsController{newsService: newsService}
}

func (c *NewsController) GetLatest(context echo.Context) error {
	news, err := c.newsService.GetLatest(context.Request().Context())
	if err != nil {
		if errors.Is(err, domain.ErrNewsNotFound) {
			return context.NoContent(_http.StatusNoContent)
		}
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, viewmodel.ToNewsViewModel(news))
}
