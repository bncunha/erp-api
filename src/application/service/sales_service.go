package service

import (
	"context"
	"time"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/usecase/sales_usecase"
)

type SalesService interface {
	CreateSales(ctx context.Context, request request.CreateSaleRequest) error
}

type salesService struct {
	salesUsecase sales_usecase.SalesUseCase
}

func NewSalesService(salesUsecase sales_usecase.SalesUseCase) SalesService {
	return &salesService{
		salesUsecase,
	}
}

func (s *salesService) CreateSales(ctx context.Context, request request.CreateSaleRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	userId := int64(ctx.Value(constants.USERID_KEY).(float64))

	items := make([]sales_usecase.DoSaleItemsInput, 0)
	payments := make([]sales_usecase.DoSalePaymentsInput, 0)
	for _, item := range request.Items {
		items = append(items, sales_usecase.DoSaleItemsInput{
			SkuId:    item.SkuId,
			Quantity: item.Quantity,
		})
	}

	for _, payment := range request.Payments {
		dates := make([]sales_usecase.DoSalePaymentDatesInput, 0)
		for j, date := range payment.Dates {
			dates = append(dates, sales_usecase.DoSalePaymentDatesInput{
				DueDate:           date.Date,
				InstallmentNumber: j + 1,
				InstallmentValue:  date.InstallmentValue,
			})
		}
		payments = append(payments, sales_usecase.DoSalePaymentsInput{
			PaymentType: payment.PaymentType,
			Dates:       dates,
		})
	}

	input := sales_usecase.DoSaleInput{
		CustomerId: request.CustomerId,
		UserId:     userId,
		Date:       time.Now(),
		Items:      items,
		Payments:   payments,
	}

	return s.salesUsecase.DoSale(ctx, input)
}
