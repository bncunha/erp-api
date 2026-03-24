package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	requests "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type stubQuoteRepository struct {
	createErr       error
	updateErr       error
	updateStatusErr error
	getByIDErr      error
	listOutput      domain.GetQuotesOutput
	listErr         error
	listInput       domain.GetQuotesInput
	company         domain.QuoteCompanyProfile
	companyErr      error
	created         domain.Quote
	updated         domain.Quote
	updatedStatus   struct {
		id     int64
		status domain.QuoteStatus
	}
	quotesByID map[int64]domain.Quote
}

func (s *stubQuoteRepository) Create(ctx context.Context, tx *sql.Tx, quote *domain.Quote) error {
	if s.createErr != nil {
		return s.createErr
	}
	if s.quotesByID == nil {
		s.quotesByID = make(map[int64]domain.Quote)
	}
	if quote.ID == 0 {
		quote.ID = int64(len(s.quotesByID) + 1)
	}
	s.created = *quote
	s.quotesByID[quote.ID] = *quote
	return nil
}

func (s *stubQuoteRepository) Update(ctx context.Context, tx *sql.Tx, quote *domain.Quote) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	if s.quotesByID == nil {
		s.quotesByID = make(map[int64]domain.Quote)
	}
	s.updated = *quote
	s.quotesByID[quote.ID] = *quote
	return nil
}

func (s *stubQuoteRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, id int64, status domain.QuoteStatus) error {
	if s.updateStatusErr != nil {
		return s.updateStatusErr
	}
	s.updatedStatus = struct {
		id     int64
		status domain.QuoteStatus
	}{id: id, status: status}
	if quote, ok := s.quotesByID[id]; ok {
		quote.Status = status
		s.quotesByID[id] = quote
	}
	return nil
}

func (s *stubQuoteRepository) GetByID(ctx context.Context, id int64) (domain.Quote, error) {
	if s.getByIDErr != nil {
		return domain.Quote{}, s.getByIDErr
	}
	if quote, ok := s.quotesByID[id]; ok {
		return quote, nil
	}
	return domain.Quote{}, nil
}

func (s *stubQuoteRepository) List(ctx context.Context, input domain.GetQuotesInput) (domain.GetQuotesOutput, error) {
	s.listInput = input
	return s.listOutput, s.listErr
}

func (s *stubQuoteRepository) GetCompanyProfile(ctx context.Context) (domain.QuoteCompanyProfile, error) {
	if s.companyErr != nil {
		return domain.QuoteCompanyProfile{}, s.companyErr
	}
	return s.company, nil
}

func newQuoteServiceForTest(t *testing.T, quoteRepo *stubQuoteRepository, customerRepo *stubCustomerRepository, skuRepo *stubSkuRepository) (*quoteService, *fakeSQLTx, func()) {
	t.Helper()
	driverName := fmt.Sprintf("quote-fake-sql-%d", atomic.AddInt64(&fakeDriverCounter, 1))
	fake := &fakeSQLTx{}
	sql.Register(driverName, &fakeDriver{tx: fake})
	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}
	repos := repository.NewRepository(db)
	service := &quoteService{
		quoteRepository:    quoteRepo,
		customerRepository: customerRepo,
		skuRepository:      skuRepo,
		repositories:       repos,
	}
	cleanup := func() {
		_ = db.Close()
	}
	return service, fake, cleanup
}

func testQuoteCreateRequest() requests.CreateQuoteRequest {
	return requests.CreateQuoteRequest{
		CustomerID:            1,
		ValidUntil:            time.Now().AddDate(0, 0, 10).Format(time.DateOnly),
		DownPaymentPercentage: 20,
		DiscountPercentage:    10,
		ShippingType:          domain.QuoteShippingTypeFree,
		Items: []requests.UpsertQuoteItemRequest{
			{SkuID: 10, Quantity: 2, UnitPrice: 50},
		},
	}
}

func testQuoteUpdateRequest() requests.UpdateQuoteRequest {
	req := testQuoteCreateRequest()
	return requests.UpdateQuoteRequest{
		CustomerID:            req.CustomerID,
		ValidUntil:            req.ValidUntil,
		DownPaymentPercentage: req.DownPaymentPercentage,
		DiscountPercentage:    req.DiscountPercentage,
		ShippingType:          req.ShippingType,
		ShippingCost:          req.ShippingCost,
		ShippingRegion:        req.ShippingRegion,
		ShippingMinValue:      req.ShippingMinValue,
		Notes:                 req.Notes,
		Items:                 req.Items,
	}
}

func TestQuoteServiceCreate(t *testing.T) {
	quoteRepo := &stubQuoteRepository{
		quotesByID: map[int64]domain.Quote{},
		company:    domain.QuoteCompanyProfile{Name: "Trinus"},
	}
	customerRepo := &stubCustomerRepository{getById: domain.Customer{Id: 1}}
	skuRepo := &stubSkuRepository{
		getByMany: []domain.Sku{
			{Id: 10, Code: "SKU-10", Product: domain.Product{Id: 2, Name: "Mesa"}},
		},
	}
	service, fakeTx, cleanup := newQuoteServiceForTest(t, quoteRepo, customerRepo, skuRepo)
	defer cleanup()

	ctx := context.WithValue(context.Background(), constants.USERID_KEY, float64(7))
	quote, err := service.Create(ctx, testQuoteCreateRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fakeTx.committed {
		t.Fatalf("expected transaction commit")
	}
	if quote.ID == 0 {
		t.Fatalf("expected created quote id")
	}
	if quoteRepo.created.CreatedByUserID != 7 {
		t.Fatalf("expected created by user id 7, got %d", quoteRepo.created.CreatedByUserID)
	}
	if quote.TotalAmount != 90 {
		t.Fatalf("expected total amount 90, got %.2f", quote.TotalAmount)
	}
	if quote.Company.Name != "Trinus" {
		t.Fatalf("expected company profile to be attached")
	}
}

func TestQuoteServiceCreateRollbackOnError(t *testing.T) {
	quoteRepo := &stubQuoteRepository{quotesByID: map[int64]domain.Quote{}}
	customerRepo := &stubCustomerRepository{getByIdErr: errors.New("customer not found")}
	skuRepo := &stubSkuRepository{
		getByMany: []domain.Sku{{Id: 10, Code: "SKU-10", Product: domain.Product{Id: 2, Name: "Mesa"}}},
	}
	service, fakeTx, cleanup := newQuoteServiceForTest(t, quoteRepo, customerRepo, skuRepo)
	defer cleanup()

	ctx := context.WithValue(context.Background(), constants.USERID_KEY, float64(7))
	if _, err := service.Create(ctx, testQuoteCreateRequest()); err == nil {
		t.Fatalf("expected error")
	}
	if !fakeTx.rolledBack {
		t.Fatalf("expected rollback on error")
	}
}

func TestQuoteServiceUpdateWhenFinalStatus(t *testing.T) {
	quoteRepo := &stubQuoteRepository{
		quotesByID: map[int64]domain.Quote{
			1: {ID: 1, Status: domain.QuoteStatusApproved, ValidUntil: time.Now().AddDate(0, 0, 1)},
		},
	}
	service, _, cleanup := newQuoteServiceForTest(t, quoteRepo, &stubCustomerRepository{}, &stubSkuRepository{})
	defer cleanup()

	_, err := service.Update(context.Background(), 1, testQuoteUpdateRequest())
	if err == nil {
		t.Fatalf("expected error when updating final status quote")
	}
}

func TestQuoteServiceUpdate(t *testing.T) {
	quoteRepo := &stubQuoteRepository{
		quotesByID: map[int64]domain.Quote{
			1: {ID: 1, Status: domain.QuoteStatusDraft, ValidUntil: time.Now().AddDate(0, 0, 5), Items: []domain.QuoteItem{{SkuID: 10, Quantity: 1, UnitPrice: 10}}},
		},
		company: domain.QuoteCompanyProfile{Name: "Trinus"},
	}
	customerRepo := &stubCustomerRepository{getById: domain.Customer{Id: 1}}
	skuRepo := &stubSkuRepository{
		getByMany: []domain.Sku{{Id: 10, Code: "SKU-10", Product: domain.Product{Id: 2, Name: "Mesa"}}},
	}
	service, fakeTx, cleanup := newQuoteServiceForTest(t, quoteRepo, customerRepo, skuRepo)
	defer cleanup()

	req := testQuoteUpdateRequest()
	req.DiscountPercentage = 0
	updated, err := service.Update(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fakeTx.committed {
		t.Fatalf("expected commit")
	}
	if updated.TotalAmount != 100 {
		t.Fatalf("expected updated total 100, got %.2f", updated.TotalAmount)
	}
}

func TestQuoteServiceDuplicate(t *testing.T) {
	quoteRepo := &stubQuoteRepository{
		quotesByID: map[int64]domain.Quote{
			1: {
				ID:                    1,
				CustomerID:            1,
				Status:                domain.QuoteStatusDraft,
				ValidUntil:            time.Now().AddDate(0, 0, 5),
				DownPaymentPercentage: 20,
				DiscountPercentage:    5,
				ShippingType:          domain.QuoteShippingTypeFree,
				Items: []domain.QuoteItem{
					{SkuID: 10, Quantity: 2, UnitPrice: 50},
				},
			},
		},
		company: domain.QuoteCompanyProfile{Name: "Trinus"},
	}
	customerRepo := &stubCustomerRepository{getById: domain.Customer{Id: 1}}
	skuRepo := &stubSkuRepository{
		getByMany: []domain.Sku{{Id: 10, Code: "SKU-10", Product: domain.Product{Id: 2, Name: "Mesa"}}},
	}
	service, _, cleanup := newQuoteServiceForTest(t, quoteRepo, customerRepo, skuRepo)
	defer cleanup()

	ctx := context.WithValue(context.Background(), constants.USERID_KEY, float64(7))
	newID, err := service.Duplicate(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newID == 0 || newID == 1 {
		t.Fatalf("expected duplicated quote with new id, got %d", newID)
	}
}

func TestQuoteServiceGetByID(t *testing.T) {
	quoteRepo := &stubQuoteRepository{
		quotesByID: map[int64]domain.Quote{
			1: {
				ID:               1,
				Status:           domain.QuoteStatusSent,
				ValidUntil:       time.Now().AddDate(0, 0, -1),
				ShippingType:     domain.QuoteShippingTypeFreeRegion,
				ShippingRegion:   ptrString("Sul"),
				ShippingMinValue: nil,
			},
		},
		company: domain.QuoteCompanyProfile{Name: "Trinus"},
	}
	service, _, cleanup := newQuoteServiceForTest(t, quoteRepo, &stubCustomerRepository{}, &stubSkuRepository{})
	defer cleanup()

	quote, err := service.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quote.Status != domain.QuoteStatusExpired {
		t.Fatalf("expected expired status, got %s", quote.Status)
	}
	if quote.Company.Name != "Trinus" {
		t.Fatalf("expected company in response")
	}
	if quote.ShippingDescription == "" {
		t.Fatalf("expected shipping description")
	}
}

func TestQuoteServiceListFiltersAndMapsComputedFields(t *testing.T) {
	quoteRepo := &stubQuoteRepository{
		listOutput: domain.GetQuotesOutput{
			Items: []domain.QuoteListItemOutput{
				{
					ID:                    1,
					Status:                domain.QuoteStatusSent,
					ValidUntil:            time.Now().AddDate(0, 0, -2),
					TotalAmount:           200,
					DownPaymentPercentage: 25,
				},
				{
					ID:                    2,
					Status:                domain.QuoteStatusDraft,
					ValidUntil:            time.Now().AddDate(0, 0, 2),
					TotalAmount:           120,
					DownPaymentPercentage: 10,
				},
			},
		},
	}
	service, _, cleanup := newQuoteServiceForTest(t, quoteRepo, &stubCustomerRepository{}, &stubSkuRepository{})
	defer cleanup()

	output, err := service.List(context.Background(), domain.GetQuotesInput{
		Statuses:    []domain.QuoteStatus{domain.QuoteStatusExpired},
		OnlyExpired: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quoteRepo.listInput.Statuses != nil {
		t.Fatalf("expected statuses to be removed when requesting EXPIRED")
	}
	if len(output.Items) != 1 {
		t.Fatalf("expected one filtered item, got %d", len(output.Items))
	}
	if output.Items[0].DownPaymentAmount != 50 {
		t.Fatalf("expected down payment amount 50, got %.2f", output.Items[0].DownPaymentAmount)
	}
}

func TestQuoteServicePatchStatus(t *testing.T) {
	quoteRepo := &stubQuoteRepository{
		quotesByID: map[int64]domain.Quote{
			1: {
				ID:          1,
				Status:      domain.QuoteStatusDraft,
				ValidUntil:  time.Now().AddDate(0, 0, 2),
				ShippingType: domain.QuoteShippingTypeCalculate,
			},
		},
		company: domain.QuoteCompanyProfile{Name: "Trinus"},
	}
	service, fakeTx, cleanup := newQuoteServiceForTest(t, quoteRepo, &stubCustomerRepository{}, &stubSkuRepository{})
	defer cleanup()

	updated, err := service.PatchStatus(context.Background(), 1, requests.PatchQuoteStatusRequest{Status: domain.QuoteStatusSent})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fakeTx.committed {
		t.Fatalf("expected commit")
	}
	if quoteRepo.updatedStatus.status != domain.QuoteStatusSent {
		t.Fatalf("expected SENT status update")
	}
	if updated.Status != domain.QuoteStatusSent {
		t.Fatalf("expected returned status SENT, got %s", updated.Status)
	}
}

func TestQuoteServiceBuildItemsSkuNotFound(t *testing.T) {
	service := &quoteService{
		skuRepository: &stubSkuRepository{
			getByMany: []domain.Sku{},
		},
	}
	_, err := service.buildItems(context.Background(), []requests.UpsertQuoteItemRequest{
		{SkuID: 999, Quantity: 1, UnitPrice: 1},
	})
	if err == nil {
		t.Fatalf("expected sku not found error")
	}
}

func ptrString(value string) *string {
	return &value
}
