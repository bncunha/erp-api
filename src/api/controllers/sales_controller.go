package controller

import (
	_http "net/http"
	"time"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/api/viewmodel"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/labstack/echo/v4"
)

type SalesController struct {
	salesService service.SalesService
}

func NewSalesController(salesService service.SalesService) *SalesController {
	return &SalesController{salesService}
}

func (c *SalesController) Create(context echo.Context) error {
	var salesRequeste request.CreateSaleRequest
	if err := context.Bind(&salesRequeste); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	err := c.salesService.CreateSales(context.Request().Context(), salesRequeste)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusCreated, nil)
}

func (c *SalesController) GetAll(context echo.Context) error {
	var customerId, userId *int64
	var minDate, maxDate *time.Time
	var paymentStatus *domain.PaymentStatus
	if context.QueryParam("customer_id") != "" {
		id := helper.ParseInt64(context.QueryParam("customer_id"))
		customerId = &id
	}
	if context.QueryParam("min_date") != "" {
		date, _ := time.Parse(time.DateOnly, context.QueryParam("min_date"))
		minDate = &date
	}
	if context.QueryParam("max_date") != "" {
		date, _ := time.Parse(time.DateOnly, context.QueryParam("max_date"))
		maxDate = &date
	}
	if context.QueryParam("user_id") != "" {
		id := helper.ParseInt64(context.QueryParam("user_id"))
		userId = &id
	}
	if context.QueryParam("payment_status") != "" {
		status := domain.PaymentStatus(context.QueryParam("payment_status"))
		paymentStatus = &status
	}

	sales, err := c.salesService.GetSales(context.Request().Context(), request.ListSalesRequest{
		CustomerId:    customerId,
		MinDate:       minDate,
		MaxDate:       maxDate,
		UserId:        userId,
		PaymentStatus: paymentStatus,
	})
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, viewmodel.ToSalesViewModel(sales))
}

func (c *SalesController) GetById(context echo.Context) error {
	id := helper.ParseInt64(context.Param("id"))

	saleOutput, paymentGroupOutput, itemsOutput, err := c.salesService.GetById(context.Request().Context(), id)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, viewmodel.ToSaleByIdViewModel(saleOutput, paymentGroupOutput, itemsOutput))
}
