package service

import (
	"context"
	"time"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/service/input"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/application/usecase/sales_usecase"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

var (
	ErrPermissionDenied = errors.New("Acesso negado.")
)

type SalesService interface {
	CreateSales(ctx context.Context, request request.CreateSaleRequest) error
	GetSales(ctx context.Context, request request.ListSalesRequest) (output output.GetSalesOutput, err error)
}

type salesService struct {
	salesUsecase    sales_usecase.SalesUseCase
	salesRepository repository.SalesRepository
}

func NewSalesService(salesUsecase sales_usecase.SalesUseCase, salesRepository repository.SalesRepository) SalesService {
	return &salesService{
		salesUsecase,
		salesRepository,
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

func (s *salesService) GetSales(ctx context.Context, request request.ListSalesRequest) (output output.GetSalesOutput, err error) {
	userRole := ctx.Value(constants.ROLE_KEY).(string)
	if userRole == "" {
		return output, ErrPermissionDenied
	}

	var userId *int64
	if userRole == string(domain.UserRoleAdmin) {
		userId = request.UserId
	} else {
		if v, ok := ctx.Value(constants.USERID_KEY).(int64); ok {
			userId = &v
		}
	}

	sales, err := s.salesRepository.GetSales(ctx, input.GetSalesInput{
		InitialDate:   request.MinDate,
		FinalDate:     request.MaxDate,
		UserId:        userId,
		CustomerId:    request.CustomerId,
		PaymentStatus: request.PaymentStatus,
	})
	if err != nil {
		return output, err
	}

	output.Sales = sales
	return output, nil
}
