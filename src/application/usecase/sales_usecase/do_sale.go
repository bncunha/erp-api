package sales_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/domain"
)

var (
	ErrSkusNotFound = errors.New("SKUs não encontrados")
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
	err = s.validateDuplicatedSkus(skus, skusIds)
	if err != nil {
		return err
	}

	inventoryOrigin, err := s.inventoryRepository.GetByUserId(ctx, user.Id)
	if err != nil && !errors.Is(err, domain.ErrInventoryNotFound) {
		return err
	}
	if err != nil && errors.Is(err, domain.ErrInventoryNotFound) && user.Role == string(domain.UserRoleAdmin) {
		inventoryOrigin, err = s.inventoryRepository.GetPrimaryInventory(ctx)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	inventoryItems, err := s.inventoryItemRepository.GetByManySkuIdsAndInventoryId(ctx, skusIds, inventoryOrigin.Id)
	if err != nil {
		return err
	}
	err = s.validateExistsInventoryItem(inventoryItems, skusIds)
	if err != nil {
		return err
	}

	sale := s.createSale(user, customer, inventoryItems, input.Items, input.Payments)

	err = sale.ValidateSale()
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

	skusInventoryInput := make([]inventory_usecase.DoTransactionSkusInput, len(sale.Items))
	for i, item := range sale.Items {
		skusInventoryInput[i] = inventory_usecase.DoTransactionSkusInput{
			SkuId:    item.Sku.Id,
			Quantity: item.Quantity,
		}
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

	err = s.inventoryUseCase.DoTransaction(ctx, tx, inventory_usecase.DoTransactionInput{
		Type:                   domain.InventoryTransactionTypeOut,
		InventoryOriginId:      inventoryOrigin.Id,
		InventoryDestinationId: 0,
		Skus:                   skusInventoryInput,
		Sale:                   sale,
		Justification:          "Vendido em " + time.Now().Format("02/01/2006"),
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *salesUseCase) createSale(user domain.User, customer domain.Customer, inventoryItems []domain.InventoryItem, itemsInput []DoSaleItemsInput, paymentsInput []DoSalePaymentsInput) domain.Sales {
	items := make([]domain.SalesItem, len(itemsInput))
	payments := make([]domain.SalesPayment, len(paymentsInput))

	for i, input := range itemsInput {
		for _, item := range inventoryItems {
			if item.Sku.Id == input.SkuId {
				items[i] = domain.NewSalesItem(item.Sku, input.Quantity)
				continue
			}
		}
	}
	for i, payment := range paymentsInput {
		payments[i] = domain.NewSalesPayment(payment.PaymentType)
		for _, date := range payment.Dates {
			payments[i].AppendNewSalesDate(date.DueDate, date.InstallmentNumber, date.InstallmentValue, date.DateInformed)
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

func (s *salesUseCase) validateExistsInventoryItem(inventoryItems []domain.InventoryItem, skusIds []int64) error {
	// Cria um mapa para marcar os SKUs encontrados
	found := make(map[int64]bool)
	for _, item := range inventoryItems {
		found[item.Sku.Id] = true
	}

	// Verifica se todos os IDs enviados pelo usuário estão no mapa
	for _, id := range skusIds {
		if !found[id] {
			return errors.New(ErrSkusNotFound.Error() + fmt.Sprintf(": %v", id))
		}
	}
	return nil
}

func (s *salesUseCase) validateDuplicatedSkus(skus []domain.Sku, skusIds []int64) error {
	seen := make(map[int64]bool)
	duplicates := []int64{}

	for _, id := range skusIds {
		if seen[id] {
			duplicates = append(duplicates, id)
		} else {
			seen[id] = true
		}
	}

	if len(duplicates) > 0 {
		errMEssage := domain.ErrSkusDuplicated.Error() + ":"
		for _, id := range duplicates {
			for _, sku := range skus {
				if sku.Id == id {
					errMEssage = errMEssage + fmt.Sprintf("- (%s) %s | ", sku.Code, sku.GetName())
				}
			}
		}
		return errors.New(errMEssage)
	}
	return nil
}
