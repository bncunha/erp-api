package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	requests "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type QuoteService interface {
	Create(ctx context.Context, request requests.CreateQuoteRequest) (domain.Quote, error)
	Update(ctx context.Context, id int64, request requests.UpdateQuoteRequest) (domain.Quote, error)
	Duplicate(ctx context.Context, id int64) (int64, error)
	GetByID(ctx context.Context, id int64) (domain.Quote, error)
	List(ctx context.Context, input domain.GetQuotesInput) (domain.GetQuotesOutput, error)
	PatchStatus(ctx context.Context, id int64, request requests.PatchQuoteStatusRequest) (domain.Quote, error)
}

type quoteService struct {
	quoteRepository    domain.QuoteRepository
	customerRepository domain.CustomerRepository
	skuRepository      domain.SkuRepository
	repositories       *repository.Repository
}

func NewQuoteService(
	quoteRepository domain.QuoteRepository,
	customerRepository domain.CustomerRepository,
	skuRepository domain.SkuRepository,
	repositories *repository.Repository,
) QuoteService {
	return &quoteService{
		quoteRepository:    quoteRepository,
		customerRepository: customerRepository,
		skuRepository:      skuRepository,
		repositories:       repositories,
	}
}

func (s *quoteService) Create(ctx context.Context, request requests.CreateQuoteRequest) (domain.Quote, error) {
	var quote domain.Quote
	if err := request.Validate(); err != nil {
		return quote, err
	}
	validUntil, err := request.ParseValidUntil()
	if err != nil {
		return quote, err
	}
	userID := int64(ctx.Value(constants.USERID_KEY).(float64))
	quote = domain.Quote{
		CustomerID:            request.CustomerID,
		Status:                domain.QuoteStatusDraft,
		ValidUntil:            validUntil,
		DownPaymentPercentage: request.DownPaymentPercentage,
		DiscountPercentage:    request.DiscountPercentage,
		Notes:                 request.Notes,
		ShippingType:          request.ShippingType,
		ShippingCost:          request.ShippingCost,
		ShippingRegion:        request.ShippingRegion,
		ShippingMinValue:      request.ShippingMinValue,
		CreatedByUserID:       userID,
	}
	items, err := s.buildItems(ctx, request.Items)
	if err != nil {
		return quote, err
	}
	quote.Items = items
	if err := quote.RecalculateTotals(); err != nil {
		return quote, err
	}

	tx, err := s.repositories.BeginTx(ctx)
	if err != nil {
		return quote, err
	}
	defer rollbackOnError(tx, &err)

	if _, err = s.customerRepository.GetById(ctx, request.CustomerID); err != nil {
		return quote, err
	}
	if err = s.quoteRepository.Create(ctx, tx, &quote); err != nil {
		return quote, err
	}
	if err = tx.Commit(); err != nil {
		return quote, err
	}
	return s.GetByID(ctx, quote.ID)
}

func (s *quoteService) Update(ctx context.Context, id int64, request requests.UpdateQuoteRequest) (domain.Quote, error) {
	var quote domain.Quote
	if err := request.Validate(); err != nil {
		return quote, err
	}
	validUntil, err := request.ParseValidUntil()
	if err != nil {
		return quote, err
	}
	existing, err := s.quoteRepository.GetByID(ctx, id)
	if err != nil {
		return quote, err
	}
	if domain.IsQuoteStatusFinal(domain.EffectiveQuoteStatus(existing.Status, existing.ValidUntil, time.Now())) {
		return quote, errors.New("não é possível editar orçamento finalizado")
	}
	items, err := s.buildItems(ctx, request.Items)
	if err != nil {
		return quote, err
	}
	existing.CustomerID = request.CustomerID
	existing.ValidUntil = validUntil
	existing.DownPaymentPercentage = request.DownPaymentPercentage
	existing.DiscountPercentage = request.DiscountPercentage
	existing.Notes = request.Notes
	existing.ShippingType = request.ShippingType
	existing.ShippingCost = request.ShippingCost
	existing.ShippingRegion = request.ShippingRegion
	existing.ShippingMinValue = request.ShippingMinValue
	existing.Items = items
	if err := existing.RecalculateTotals(); err != nil {
		return quote, err
	}
	tx, err := s.repositories.BeginTx(ctx)
	if err != nil {
		return quote, err
	}
	defer rollbackOnError(tx, &err)

	if _, err = s.customerRepository.GetById(ctx, request.CustomerID); err != nil {
		return quote, err
	}
	if err = s.quoteRepository.Update(ctx, tx, &existing); err != nil {
		return quote, err
	}
	if err = tx.Commit(); err != nil {
		return quote, err
	}
	return s.GetByID(ctx, id)
}

func (s *quoteService) Duplicate(ctx context.Context, id int64) (int64, error) {
	existing, err := s.quoteRepository.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}

	items := make([]requests.UpsertQuoteItemRequest, 0, len(existing.Items))
	for _, item := range existing.Items {
		items = append(items, requests.UpsertQuoteItemRequest{
			SkuID:     item.SkuID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	request := requests.CreateQuoteRequest{
		CustomerID:            existing.CustomerID,
		ValidUntil:            existing.ValidUntil.Format(time.DateOnly),
		DownPaymentPercentage: existing.DownPaymentPercentage,
		DiscountPercentage:    existing.DiscountPercentage,
		ShippingType:          existing.ShippingType,
		ShippingCost:          existing.ShippingCost,
		ShippingRegion:        existing.ShippingRegion,
		ShippingMinValue:      existing.ShippingMinValue,
		Notes:                 existing.Notes,
		Items:                 items,
	}

	created, err := s.Create(ctx, request)
	if err != nil {
		return 0, err
	}

	return created.ID, nil
}

func (s *quoteService) GetByID(ctx context.Context, id int64) (domain.Quote, error) {
	quote, err := s.quoteRepository.GetByID(ctx, id)
	if err != nil {
		return quote, err
	}
	company, err := s.quoteRepository.GetCompanyProfile(ctx)
	if err != nil {
		return quote, err
	}
	quote.Company = company
	quote.Status = domain.EffectiveQuoteStatus(quote.Status, quote.ValidUntil, time.Now())
	quote.ShippingDescription = domain.BuildQuoteShippingDescription(quote.ShippingType, quote.ShippingRegion, quote.ShippingMinValue)
	return quote, nil
}

func (s *quoteService) List(ctx context.Context, input domain.GetQuotesInput) (domain.GetQuotesOutput, error) {
	requestedStatuses := append([]domain.QuoteStatus(nil), input.Statuses...)
	requestedStatusSet := make(map[domain.QuoteStatus]struct{}, len(requestedStatuses))
	hasExpiredStatus := false
	for _, status := range requestedStatuses {
		requestedStatusSet[status] = struct{}{}
		if status == domain.QuoteStatusExpired {
			hasExpiredStatus = true
		}
	}
	if hasExpiredStatus {
		// Expired is a computed status, so DB status filter cannot represent it safely.
		input.Statuses = nil
	}
	output, err := s.quoteRepository.List(ctx, input)
	if err != nil {
		return output, err
	}
	now := time.Now()
	filtered := make([]domain.QuoteListItemOutput, 0, len(output.Items))
	for _, item := range output.Items {
		item.Status = domain.EffectiveQuoteStatus(item.Status, item.ValidUntil, now)
		if len(requestedStatusSet) > 0 {
			if _, ok := requestedStatusSet[item.Status]; !ok {
				continue
			}
		}
		if input.OnlyExpired && item.Status != domain.QuoteStatusExpired {
			continue
		}
		item.DownPaymentAmount = item.TotalAmount * item.DownPaymentPercentage / 100
		filtered = append(filtered, item)
	}
	output.Items = filtered
	return output, nil
}

func (s *quoteService) PatchStatus(ctx context.Context, id int64, request requests.PatchQuoteStatusRequest) (domain.Quote, error) {
	var quote domain.Quote
	if err := request.Validate(); err != nil {
		return quote, err
	}
	existing, err := s.quoteRepository.GetByID(ctx, id)
	if err != nil {
		return quote, err
	}
	fromStatus := domain.EffectiveQuoteStatus(existing.Status, existing.ValidUntil, time.Now())
	if err := domain.ValidateQuoteStatusTransition(fromStatus, request.Status); err != nil {
		return quote, err
	}
	tx, err := s.repositories.BeginTx(ctx)
	if err != nil {
		return quote, err
	}
	defer rollbackOnError(tx, &err)

	if err = s.quoteRepository.UpdateStatus(ctx, tx, id, request.Status); err != nil {
		return quote, err
	}
	if err = tx.Commit(); err != nil {
		return quote, err
	}
	return s.GetByID(ctx, id)
}

func (s *quoteService) buildItems(ctx context.Context, itemsRequest []requests.UpsertQuoteItemRequest) ([]domain.QuoteItem, error) {
	if len(itemsRequest) == 0 {
		return nil, domain.ErrQuoteItemsRequired
	}
	skuIDs := make([]int64, 0, len(itemsRequest))
	for _, item := range itemsRequest {
		skuIDs = append(skuIDs, item.SkuID)
	}
	skus, err := s.skuRepository.GetByManyIds(ctx, skuIDs)
	if err != nil {
		return nil, err
	}
	skuMap := make(map[int64]domain.Sku, len(skus))
	for _, sku := range skus {
		skuMap[sku.Id] = sku
	}
	items := make([]domain.QuoteItem, 0, len(itemsRequest))
	for _, itemReq := range itemsRequest {
		sku, ok := skuMap[itemReq.SkuID]
		if !ok {
			return nil, errors.New("sku não encontrado")
		}
		snapshot := strings.TrimSpace(sku.GetName())
		if snapshot == "" {
			snapshot = fmt.Sprintf("SKU %d", sku.Id)
		}
		items = append(items, domain.QuoteItem{
			SkuID:                      itemReq.SkuID,
			ProductID:                  sku.Product.Id,
			SkuCode:                    sku.Code,
			ProductDescriptionSnapshot: snapshot,
			Quantity:                   itemReq.Quantity,
			UnitPrice:                  itemReq.UnitPrice,
		})
	}
	return items, nil
}

func rollbackOnError(tx *sql.Tx, err *error) {
	if *err != nil {
		_ = tx.Rollback()
	}
}
