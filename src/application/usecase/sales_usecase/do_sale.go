package sales_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

var (
	ErrSkusNotFound                    = errors.New("SKUs não encontrados")
	ErrPaymentValuesNotMatchTotalValue = errors.New("Valores de pagamento não correspondem ao valor total")
	ErrQuantityNotValid                = errors.New("Não há quantidade suficiente no estoque")
	ErrPaymentDatesPast                = errors.New("As datas de pagamento devem ser maiores que a data atual")
	ErrPaymentDatesOrderInvalid        = errors.New("As datas de pagamento devem ser ordenadas")
)

func (s *salesUseCase) DoSale(ctx context.Context, input DoSaleInput) error {

	user, err := s.userRepository.GetById(ctx, input.UserId)
	if err != nil {
		return err
	}

	customer, err := s.customerRepository.GetById(ctx, input.CustomerId)
	if err != nil {
		return err
	}

	skusIds := s.detachIds(input.Items)
	skus, err := s.skuRepository.GetByManyIds(ctx, skusIds)
	if err != nil {
		return err
	}
	if len(skus) != len(skusIds) {
		return ErrSkusNotFound
	}

	sale := s.createSale(user, customer, skus, input.Items, input.Payments)

	err = s.validateSale(sale)
	if err != nil {
		return err
	}

	inventoryOrigin, err := s.inventoryRepository.GetByUserId(ctx, user.Id)
	if err != nil && !errors.Is(err, repository.ErrInventoryNotFound) {
		return err
	}
	if err != nil && errors.Is(err, repository.ErrInventoryNotFound) && user.Role == string(domain.UserRoleAdmin) {
		inventoryOrigin, err = s.inventoryRepository.GetPrimaryInventory(ctx)
		if err != nil {
			return err
		}
	}

	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	skusInventoryInput := make([]inventory_usecase.DoTransactionSkusInput, len(sale.Items))
	for i, item := range sale.Items {
		skusInventoryInput[i] = inventory_usecase.DoTransactionSkusInput{
			SkuId:    item.Sku.Id,
			Quantity: item.Quantity,
		}
	}

	err = s.inventoryUseCase.DoTransaction(ctx, tx, inventory_usecase.DoTransactionInput{
		Type:                   domain.InventoryTransactionTypeOut,
		InventoryOriginId:      inventoryOrigin.Id,
		InventoryDestinationId: 0,
		Skus:                   skusInventoryInput,
		Justification:          "Vendido em " + time.Now().Format("02/01/2006"),
	})
	if err != nil {
		return err
	}

	sale.Id, err = s.saleRepository.CreateSale(ctx, tx, sale)
	if err != nil {
		return err
	}

	_, err = s.saleRepository.CreateManySaleItem(ctx, tx, sale, sale.Items)
	if err != nil {
		return err
	}

	for _, payment := range sale.Payments {
		payment.Id, err = s.saleRepository.CreatePayment(ctx, tx, sale, payment)
		if err != nil {
			return err
		}

		_, err = s.saleRepository.CreateManyPaymentDates(ctx, tx, payment, payment.Dates)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *salesUseCase) createSale(user domain.User, customer domain.Customer, skus []domain.Sku, itemsInput []DoSaleItemsInput, paymentsInput []DoSalePaymentsInput) domain.Sales {
	items := make([]domain.SalesItem, len(itemsInput))
	payments := make([]domain.SalesPayment, len(paymentsInput))

	for i, item := range itemsInput {
		for _, sku := range skus {
			if sku.Id == item.SkuId {
				items[i] = domain.NewSalesItem(sku, *sku.Price, item.Quantity)
				continue
			}
		}
	}
	for i, payment := range paymentsInput {
		paymentDates := make([]domain.SalesPaymentDates, len(payment.Dates))
		payments[i] = domain.NewSalesPayment(payment.PaymentType, paymentDates)
		for _, date := range payment.Dates {
			payments[i].AppendNewSalesDate(date.DueDate, date.PaidDate, date.InstallmentNumber, date.InstallmentValue)
		}
	}
	return domain.NewSales(time.Now(), user, customer, items, payments)
}

func (s *salesUseCase) detachIds(items []DoSaleItemsInput) []int64 {
	var skuIds []int64
	for _, item := range items {
		skuIds = append(skuIds, item.SkuId)
	}
	return skuIds
}

func (s *salesUseCase) validateSale(sale domain.Sales) error {
	if !sale.IsPaymentValuesMatchTotalValue() {
		return ErrPaymentValuesNotMatchTotalValue
	}
	for _, item := range sale.Items {
		if !item.IsQuantityValid() {
			return errors.New(ErrQuantityNotValid.Error() + fmt.Sprintf(": (%d) %s", item.Sku.Id, item.Sku.GetName()))
		}
	}
	for _, payment := range sale.Payments {
		if !payment.IsPaymentDatesOrderValid() {
			return ErrPaymentDatesOrderInvalid
		}
		if !payment.IsPaymentDatesGraterThanToday() {
			return ErrPaymentDatesPast
		}
	}
	return nil
}
