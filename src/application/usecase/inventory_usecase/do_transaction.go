package inventory_usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

var (
	ErrEnventoryItemDestinationNotFound = errors.New("Item de estoque não encontrado no destino")
	ErrInventoryItemOriginNotFound      = errors.New("Item de estoque não encontrado na origem")
	ErrQuantityInsufficient             = errors.New("Quantidade insuficiente")
	ErrInventoryesTransferEquais        = errors.New("Inventários de origem e de destino precisam ser diferentes")
	ErrInventoryItemNotFound            = errors.New("Item de estoque não encontrado")
	ErrSkusNotFound                     = errors.New("SKUs não encontrados")
	ErrSkusDuplicated                   = errors.New("SKUs duplicados encontrados")
)

type InventoryUseCase interface {
	DoTransaction(ctx context.Context, input DoTransactionInput) error
}

type inventoryUseCase struct {
	repository               *repository.Repository
	inventoryRepository      repository.InventoryRepository
	inventoryItemRepository  repository.InventoryItemRepository
	inventoryTransactionRepo repository.InventoryTransactionRepository
	skuRepository            repository.SkuRepository
}

func NewInventoryUseCase(repository *repository.Repository, inventoryRepository repository.InventoryRepository, inventoryItemRepository repository.InventoryItemRepository, inventoryTransactionRepo repository.InventoryTransactionRepository, skuRepository repository.SkuRepository) InventoryUseCase {
	return &inventoryUseCase{repository, inventoryRepository, inventoryItemRepository, inventoryTransactionRepo, skuRepository}
}

func (s *inventoryUseCase) DoTransaction(ctx context.Context, input DoTransactionInput) (err error) {

	var inventoryOut domain.Inventory
	var inventoryIn domain.Inventory
	var inventoryItemOut []domain.InventoryItem
	var inventoryItemIn []domain.InventoryItem
	skusIds := s.detachIds(input.Skus)

	if input.InventoryDestinationId != 0 {
		inventoryIn, err = s.inventoryRepository.GetById(ctx, input.InventoryDestinationId)
		if err != nil {
			return err
		}
	}
	if input.InventoryOriginId != 0 {
		inventoryOut, err = s.inventoryRepository.GetById(ctx, input.InventoryOriginId)
		if err != nil {
			return err
		}
	}

	skus, err := s.skuRepository.GetByManyIds(ctx, skusIds)
	if err != nil {
		return err
	}

	err = s.validateDuplicatedSkus(skus, skusIds)
	if err != nil {
		return err
	}

	err = s.validateExistsSkus(skus, skusIds)
	if err != nil {
		return err
	}

	for i, sku := range skus {
		for _, inputSku := range input.Skus {
			if sku.Id == inputSku.SkuId {
				skus[i].Quantity = inputSku.Quantity
			}
		}
	}

	inventoryItemOut, err = s.inventoryItemRepository.GetByManySkuIdsAndInventoryId(ctx, skusIds, inventoryOut.Id)
	if err != nil {
		return err
	}

	inventoryItemIn, err = s.inventoryItemRepository.GetByManySkuIdsAndInventoryId(ctx, skusIds, inventoryIn.Id)
	if err != nil {
		return err
	}

	err = s.validateInventoryTransaction(inventoryIn, inventoryOut, inventoryItemOut, inventoryItemIn, skus, input.Type)
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

	inventoryItemIn, err = s.createInventoryItemInIfNotExists(ctx, tx, skus, inventoryItemIn, domain.InventoryTransactionType(input.Type), inventoryIn)
	if err != nil {
		return err
	}

	err = s.createTransactions(ctx, tx, inventoryItemOut, inventoryItemIn, inventoryOut, inventoryIn, skus, input.Type, input.Justification)
	if err != nil {
		return err
	}

	switch input.Type {
	case domain.InventoryTransactionTypeTransfer:
		err = s.transferQuantity(ctx, tx, inventoryItemOut, inventoryItemIn, skus)
		if err != nil {
			return err
		}
	case domain.InventoryTransactionTypeIn:
		err = s.addQuantity(ctx, tx, inventoryItemIn, skus)
		if err != nil {
			return err
		}
	case domain.InventoryTransactionTypeOut:
		err = s.subQuantity(ctx, tx, inventoryItemOut, skus)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *inventoryUseCase) detachIds(skusInput []DoTransactionSkusInput) []int64 {
	var skuIds []int64
	for _, sku := range skusInput {
		skuIds = append(skuIds, sku.SkuId)
	}
	return skuIds
}

func (s *inventoryUseCase) findInventoryItem(inventoryItems []domain.InventoryItem, skuId int64) *domain.InventoryItem {
	for _, inventoryItem := range inventoryItems {
		if inventoryItem.SkuId == skuId {
			return &inventoryItem
		}
	}
	return nil
}

func (s *inventoryUseCase) createTransactions(ctx context.Context, tx *sql.Tx, inventoryItemsOut []domain.InventoryItem, inventoryItemsIn []domain.InventoryItem, inventoryOut domain.Inventory, inventoryIn domain.Inventory, inputSkus []domain.Sku, transactionType domain.InventoryTransactionType, justification string) error {
	for _, inputSku := range inputSkus {
		findedInventoryItemOut := s.findInventoryItem(inventoryItemsOut, inputSku.Id)
		findedInventoryItemIn := s.findInventoryItem(inventoryItemsIn, inputSku.Id)

		transaction := domain.InventoryTransaction{
			Quantity:      inputSku.Quantity,
			Date:          time.Now(),
			InventoryOut:  inventoryOut,
			InventoryIn:   inventoryIn,
			Justification: justification,
			Type:          transactionType,
		}
		if transactionType == domain.InventoryTransactionTypeTransfer || transactionType == domain.InventoryTransactionTypeOut {
			transaction.InventoryItem = *findedInventoryItemOut
		} else if transactionType == domain.InventoryTransactionTypeIn {
			transaction.InventoryItem = *findedInventoryItemIn
		}
		_, err := s.inventoryTransactionRepo.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *inventoryUseCase) createInventoryItemInIfNotExists(ctx context.Context, tx *sql.Tx, inputSku []domain.Sku, inventoriesItemIn []domain.InventoryItem, transactionType domain.InventoryTransactionType, inventoryIn domain.Inventory) ([]domain.InventoryItem, error) {
	var createdInventoryItemIn []domain.InventoryItem
	for _, inputSku := range inputSku {
		if s.findInventoryItem(inventoriesItemIn, inputSku.Id) == nil {
			if transactionType == domain.InventoryTransactionTypeIn || transactionType == domain.InventoryTransactionTypeTransfer {
				insertInventoryItem := domain.InventoryItem{
					InventoryId: inventoryIn.Id,
					SkuId:       inputSku.Id,
					Quantity:    0,
				}
				inventoryItemId, err := s.inventoryItemRepository.Create(ctx, tx, insertInventoryItem)
				if err != nil {
					return []domain.InventoryItem{}, err
				}
				inventoryItemIn, err := s.inventoryItemRepository.GetByIdWithTransaction(ctx, tx, inventoryItemId)
				if err != nil {
					return []domain.InventoryItem{}, err
				}
				createdInventoryItemIn = append(createdInventoryItemIn, inventoryItemIn)
			}
		}
	}
	return append(createdInventoryItemIn, inventoriesItemIn...), nil
}

func (s *inventoryUseCase) validateInventoryTransaction(inventoryIn domain.Inventory, inventoryOut domain.Inventory, inventoriesItemOut []domain.InventoryItem, inventoriesItemIn []domain.InventoryItem, skusInput []domain.Sku, inventoryType domain.InventoryTransactionType) error {
	if inventoryType == domain.InventoryTransactionTypeTransfer || inventoryType == domain.InventoryTransactionTypeOut {
		if validateExistingInventoryItemOutErr := s.validateExistingInventoryItemOut(inventoriesItemOut, skusInput); validateExistingInventoryItemOutErr != nil {
			return validateExistingInventoryItemOutErr
		}

		if validateInventoryItemOutQuantitiesErr := s.validateIInventotyItemOutQuantities(inventoriesItemOut, skusInput); validateInventoryItemOutQuantitiesErr != nil {
			return validateInventoryItemOutQuantitiesErr
		}
	}
	if inventoryType == domain.InventoryTransactionTypeTransfer {
		if inventoryIn.Id == inventoryOut.Id {
			return ErrInventoryesTransferEquais
		}
	}
	return nil
}

func (s *inventoryUseCase) validateExistingInventoryItemOut(inventoryItemOut []domain.InventoryItem, skusInput []domain.Sku) error {
	notExistingSkus := make([]string, 0)
	for _, sku := range skusInput {
		if s.findInventoryItem(inventoryItemOut, sku.Id) == nil {
			notExistingSkus = append(notExistingSkus, fmt.Sprintf(`(%d) %s`, sku.Id, sku.GetName()))
		}
	}
	if len(notExistingSkus) > 0 {
		return errors.New(ErrInventoryItemOriginNotFound.Error() + fmt.Sprintf(": %v", notExistingSkus))
	}
	return nil
}

func (s *inventoryUseCase) validateIInventotyItemOutQuantities(inventoryItemOut []domain.InventoryItem, skusInput []domain.Sku) error {
	for _, sku := range skusInput {
		for _, inventoryItem := range inventoryItemOut {
			if inventoryItem.SkuId == sku.Id {
				if sku.Quantity > inventoryItem.Quantity {
					return errors.New(ErrQuantityInsufficient.Error() + fmt.Sprintf(": (%d) %s", sku.Id, sku.GetName()))
				}
			}
		}
	}
	return nil
}

func (s *inventoryUseCase) validateExistsSkus(skus []domain.Sku, skusIds []int64) error {
	// Cria um mapa para marcar os SKUs encontrados
	found := make(map[int64]bool)
	for _, sku := range skus {
		found[sku.Id] = true
	}

	// Verifica se todos os IDs enviados pelo usuário estão no mapa
	for _, id := range skusIds {
		if !found[id] {
			return errors.New(ErrSkusNotFound.Error() + fmt.Sprintf(": %v", id))
		}
	}
	return nil
}

func (s *inventoryUseCase) validateDuplicatedSkus(skus []domain.Sku, skusIds []int64) error {
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
		errMEssage := ErrSkusDuplicated.Error() + ":"
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

func (s *inventoryUseCase) transferQuantity(ctx context.Context, tx *sql.Tx, inventoryItemsOut []domain.InventoryItem, inventoryItemsIn []domain.InventoryItem, inputSkus []domain.Sku) error {
	err := s.subQuantity(ctx, tx, inventoryItemsOut, inputSkus)
	if err != nil {
		return err
	}
	err = s.addQuantity(ctx, tx, inventoryItemsIn, inputSkus)
	if err != nil {
		return err
	}
	return nil
}

func (s *inventoryUseCase) addQuantity(ctx context.Context, tx *sql.Tx, inventoryItems []domain.InventoryItem, inputSkus []domain.Sku) error {
	for _, inputSku := range inputSkus {
		findedInventoryItem := s.findInventoryItem(inventoryItems, inputSku.Id)
		if findedInventoryItem != nil {
			findedInventoryItem.Quantity = findedInventoryItem.Quantity + inputSku.Quantity
			err := s.inventoryItemRepository.UpdateQuantity(ctx, tx, *findedInventoryItem)
			if err != nil {
				return err
			}
		} else {
			return errors.New(ErrInventoryItemNotFound.Error() + fmt.Sprintf(": %v", inputSku.Id))
		}
	}
	return nil
}

func (s *inventoryUseCase) subQuantity(ctx context.Context, tx *sql.Tx, inventoryItems []domain.InventoryItem, inputSkus []domain.Sku) error {
	for _, inputSku := range inputSkus {
		findedInventoryItem := s.findInventoryItem(inventoryItems, inputSku.Id)
		if findedInventoryItem != nil {
			findedInventoryItem.Quantity = findedInventoryItem.Quantity - inputSku.Quantity
			if findedInventoryItem.Quantity < 0 {
				return ErrQuantityInsufficient
			}
			err := s.inventoryItemRepository.UpdateQuantity(ctx, tx, *findedInventoryItem)
			if err != nil {
				return err
			}
		} else {
			return errors.New(ErrInventoryItemNotFound.Error() + fmt.Sprintf(": %v", inputSku.Id))
		}
	}
	return nil
}
