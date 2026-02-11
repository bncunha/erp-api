package sales_usecase

import (
	"context"
	"database/sql"
	"database/sql/driver"
	stdErrors "errors"
	"strings"
	"testing"
	"time"

	serviceInput "github.com/bncunha/erp-api/src/application/service/input"
	serviceOutput "github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

func init() {
	sql.Register("sales_usecase_stub", stubDriver{})
}

var stubBeginErr error

type stubDriver struct{}

func (stubDriver) Open(string) (driver.Conn, error) {
	return &stubConn{}, nil
}

type stubConn struct{}

func (c *stubConn) Prepare(string) (driver.Stmt, error) {
	return &stubStmt{}, nil
}

func (c *stubConn) Close() error { return nil }

func (c *stubConn) Begin() (driver.Tx, error) {
	if stubBeginErr != nil {
		return nil, stubBeginErr
	}
	return &stubTx{}, nil
}

func (c *stubConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	if stubBeginErr != nil {
		return nil, stubBeginErr
	}
	return &stubTx{}, nil
}

type stubStmt struct{}

func (s *stubStmt) Close() error { return nil }

func (s *stubStmt) NumInput() int { return 0 }

func (s *stubStmt) Exec([]driver.Value) (driver.Result, error) {
	return nil, stdErrors.New("not implemented")
}

func (s *stubStmt) Query([]driver.Value) (driver.Rows, error) {
	return nil, stdErrors.New("not implemented")
}

type stubTx struct{}

func (t *stubTx) Commit() error { return nil }

func (t *stubTx) Rollback() error { return nil }

type fakeUserRepository struct {
	user domain.User
	err  error
}

func (f *fakeUserRepository) GetByUsername(context.Context, string) (domain.User, error) {
	return domain.User{}, nil
}

func (f *fakeUserRepository) Create(context.Context, domain.User) (int64, error) { return 0, nil }
func (f *fakeUserRepository) CreateWithTx(context.Context, *sql.Tx, domain.User) (int64, error) {
	return 0, nil
}

func (f *fakeUserRepository) Update(context.Context, domain.User) error { return nil }

func (f *fakeUserRepository) Inactivate(context.Context, int64) error { return nil }

func (f *fakeUserRepository) GetAll(context.Context, serviceInput.GetAllUserInput) ([]domain.User, error) {
	return nil, nil
}

func (f *fakeUserRepository) GetById(context.Context, int64) (domain.User, error) {
	return f.user, f.err
}

func (f *fakeUserRepository) UpdatePassword(context.Context, domain.User, string) error {
	return nil
}

func (f *fakeUserRepository) GetByEmail(context.Context, string) (domain.User, error) {
	return f.user, f.err
}

type fakeCustomerRepository struct {
	customer domain.Customer
	err      error
}

func (f *fakeCustomerRepository) GetById(context.Context, int64) (domain.Customer, error) {
	return f.customer, f.err
}

func (f *fakeCustomerRepository) GetAll(context.Context) ([]domain.Customer, error) { return nil, nil }

func (f *fakeCustomerRepository) Create(context.Context, domain.Customer) (int64, error) {
	return 0, nil
}

func (f *fakeCustomerRepository) Edit(context.Context, domain.Customer, int64) (int64, error) {
	return 0, nil
}

func (f *fakeCustomerRepository) Inactivate(context.Context, int64) error { return nil }

type fakeSkuRepository struct {
	skus []domain.Sku
	err  error
}

func (f *fakeSkuRepository) Create(context.Context, domain.Sku, int64) (int64, error) { return 0, nil }

func (f *fakeSkuRepository) CreateMany(context.Context, []domain.Sku, int64) ([]int64, error) {
	return nil, nil
}

func (f *fakeSkuRepository) GetByProductId(context.Context, int64) ([]domain.Sku, error) {
	return nil, nil
}

func (f *fakeSkuRepository) Update(context.Context, domain.Sku) error { return nil }

func (f *fakeSkuRepository) GetById(context.Context, int64) (domain.Sku, error) {
	return domain.Sku{}, nil
}

func (f *fakeSkuRepository) GetByManyIds(context.Context, []int64) ([]domain.Sku, error) {
	return f.skus, f.err
}

func (f *fakeSkuRepository) GetAll(context.Context, serviceInput.GetSkusInput) ([]domain.Sku, error) {
	return nil, nil
}

func (f *fakeSkuRepository) Inactivate(context.Context, int64) error { return nil }

type fakeInventoryRepository struct {
	byUser        domain.Inventory
	byUserErr     error
	primary       domain.Inventory
	primaryErr    error
	primaryCalled bool
}

func (f *fakeInventoryRepository) Create(context.Context, domain.Inventory) (int64, error) {
	return 0, nil
}

func (f *fakeInventoryRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, inventory domain.Inventory) (int64, error) {
	return f.Create(ctx, inventory)
}

func (f *fakeInventoryRepository) GetById(context.Context, int64) (domain.Inventory, error) {
	return domain.Inventory{}, nil
}

func (f *fakeInventoryRepository) GetAll(context.Context) ([]domain.Inventory, error) {
	return nil, nil
}

func (f *fakeInventoryRepository) GetByUserId(context.Context, int64) (domain.Inventory, error) {
	return f.byUser, f.byUserErr
}

func (f *fakeInventoryRepository) GetPrimaryInventory(context.Context) (domain.Inventory, error) {
	f.primaryCalled = true
	return f.primary, f.primaryErr
}

func (f *fakeInventoryRepository) GetSummary(context.Context) ([]serviceOutput.GetInventorySummaryOutput, error) {
	return nil, nil
}

func (f *fakeInventoryRepository) GetSummaryById(context.Context, int64) (serviceOutput.GetInventorySummaryByIdOutput, error) {
	return serviceOutput.GetInventorySummaryByIdOutput{}, nil
}

type fakeInventoryItemRepository struct {
	items    []domain.InventoryItem
	itemsErr error
}

func (f *fakeInventoryItemRepository) Create(context.Context, *sql.Tx, domain.InventoryItem) (int64, error) {
	return 0, nil
}

func (f *fakeInventoryItemRepository) UpdateQuantity(context.Context, *sql.Tx, domain.InventoryItem) error {
	return nil
}

func (f *fakeInventoryItemRepository) GetById(context.Context, int64) (domain.InventoryItem, error) {
	return domain.InventoryItem{}, nil
}

func (f *fakeInventoryItemRepository) GetByIdWithTransaction(context.Context, *sql.Tx, int64) (domain.InventoryItem, error) {
	return domain.InventoryItem{}, nil
}

func (f *fakeInventoryItemRepository) GetByManySkuIdsAndInventoryId(context.Context, []int64, int64) ([]domain.InventoryItem, error) {
	return f.items, f.itemsErr
}

func (f *fakeInventoryItemRepository) GetAll(context.Context) ([]serviceOutput.GetInventoryItemsOutput, error) {
	return nil, nil
}

func (f *fakeInventoryItemRepository) GetByInventoryId(context.Context, int64) ([]serviceOutput.GetInventoryItemsOutput, error) {
	return nil, nil
}

func (f *fakeInventoryItemRepository) GetBySkuId(context.Context, int64) ([]domain.GetSkuInventoryOutput, error) {
	return nil, nil
}

type fakeSalesRepository struct {
	sale             domain.Sales
	saleItems        []domain.SalesItem
	payments         []domain.SalesPayment
	paymentDates     [][]domain.SalesPaymentDates
	createSaleErr    error
	createSaleVersionErr error
	createItemsErr   error
	createPaymentErr error
	createDatesErr   error
	createReturnErr error
	createReturnItemsErr error
	cancelPaymentDatesErr error
	updateSaleLastVersionErr error
	saleByIdForUpdate domain.SaleWithVersionOutput
	saleByIdForUpdateErr error
	itemsByVersion []serviceOutput.GetItemsOutput
	itemsByVersionErr error
	paymentsByVersion []serviceOutput.GetSalesPaymentOutput
	paymentsByVersionErr error
	updateLastVersion int
	returnCreated bool
}

func (f *fakeSalesRepository) CreateSale(ctx context.Context, tx *sql.Tx, sale domain.Sales) (int64, error) {
	if f.createSaleErr != nil {
		return 0, f.createSaleErr
	}
	f.sale = sale
	return 101, nil
}

func (f *fakeSalesRepository) CreateSaleVersion(ctx context.Context, tx *sql.Tx, saleId int64, version int, date time.Time) (int64, error) {
	if f.createSaleVersionErr != nil {
		return 0, f.createSaleVersionErr
	}
	return 1001, nil
}

func (f *fakeSalesRepository) CreateManySaleItem(ctx context.Context, tx *sql.Tx, sale domain.Sales, items []domain.SalesItem) ([]int64, error) {
	if f.createItemsErr != nil {
		return nil, f.createItemsErr
	}
	f.saleItems = append([]domain.SalesItem(nil), items...)
	return []int64{1}, nil
}

func (f *fakeSalesRepository) CreatePayment(ctx context.Context, tx *sql.Tx, sale domain.Sales, payment domain.SalesPayment) (int64, error) {
	if f.createPaymentErr != nil {
		return 0, f.createPaymentErr
	}
	f.payments = append(f.payments, payment)
	return 201, nil
}

func (f *fakeSalesRepository) CreateManyPaymentDates(ctx context.Context, tx *sql.Tx, payment domain.SalesPayment, dates []domain.SalesPaymentDates) ([]int64, error) {
	if f.createDatesErr != nil {
		return nil, f.createDatesErr
	}
	copied := append([]domain.SalesPaymentDates(nil), dates...)
	f.paymentDates = append(f.paymentDates, copied)
	return []int64{1}, nil
}

func (f *fakeSalesRepository) CreateSalesReturn(ctx context.Context, tx *sql.Tx, saleId int64, fromSalesVersionId int64, toSalesVersionId int64, salesReturn domain.SalesReturn, createdByUserId int64) (int64, error) {
	if f.createReturnErr != nil {
		return 0, f.createReturnErr
	}
	f.returnCreated = true
	return 1, nil
}

func (f *fakeSalesRepository) CreateSalesReturnItems(ctx context.Context, tx *sql.Tx, salesReturnId int64, items []domain.SalesReturnItem) ([]int64, error) {
	if f.createReturnItemsErr != nil {
		return nil, f.createReturnItemsErr
	}
	return []int64{1}, nil
}

func (f *fakeSalesRepository) UpdateSaleLastVersion(ctx context.Context, tx *sql.Tx, saleId int64, version int) error {
	if f.updateSaleLastVersionErr != nil {
		return f.updateSaleLastVersionErr
	}
	f.updateLastVersion = version
	return nil
}

func (f *fakeSalesRepository) CancelPaymentDatesBySaleVersionId(ctx context.Context, tx *sql.Tx, saleVersionId int64) error {
	if f.cancelPaymentDatesErr != nil {
		return f.cancelPaymentDatesErr
	}
	return nil
}

func (f *fakeSalesRepository) GetSaleByIdForUpdate(ctx context.Context, tx *sql.Tx, id int64) (domain.SaleWithVersionOutput, error) {
	return f.saleByIdForUpdate, f.saleByIdForUpdateErr
}

func (f *fakeSalesRepository) GetSaleVersionIdBySaleIdAndVersion(ctx context.Context, tx *sql.Tx, saleId int64, version int) (int64, error) {
	return 1, nil
}

func (f *fakeSalesRepository) GetSales(context.Context, serviceInput.GetSalesInput) ([]serviceOutput.GetSalesItemOutput, error) {
	return nil, nil
}

func (f *fakeSalesRepository) GetSaleById(context.Context, int64) (serviceOutput.GetSaleByIdOutput, error) {
	return serviceOutput.GetSaleByIdOutput{}, nil
}

func (f *fakeSalesRepository) GetPaymentsBySaleId(context.Context, int64) ([]serviceOutput.GetSalesPaymentOutput, error) {
	return nil, nil
}

func (f *fakeSalesRepository) GetPaymentsBySaleVersionId(context.Context, int64) ([]serviceOutput.GetSalesPaymentOutput, error) {
	return f.paymentsByVersion, f.paymentsByVersionErr
}

func (f *fakeSalesRepository) GetItemsBySaleId(context.Context, int64) ([]serviceOutput.GetItemsOutput, error) {
	return nil, nil
}

func (f *fakeSalesRepository) GetItemsBySaleVersionId(context.Context, int64) ([]serviceOutput.GetItemsOutput, error) {
	return f.itemsByVersion, f.itemsByVersionErr
}

func (f *fakeSalesRepository) GetReturnsBySaleId(context.Context, int64) ([]serviceOutput.GetSalesReturnOutput, error) {
	return nil, nil
}

func (f *fakeSalesRepository) ChangePaymentStatus(context.Context, int64, domain.PaymentStatus) (int64, error) {
	return 0, nil
}

func (f *fakeSalesRepository) ChangePaymentDate(context.Context, int64, *time.Time) (int64, error) {
	return 0, nil
}

func (f *fakeSalesRepository) GetPaymentDatesBySaleIdAndPaymentDateId(context.Context, int64, int64) (domain.SalesPaymentDates, error) {
	return domain.SalesPaymentDates{}, nil
}

type fakeInventoryUseCase struct {
	received inventory_usecase.DoTransactionInput
	err      error
}

func (f *fakeInventoryUseCase) DoTransaction(ctx context.Context, tx *sql.Tx, input inventory_usecase.DoTransactionInput) error {
	f.received = input
	return f.err
}

func newStubRepository(t *testing.T) *repository.Repository {
	t.Helper()
	db, err := sql.Open("sales_usecase_stub", "")
	if err != nil {
		t.Fatalf("failed to open stub db: %v", err)
	}
	return repository.NewRepository(db)
}

type saleTestEnv struct {
	useCase           SalesUseCase
	userRepo          *fakeUserRepository
	customerRepo      *fakeCustomerRepository
	skuRepo           *fakeSkuRepository
	inventoryRepo     *fakeInventoryRepository
	inventoryItemRepo *fakeInventoryItemRepository
	salesRepo         *fakeSalesRepository
	inventoryUseCase  *fakeInventoryUseCase
	input             DoSaleInput
}

func newSaleTestEnv(t *testing.T) saleTestEnv {
	repo := newStubRepository(t)
	user := domain.User{Id: 1, Role: string(domain.UserRoleReseller)}
	customer := domain.Customer{Id: 2}
	sku := domain.Sku{Id: 3, Code: "SKU1", Price: 10, Quantity: 5, Product: domain.Product{Name: "Prod"}}
	userRepo := &fakeUserRepository{user: user}
	customerRepo := &fakeCustomerRepository{customer: customer}
	skuRepo := &fakeSkuRepository{skus: []domain.Sku{sku}}
	inventory := domain.Inventory{Id: 4}
	inventoryRepo := &fakeInventoryRepository{byUser: inventory}
	inventoryItems := []domain.InventoryItem{{Id: 5, InventoryId: inventory.Id, Sku: sku, Quantity: 5}}
	inventoryItemRepo := &fakeInventoryItemRepository{items: inventoryItems}
	salesRepo := &fakeSalesRepository{}
	inventoryUC := &fakeInventoryUseCase{}

	input := DoSaleInput{
		UserId:     user.Id,
		CustomerId: customer.Id,
		Items: []DoSaleItemsInput{{
			SkuId:    sku.Id,
			Quantity: 2,
		}},
		Payments: []DoSalePaymentsInput{{
			PaymentType: domain.PaymentTypeCash,
			Dates: []DoSalePaymentDatesInput{{
				DueDate:           time.Now().Add(24 * time.Hour),
				InstallmentNumber: 1,
				InstallmentValue:  20,
				DateInformed:      false,
			}},
		}},
	}

	useCase := NewSalesUseCase(userRepo, customerRepo, skuRepo, salesRepo, inventoryUC, inventoryRepo, inventoryItemRepo, repo)

	return saleTestEnv{
		useCase:           useCase,
		userRepo:          userRepo,
		customerRepo:      customerRepo,
		skuRepo:           skuRepo,
		inventoryRepo:     inventoryRepo,
		inventoryItemRepo: inventoryItemRepo,
		salesRepo:         salesRepo,
		inventoryUseCase:  inventoryUC,
		input:             input,
	}
}

func TestNewSalesUseCase(t *testing.T) {
	repo := newStubRepository(t)
	uc := NewSalesUseCase(&fakeUserRepository{}, &fakeCustomerRepository{}, &fakeSkuRepository{}, &fakeSalesRepository{}, &fakeInventoryUseCase{}, &fakeInventoryRepository{}, &fakeInventoryItemRepository{}, repo)
	impl, ok := uc.(*salesUseCase)
	if !ok {
		t.Fatalf("expected concrete sales use case type")
	}
	if impl.repository != repo {
		t.Fatalf("expected repository to be assigned")
	}
}

func TestSalesUseCaseDoSaleSuccess(t *testing.T) {
	env := newSaleTestEnv(t)

	if err := env.useCase.DoSale(context.Background(), env.input); err != nil {
		t.Fatalf("expected sale to succeed, got %v", err)
	}

	if len(env.salesRepo.sale.Items) != 1 || env.salesRepo.sale.Items[0].Sku.Id != env.skuRepo.skus[0].Id {
		t.Fatalf("expected sale items to be persisted: %+v", env.salesRepo.sale.Items)
	}
	if len(env.salesRepo.payments) != 1 || env.salesRepo.payments[0].PaymentType != domain.PaymentTypeCash {
		t.Fatalf("expected payment to be created: %+v", env.salesRepo.payments)
	}
	if len(env.salesRepo.paymentDates) != 1 || len(env.salesRepo.paymentDates[0]) != 1 {
		t.Fatalf("expected payment dates to be created: %+v", env.salesRepo.paymentDates)
	}
	if env.inventoryUseCase.received.Type != domain.InventoryTransactionTypeOut {
		t.Fatalf("expected inventory transaction type out, got %v", env.inventoryUseCase.received.Type)
	}
	if env.inventoryUseCase.received.InventoryOriginId != env.inventoryRepo.byUser.Id {
		t.Fatalf("expected inventory origin %d, got %d", env.inventoryRepo.byUser.Id, env.inventoryUseCase.received.InventoryOriginId)
	}
	if len(env.inventoryUseCase.received.Skus) != 1 || env.inventoryUseCase.received.Skus[0].Quantity != 2 {
		t.Fatalf("unexpected inventory sku payload: %+v", env.inventoryUseCase.received.Skus)
	}
	if env.inventoryUseCase.received.Sale.Id != 101 {
		t.Fatalf("expected sale id to be propagated, got %d", env.inventoryUseCase.received.Sale.Id)
	}
	if !strings.HasPrefix(env.inventoryUseCase.received.Justification, "Vendido em") {
		t.Fatalf("expected justification to mention sale date, got %s", env.inventoryUseCase.received.Justification)
	}
}

func TestSalesUseCaseDoSaleUserError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("user not found")
	env.userRepo.err = expectedErr

	if err := env.useCase.DoSale(context.Background(), env.input); err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoSaleInventoryFallback(t *testing.T) {
	env := newSaleTestEnv(t)
	env.userRepo.user.Role = string(domain.UserRoleAdmin)
	env.inventoryRepo.byUserErr = domain.ErrInventoryNotFound
	env.inventoryRepo.primary = domain.Inventory{Id: 9}
	env.inventoryItemRepo.items = []domain.InventoryItem{{Id: 10, InventoryId: 9, Sku: env.skuRepo.skus[0], Quantity: 5}}
	env.input.Items = []DoSaleItemsInput{{SkuId: env.skuRepo.skus[0].Id, Quantity: 1}}
	env.input.Payments[0].Dates[0].InstallmentValue = env.skuRepo.skus[0].Price

	if err := env.useCase.DoSale(context.Background(), env.input); err != nil {
		t.Fatalf("expected sale to succeed with fallback, got %v", err)
	}
	if !env.inventoryRepo.primaryCalled {
		t.Fatalf("expected primary inventory to be requested")
	}
	if env.inventoryUseCase.received.InventoryOriginId != env.inventoryRepo.primary.Id {
		t.Fatalf("expected inventory origin id %d, got %d", env.inventoryRepo.primary.Id, env.inventoryUseCase.received.InventoryOriginId)
	}
}

func TestSalesUseCaseDoSaleCustomerError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("customer error")
	env.customerRepo.err = expectedErr

	if err := env.useCase.DoSale(context.Background(), env.input); err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoSaleSkuRepositoryError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("sku error")
	env.skuRepo.err = expectedErr

	if err := env.useCase.DoSale(context.Background(), env.input); err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoSaleDuplicatedSkusError(t *testing.T) {
	env := newSaleTestEnv(t)
	env.input.Items = []DoSaleItemsInput{{SkuId: env.skuRepo.skus[0].Id}, {SkuId: env.skuRepo.skus[0].Id}}

	if err := env.useCase.DoSale(context.Background(), env.input); err == nil || !strings.Contains(err.Error(), domain.ErrSkusDuplicated.Error()) {
		t.Fatalf("expected duplicated sku error, got %v", err)
	}
}

func TestSalesUseCaseDoSaleInventoryRepositoryError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("inventory error")
	env.inventoryRepo.byUserErr = expectedErr

	if err := env.useCase.DoSale(context.Background(), env.input); err != expectedErr {
		t.Fatalf("expected inventory error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoSaleInventoryItemsRepositoryError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("inventory items error")
	env.inventoryItemRepo.itemsErr = expectedErr

	if err := env.useCase.DoSale(context.Background(), env.input); err != expectedErr {
		t.Fatalf("expected inventory items error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoSaleInventoryItemsValidationError(t *testing.T) {
	env := newSaleTestEnv(t)
	env.inventoryItemRepo.items = []domain.InventoryItem{{Sku: domain.Sku{Id: 999}}}

	if err := env.useCase.DoSale(context.Background(), env.input); err == nil || !strings.Contains(err.Error(), ErrSkusNotFound.Error()) {
		t.Fatalf("expected missing sku error, got %v", err)
	}
}

func TestSalesUseCaseDoSaleValidateSaleError(t *testing.T) {
	env := newSaleTestEnv(t)
	env.input.Payments[0].Dates[0].InstallmentValue = 1

	if err := env.useCase.DoSale(context.Background(), env.input); err == nil || !strings.Contains(err.Error(), domain.ErrPaymentValueIsMissing.Error()) {
		t.Fatalf("expected payment validation error, got %v", err)
	}
}

func TestSalesUseCaseDoSaleBeginTxError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("begin error")
	stubBeginErr = expectedErr
	t.Cleanup(func() { stubBeginErr = nil })

	if err := env.useCase.DoSale(context.Background(), env.input); err != expectedErr {
		t.Fatalf("expected begin tx error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoSaleCreateSaleError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("create sale error")
	env.salesRepo.createSaleErr = expectedErr

	if err := env.useCase.DoSale(context.Background(), env.input); err != expectedErr {
		t.Fatalf("expected create sale error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoSaleCreateItemsError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("create items error")
	env.salesRepo.createItemsErr = expectedErr

	if err := env.useCase.DoSale(context.Background(), env.input); err != expectedErr {
		t.Fatalf("expected create items error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoSaleCreatePaymentError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("create payment error")
	env.salesRepo.createPaymentErr = expectedErr

	if err := env.useCase.DoSale(context.Background(), env.input); err != expectedErr {
		t.Fatalf("expected create payment error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoSaleCreatePaymentDatesError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("create payment dates error")
	env.salesRepo.createDatesErr = expectedErr

	if err := env.useCase.DoSale(context.Background(), env.input); err != expectedErr {
		t.Fatalf("expected create payment dates error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoSaleInventoryTransactionError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("inventory transaction error")
	env.inventoryUseCase.err = expectedErr

	if err := env.useCase.DoSale(context.Background(), env.input); err != expectedErr {
		t.Fatalf("expected inventory transaction error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseValidateDuplicatedSkus(t *testing.T) {
	sku := domain.Sku{Id: 1, Code: "SKU1", Product: domain.Product{Name: "Prod"}}
	useCase := &salesUseCase{}
	err := useCase.validateDuplicatedSkus([]domain.Sku{sku}, []int64{sku.Id, sku.Id})
	if err == nil || !strings.Contains(err.Error(), domain.ErrSkusDuplicated.Error()) || !strings.Contains(err.Error(), sku.Code) {
		t.Fatalf("expected duplicated sku error, got %v", err)
	}

	if err := useCase.validateDuplicatedSkus([]domain.Sku{sku}, []int64{sku.Id}); err != nil {
		t.Fatalf("expected no error for unique skus, got %v", err)
	}
}

func TestSalesUseCaseValidateExistsInventoryItem(t *testing.T) {
	useCase := &salesUseCase{}
	inventoryItems := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}}}
	if err := useCase.validateExistsInventoryItem(inventoryItems, []int64{1, 2}); err == nil || !strings.Contains(err.Error(), ErrSkusNotFound.Error()) {
		t.Fatalf("expected missing sku error, got %v", err)
	}

	if err := useCase.validateExistsInventoryItem(inventoryItems, []int64{1}); err != nil {
		t.Fatalf("expected success when all inventory items exist, got %v", err)
	}
}

func TestSalesUseCaseCreateSale(t *testing.T) {
	useCase := &salesUseCase{}
	user := domain.User{Id: 1}
	customer := domain.Customer{Id: 2}
	sku := domain.Sku{Id: 3, Price: 10, Product: domain.Product{Name: "Prod"}}
	inventoryItems := []domain.InventoryItem{{Sku: sku, Quantity: 5}}
	due := time.Now().Add(48 * time.Hour)
	paymentsInput := []DoSalePaymentsInput{{
		PaymentType: domain.PaymentTypeCreditStore,
		Dates: []DoSalePaymentDatesInput{{
			DueDate:           due,
			InstallmentNumber: 1,
			InstallmentValue:  20,
			DateInformed:      false,
		}},
	}}

	sale := useCase.createSale(user, customer, inventoryItems, []DoSaleItemsInput{{SkuId: sku.Id, Quantity: 2}}, paymentsInput)

	if sale.User.Id != user.Id || sale.Customer.Id != customer.Id {
		t.Fatalf("expected sale to copy user and customer")
	}
	if len(sale.Items) != 1 || sale.Items[0].Quantity != 2 || sale.Items[0].Sku.Id != sku.Id {
		t.Fatalf("unexpected sale items: %+v", sale.Items)
	}
	if len(sale.Payments) != 1 || sale.Payments[0].PaymentType != domain.PaymentTypeCreditStore {
		t.Fatalf("unexpected sale payments: %+v", sale.Payments)
	}
	if len(sale.Payments[0].Dates) != 1 || sale.Payments[0].Dates[0].Status != domain.PaymentStatusPending {
		t.Fatalf("expected pending payment dates for credit store: %+v", sale.Payments[0].Dates)
	}
	if sale.Payments[0].Dates[0].InstallmentValue != 20 {
		t.Fatalf("expected installment value 20, got %v", sale.Payments[0].Dates[0].InstallmentValue)
	}
}

func TestSalesUseCaseDetachIds(t *testing.T) {
	useCase := &salesUseCase{}
	ids := useCase.detachIds([]DoSaleItemsInput{{SkuId: 1}, {SkuId: 2}})
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestSalesUseCaseDoReturnSuccess(t *testing.T) {
	env := newSaleTestEnv(t)
	env.salesRepo.saleByIdForUpdate = domain.SaleWithVersionOutput{
		Id:             101,
		LastVersion:    1,
		SalesVersionId: 1001,
	}
	env.salesRepo.itemsByVersion = []serviceOutput.GetItemsOutput{
		{Sku: domain.Sku{Id: 3, Price: 10, Product: domain.Product{Name: "Prod"}}, Quantity: 2, UnitPrice: 10},
	}
	env.salesRepo.paymentsByVersion = []serviceOutput.GetSalesPaymentOutput{
		{PaymentType: domain.PaymentTypeCash, InstallmentNumber: 1, InstallmentValue: 20, DueDate: time.Now(), PaymentStatus: domain.PaymentStatusPaid},
	}

	err := env.useCase.DoReturn(context.Background(), DoReturnInput{
		SaleId:                 101,
		UserId:                 1,
		InventoryDestinationId: 4,
		ReturnerName:           "Cliente",
		Reason:                 "Defeito",
		Items: []DoReturnItemInput{
			{SkuId: 3, Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !env.salesRepo.returnCreated {
		t.Fatalf("expected sales return to be created")
	}
	if env.salesRepo.updateLastVersion != 2 {
		t.Fatalf("expected last version to be updated to 2, got %d", env.salesRepo.updateLastVersion)
	}
	if env.inventoryUseCase.received.Type != domain.InventoryTransactionTypeIn {
		t.Fatalf("expected inventory IN transaction")
	}
}

func TestSalesUseCaseDoReturnValidationError(t *testing.T) {
	env := newSaleTestEnv(t)
	env.salesRepo.saleByIdForUpdate = domain.SaleWithVersionOutput{
		Id:             101,
		LastVersion:    1,
		SalesVersionId: 1001,
	}
	env.salesRepo.itemsByVersion = []serviceOutput.GetItemsOutput{
		{Sku: domain.Sku{Id: 3, Price: 10, Product: domain.Product{Name: "Prod"}}, Quantity: 1, UnitPrice: 10},
	}

	err := env.useCase.DoReturn(context.Background(), DoReturnInput{
		SaleId:                 101,
		UserId:                 1,
		InventoryDestinationId: 4,
		ReturnerName:           "Cliente",
		Reason:                 "Defeito",
		Items: []DoReturnItemInput{
			{SkuId: 3, Quantity: 2},
		},
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestSalesUseCaseDoReturnGetSaleByIdError(t *testing.T) {
	env := newSaleTestEnv(t)
	expectedErr := stdErrors.New("sale for update error")
	env.salesRepo.saleByIdForUpdateErr = expectedErr

	err := env.useCase.DoReturn(context.Background(), DoReturnInput{SaleId: 101, UserId: 1})
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoReturnGetItemsByVersionError(t *testing.T) {
	env := newSaleTestEnv(t)
	env.salesRepo.saleByIdForUpdate = domain.SaleWithVersionOutput{Id: 101, LastVersion: 1, SalesVersionId: 1001}
	expectedErr := stdErrors.New("items by version error")
	env.salesRepo.itemsByVersionErr = expectedErr

	err := env.useCase.DoReturn(context.Background(), DoReturnInput{SaleId: 101, UserId: 1})
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoReturnCreateVersionError(t *testing.T) {
	env := newSaleTestEnv(t)
	env.salesRepo.saleByIdForUpdate = domain.SaleWithVersionOutput{Id: 101, LastVersion: 1, SalesVersionId: 1001}
	env.salesRepo.itemsByVersion = []serviceOutput.GetItemsOutput{
		{Sku: domain.Sku{Id: 3, Price: 10, Product: domain.Product{Name: "Prod"}}, Quantity: 2, UnitPrice: 10},
	}
	expectedErr := stdErrors.New("create version error")
	env.salesRepo.createSaleVersionErr = expectedErr

	err := env.useCase.DoReturn(context.Background(), DoReturnInput{
		SaleId:                 101,
		UserId:                 1,
		InventoryDestinationId: 4,
		ReturnerName:           "Cliente",
		Reason:                 "Defeito",
		Items:                  []DoReturnItemInput{{SkuId: 3, Quantity: 1}},
	})
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestSalesUseCaseDoReturnCancelPaymentDatesError(t *testing.T) {
	env := newSaleTestEnv(t)
	env.salesRepo.saleByIdForUpdate = domain.SaleWithVersionOutput{Id: 101, LastVersion: 1, SalesVersionId: 1001}
	env.salesRepo.itemsByVersion = []serviceOutput.GetItemsOutput{
		{Sku: domain.Sku{Id: 3, Price: 10, Product: domain.Product{Name: "Prod"}}, Quantity: 2, UnitPrice: 10},
	}
	env.salesRepo.paymentsByVersion = []serviceOutput.GetSalesPaymentOutput{
		{PaymentType: domain.PaymentTypeCash, InstallmentNumber: 1, InstallmentValue: 20, DueDate: time.Now(), PaymentStatus: domain.PaymentStatusPaid},
	}
	expectedErr := stdErrors.New("cancel dates error")
	env.salesRepo.cancelPaymentDatesErr = expectedErr

	err := env.useCase.DoReturn(context.Background(), DoReturnInput{
		SaleId:                 101,
		UserId:                 1,
		InventoryDestinationId: 4,
		ReturnerName:           "Cliente",
		Reason:                 "Defeito",
		Items:                  []DoReturnItemInput{{SkuId: 3, Quantity: 1}},
	})
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestSplitAmount(t *testing.T) {
	values := splitAmount(10, 3)
	if len(values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(values))
	}
	if values[0] != 3.34 || values[1] != 3.33 || values[2] != 3.33 {
		t.Fatalf("unexpected split values: %v", values)
	}
}

func TestSalesUseCaseRecalculatePaymentsCreatesCreditStoreWhenNoPending(t *testing.T) {
	uc := &salesUseCase{}
	old := []domain.GetSalesPaymentOutput{
		{
			PaymentType:       domain.PaymentTypeCash,
			InstallmentNumber: 1,
			InstallmentValue:  10,
			DueDate:           time.Now(),
			PaymentStatus:     domain.PaymentStatusPaid,
		},
	}
	newItems := []domain.SalesItem{
		{Sku: domain.Sku{Price: 30}, Quantity: 1},
	}

	payments := uc.recalculatePayments(old, newItems)
	if len(payments) != 2 {
		t.Fatalf("expected 2 payment groups, got %d", len(payments))
	}
	foundCreditStore := false
	for _, p := range payments {
		if p.PaymentType == domain.PaymentTypeCreditStore {
			foundCreditStore = true
			if len(p.Dates) != 1 || p.Dates[0].InstallmentValue != 20 || p.Dates[0].Status != domain.PaymentStatusPending {
				t.Fatalf("unexpected credit store dates: %+v", p.Dates)
			}
		}
	}
	if !foundCreditStore {
		t.Fatalf("expected credit store payment to be created")
	}
}

func TestSalesUseCaseRecalculatePaymentsCreatesReturnWhenOverpaid(t *testing.T) {
	uc := &salesUseCase{}
	now := time.Now()
	old := []domain.GetSalesPaymentOutput{
		{
			PaymentType:       domain.PaymentTypeCash,
			InstallmentNumber: 1,
			InstallmentValue:  50,
			DueDate:           now,
			PaidDate:          &now,
			PaymentStatus:     domain.PaymentStatusPaid,
		},
	}
	newItems := []domain.SalesItem{
		{Sku: domain.Sku{Price: 30}, Quantity: 1},
	}

	payments := uc.recalculatePayments(old, newItems)
	foundReturn := false
	for _, p := range payments {
		if p.PaymentType == domain.PaymentTypeReturn {
			foundReturn = true
			if len(p.Dates) != 1 || p.Dates[0].InstallmentValue != -20 || p.Dates[0].Status != domain.PaymentStatusPaid {
				t.Fatalf("unexpected return payment dates: %+v", p.Dates)
			}
		}
	}
	if !foundReturn {
		t.Fatalf("expected return payment to be created")
	}
}

func TestSalesUseCaseRecalculatePaymentsDistributesPendingAmounts(t *testing.T) {
	uc := &salesUseCase{}
	old := []domain.GetSalesPaymentOutput{
		{
			PaymentType:       domain.PaymentTypeCreditCard,
			InstallmentNumber: 2,
			InstallmentValue:  0,
			DueDate:           time.Now().AddDate(0, 0, 2),
			PaymentStatus:     domain.PaymentStatusPending,
		},
		{
			PaymentType:       domain.PaymentTypeCreditCard,
			InstallmentNumber: 1,
			InstallmentValue:  0,
			DueDate:           time.Now().AddDate(0, 0, 1),
			PaymentStatus:     domain.PaymentStatusDelayed,
		},
		{
			PaymentType:       domain.PaymentTypeCreditCard,
			InstallmentNumber: 3,
			InstallmentValue:  0,
			DueDate:           time.Now().AddDate(0, 0, 3),
			PaymentStatus:     domain.PaymentStatusPending,
		},
	}
	newItems := []domain.SalesItem{
		{Sku: domain.Sku{Price: 10}, Quantity: 1},
	}

	payments := uc.recalculatePayments(old, newItems)
	if len(payments) != 1 {
		t.Fatalf("expected 1 payment group, got %d", len(payments))
	}
	if len(payments[0].Dates) != 3 {
		t.Fatalf("expected 3 payment dates, got %d", len(payments[0].Dates))
	}
	if payments[0].Dates[0].InstallmentValue != 3.34 || payments[0].Dates[1].InstallmentValue != 3.33 || payments[0].Dates[2].InstallmentValue != 3.33 {
		t.Fatalf("unexpected split across pending dates: %+v", payments[0].Dates)
	}
}
