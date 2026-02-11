package sales_usecase

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/domain"
)

func (s *salesUseCase) DoReturn(ctx context.Context, input DoReturnInput) (err error) {
	_, err = s.userRepository.GetById(ctx, input.UserId)
	if err != nil {
		return err
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

	sale, err := s.saleRepository.GetSaleByIdForUpdate(ctx, tx, input.SaleId)
	if err != nil {
		return err
	}

	currentItems, err := s.saleRepository.GetItemsBySaleVersionId(ctx, sale.SalesVersionId)
	if err != nil {
		return err
	}

	returnItems := make([]domain.SalesReturnItem, 0, len(input.Items))
	for _, item := range input.Items {
		for _, saleItem := range currentItems {
			if saleItem.Sku.Id == item.SkuId {
				returnItems = append(returnItems, domain.SalesReturnItem{
					Sku:       saleItem.Sku,
					Quantity:  item.Quantity,
					UnitPrice: saleItem.UnitPrice,
				})
				break
			}
		}
	}

	saleDomainItems := make([]domain.SalesItem, 0, len(currentItems))
	for _, item := range currentItems {
		saleDomainItems = append(saleDomainItems, domain.SalesItem{
			Sku:      item.Sku,
			Quantity: item.Quantity,
		})
	}

	salesReturn := domain.NewSalesReturn(input.ReturnerName, input.Reason, returnItems)
	if err = salesReturn.Validate(saleDomainItems); err != nil {
		return err
	}

	nextVersion := sale.LastVersion + 1
	nextSaleVersionId, err := s.saleRepository.CreateSaleVersion(ctx, tx, sale.Id, nextVersion, time.Now())
	if err != nil {
		return err
	}

	newSaleItems := s.buildRemainingItems(currentItems, salesReturn.Items)
	saleToCreate := domain.Sales{
		Id:             sale.Id,
		SalesVersionId: nextSaleVersionId,
		Items:          newSaleItems,
	}
	if len(newSaleItems) > 0 {
		if _, err = s.saleRepository.CreateManySaleItem(ctx, tx, saleToCreate, newSaleItems); err != nil {
			return err
		}
	}

	oldPayments, err := s.saleRepository.GetPaymentsBySaleVersionId(ctx, sale.SalesVersionId)
	if err != nil {
		return err
	}
	newPayments := s.recalculatePayments(oldPayments, newSaleItems)
	for _, payment := range newPayments {
		payment.Id, err = s.saleRepository.CreatePayment(ctx, tx, saleToCreate, payment)
		if err != nil {
			return err
		}
		if _, err = s.saleRepository.CreateManyPaymentDates(ctx, tx, payment, payment.Dates); err != nil {
			return err
		}
	}

	salesReturnId, err := s.saleRepository.CreateSalesReturn(ctx, tx, sale.Id, sale.SalesVersionId, nextSaleVersionId, salesReturn, input.UserId)
	if err != nil {
		return err
	}
	if _, err = s.saleRepository.CreateSalesReturnItems(ctx, tx, salesReturnId, salesReturn.Items); err != nil {
		return err
	}

	if err = s.saleRepository.CancelPaymentDatesBySaleVersionId(ctx, tx, sale.SalesVersionId); err != nil {
		return err
	}
	if err = s.saleRepository.UpdateSaleLastVersion(ctx, tx, sale.Id, nextVersion); err != nil {
		return err
	}

	stockItems := make([]inventory_usecase.DoTransactionSkusInput, 0, len(salesReturn.Items))
	for _, item := range salesReturn.Items {
		stockItems = append(stockItems, inventory_usecase.DoTransactionSkusInput{
			SkuId:    item.Sku.Id,
			Quantity: item.Quantity,
		})
	}

	err = s.inventoryUseCase.DoTransaction(ctx, tx, inventory_usecase.DoTransactionInput{
		Type:                   domain.InventoryTransactionTypeIn,
		InventoryDestinationId: input.InventoryDestinationId,
		Skus:                   stockItems,
		Sale: domain.Sales{
			Id:             sale.Id,
			SalesVersionId: nextSaleVersionId,
		},
		Justification: fmt.Sprintf("Produto devolvido no dia %s pelo cliente %s", time.Now().Format("02/01/2006"), input.ReturnerName),
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *salesUseCase) buildRemainingItems(current []domain.GetItemsOutput, returns []domain.SalesReturnItem) []domain.SalesItem {
	returnMap := make(map[int64]float64)
	for _, item := range returns {
		returnMap[item.Sku.Id] = item.Quantity
	}

	output := make([]domain.SalesItem, 0, len(current))
	for _, item := range current {
		remaining := item.Quantity - returnMap[item.Sku.Id]
		if remaining <= 0 {
			continue
		}
		output = append(output, domain.SalesItem{
			Sku:       item.Sku,
			Quantity:  remaining,
			UnitPrice: item.UnitPrice,
		})
	}
	return output
}

func (s *salesUseCase) recalculatePayments(old []domain.GetSalesPaymentOutput, newItems []domain.SalesItem) []domain.SalesPayment {
	newTotal := 0.0
	for _, item := range newItems {
		newTotal += item.Quantity * item.Sku.Price
	}
	newTotal = round2(newTotal)

	paidTotal := 0.0
	pending := make([]domain.GetSalesPaymentOutput, 0)
	settled := make([]domain.GetSalesPaymentOutput, 0)
	for _, date := range old {
		if date.PaymentStatus == domain.PaymentStatusPaid || date.PaymentStatus == domain.PaymentStatusReversal {
			settled = append(settled, date)
			paidTotal += date.InstallmentValue
			continue
		}
		if date.PaymentStatus == domain.PaymentStatusPending || date.PaymentStatus == domain.PaymentStatusDelayed {
			pending = append(pending, date)
		}
	}
	paidTotal = round2(paidTotal)

	paymentsMap := make(map[domain.PaymentType]*domain.SalesPayment)
	ensurePayment := func(t domain.PaymentType) *domain.SalesPayment {
		if p, ok := paymentsMap[t]; ok {
			return p
		}
		p := domain.NewSalesPayment(t)
		paymentsMap[t] = &p
		return &p
	}

	for _, p := range settled {
		payment := ensurePayment(p.PaymentType)
		dueDate := p.DueDate
		d := domain.NewSalesPaymentDates(dueDate, p.PaidDate, int(p.InstallmentNumber), p.InstallmentValue, p.PaymentStatus)
		d.PaymentType = p.PaymentType
		payment.Dates = append(payment.Dates, d)
		paymentsMap[p.PaymentType] = payment
	}

	remaining := round2(newTotal - paidTotal)
	if remaining > 0 {
		sort.Slice(pending, func(i, j int) bool {
			return pending[i].DueDate.Before(pending[j].DueDate)
		})
		if len(pending) == 0 {
			payment := ensurePayment(domain.PaymentTypeCreditStore)
			d := domain.NewSalesPaymentDates(time.Now(), nil, 1, remaining, domain.PaymentStatusPending)
			d.PaymentType = domain.PaymentTypeCreditStore
			payment.Dates = append(payment.Dates, d)
			paymentsMap[domain.PaymentTypeCreditStore] = payment
		} else {
			values := splitAmount(remaining, len(pending))
			for i, p := range pending {
				payment := ensurePayment(p.PaymentType)
				d := domain.NewSalesPaymentDates(p.DueDate, nil, int(p.InstallmentNumber), values[i], domain.PaymentStatusPending)
				d.PaymentType = p.PaymentType
				payment.Dates = append(payment.Dates, d)
				paymentsMap[p.PaymentType] = payment
			}
		}
	}

	if paidTotal > newTotal {
		diff := round2(paidTotal - newTotal)
		payment := ensurePayment(domain.PaymentTypeReturn)
		now := time.Now()
		d := domain.NewSalesPaymentDates(now, &now, 1, -diff, domain.PaymentStatusReversal)
		d.PaymentType = domain.PaymentTypeReturn
		payment.Dates = append(payment.Dates, d)
		paymentsMap[domain.PaymentTypeReturn] = payment
	}

	payments := make([]domain.SalesPayment, 0, len(paymentsMap))
	for _, payment := range paymentsMap {
		payments = append(payments, *payment)
	}

	sort.Slice(payments, func(i, j int) bool {
		return payments[i].PaymentType < payments[j].PaymentType
	})
	return payments
}

func splitAmount(total float64, n int) []float64 {
	totalCents := int64(math.Round(total * 100))
	base := totalCents / int64(n)
	rem := totalCents % int64(n)

	out := make([]float64, n)
	for i := 0; i < n; i++ {
		value := base
		if int64(i) < rem {
			value++
		}
		out[i] = float64(value) / 100
	}
	return out
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
