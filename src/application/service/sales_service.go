package service

import (
	"context"
	"math"
	"time"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/service/input"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/application/usecase/sales_usecase"
	"github.com/bncunha/erp-api/src/domain"
)

var (
	ErrPermissionDenied = errors.New("Acesso negado.")
)

type SalesService interface {
	CreateSales(ctx context.Context, request request.CreateSaleRequest) error
	GetSales(ctx context.Context, request request.ListSalesRequest) (output output.GetSalesOutput, err error)
	GetById(ctx context.Context, id int64) (saleOutput output.GetSaleByIdOutput, paymentGroupOutput []output.GetSalesPaymentGroupOutput, itemsOutput []output.GetItemsOutput, err error)
	ChangePaymentStatus(ctx context.Context, id int64, paymentId int64, request request.ChangePaymentStatusRequest) error
}

type salesService struct {
	salesUsecase    sales_usecase.SalesUseCase
	salesRepository domain.SalesRepository
}

func NewSalesService(salesUsecase sales_usecase.SalesUseCase, salesRepository domain.SalesRepository) SalesService {
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
		installmentQuantity := 1
		if payment.InstallmentsQuantity != nil {
			installmentQuantity = *payment.InstallmentsQuantity
		}
		installmentsValue := s.calculateTotalValue(payment.Value, installmentQuantity)
		dateInformed := payment.FirstInstallmentDate != nil
		for i := 0; i < installmentQuantity; i++ {
			dueDate := time.Now()
			if payment.FirstInstallmentDate != nil {
				dueDate = *payment.FirstInstallmentDate
			}
			dates = append(dates, sales_usecase.DoSalePaymentDatesInput{
				DueDate:           dueDate.AddDate(0, i, 0),
				InstallmentNumber: i + 1,
				InstallmentValue:  installmentsValue[i],
				DateInformed:      dateInformed,
			})
		}
		payments = append(payments, sales_usecase.DoSalePaymentsInput{
			PaymentType: payment.PaymentType,
			Value:       payment.Value,
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

func (s *salesService) calculateTotalValue(total float64, n int) []float64 {
	totalCents := int64(math.Round(total * 100)) // ex.: 14.00 -> 1400
	base := totalCents / int64(n)
	rem := totalCents % int64(n)

	out := make([]int64, n)
	outFloat := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = base
		if int64(i) < rem {
			out[i]++ // distribui 1 centavo nas primeiras 'rem' parcelas
		}
	}
	for i, v := range out {
		outFloat[i] = float64(v) / 100
	}
	return outFloat
}

func (s *salesService) GetSales(ctx context.Context, request request.ListSalesRequest) (output output.GetSalesOutput, err error) {
	userRole := ctx.Value(constants.ROLE_KEY).(string)
	if userRole == "" {
		return output, ErrPermissionDenied
	}

	var userId []int64
	if userRole == string(domain.UserRoleAdmin) {
		userId = request.UserId
	} else {
		if v, ok := ctx.Value(constants.USERID_KEY).(float64); ok {
			userId = []int64{int64(v)}
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

func (s *salesService) GetById(ctx context.Context, id int64) (saleOutput output.GetSaleByIdOutput, paymentGroupOutput []output.GetSalesPaymentGroupOutput, itemsOutput []output.GetItemsOutput, err error) {
	saleOutput, err = s.salesRepository.GetSaleById(ctx, id)
	if err != nil {
		return saleOutput, paymentGroupOutput, itemsOutput, err
	}

	paymentOutput, err := s.salesRepository.GetPaymentsBySaleId(ctx, id)
	if err != nil {
		return saleOutput, paymentGroupOutput, itemsOutput, err
	}
	paymentGroupOutput = s.groupPaymentsByPaymentType(paymentOutput)

	itemsOutput, err = s.salesRepository.GetItemsBySaleId(ctx, id)
	if err != nil {
		return saleOutput, paymentGroupOutput, itemsOutput, err
	}

	return saleOutput, paymentGroupOutput, itemsOutput, err
}

func (s *salesService) groupPaymentsByPaymentType(payments []output.GetSalesPaymentOutput) []output.GetSalesPaymentGroupOutput {
	items := make([]output.GetSalesPaymentGroupOutput, 0)

	for _, payment := range payments {
		found := false
		for i, item := range items {
			if item.PaymentType == payment.PaymentType {
				items[i].Installments = append(items[i].Installments, payment)
				found = true
				break
			}
		}
		if !found {
			items = append(items, output.GetSalesPaymentGroupOutput{
				PaymentType:  payment.PaymentType,
				Installments: []output.GetSalesPaymentOutput{payment},
			})
		}
	}

	return items
}

func (s *salesService) ChangePaymentStatus(ctx context.Context, id int64, paymentId int64, request request.ChangePaymentStatusRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	_, err := s.salesRepository.GetPaymentDatesBySaleIdAndPaymentDateId(ctx, id, paymentId)
	if err != nil {
		return err
	}

	_, err = s.salesRepository.ChangePaymentStatus(ctx, paymentId, domain.PaymentStatus(request.Status))
	if err != nil {
		return err
	}

	var date *time.Time
	if request.Status == string(domain.PaymentStatusPaid) {
		date = &request.Date
	}
	_, err = s.salesRepository.ChangePaymentDate(ctx, paymentId, date)
	if err != nil {
		return err
	}

	return err
}
