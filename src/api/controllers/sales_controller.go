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
	var customerId, userId []int64
	var minDate, maxDate *time.Time
	var paymentStatus *domain.PaymentStatus

	customersParams := context.QueryParams()["customer_id"]
	usersParams := context.QueryParams()["user_id"]
	if len(customersParams) > 0 {
		for _, customerParam := range customersParams {
			customerId = append(customerId, helper.ParseInt64(customerParam))
		}
	}
	if context.QueryParam("min_date") != "" {
		date, _ := time.Parse(time.DateOnly, context.QueryParam("min_date"))
		minDate = &date
	}
	if context.QueryParam("max_date") != "" {
		date, _ := time.Parse(time.DateOnly, context.QueryParam("max_date"))
		maxDate = &date
	}
	if len(usersParams) > 0 {
		for _, userParam := range usersParams {
			userId = append(userId, helper.ParseInt64(userParam))
		}
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

func (c *SalesController) ChangePaymentStatus(context echo.Context) error {
	id := helper.ParseInt64(context.Param("id"))
	paymentId := helper.ParseInt64(context.Param("payment_id"))
	var request request.ChangePaymentStatusRequest
	if err := context.Bind(&request); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	err := c.salesService.ChangePaymentStatus(context.Request().Context(), id, paymentId, request)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, nil)
}
