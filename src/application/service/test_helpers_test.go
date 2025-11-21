package service

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/bncunha/erp-api/src/application/service/input"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/application/usecase/sales_usecase"
	"github.com/bncunha/erp-api/src/domain"
)

var fakeDriverCounter int64

type fakeSQLTx struct {
	committed   bool
	rolledBack  bool
	commitErr   error
	rollbackErr error
}

func (f *fakeSQLTx) Commit() error {
	f.committed = true
	return f.commitErr
}

func (f *fakeSQLTx) Rollback() error {
	f.rolledBack = true
	return f.rollbackErr
}

type fakeDriver struct {
	tx *fakeSQLTx
}

func (d *fakeDriver) Open(name string) (driver.Conn, error) {
	return &fakeConn{tx: d.tx}, nil
}

type fakeConn struct {
	tx *fakeSQLTx
}

func (c *fakeConn) Prepare(query string) (driver.Stmt, error) { return fakeStmt{}, nil }
func (c *fakeConn) Close() error                              { return nil }
func (c *fakeConn) Begin() (driver.Tx, error)                 { return c.tx, nil }
func (c *fakeConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.tx, nil
}

type fakeStmt struct{}

func (fakeStmt) Close() error                                    { return nil }
func (fakeStmt) NumInput() int                                   { return 0 }
func (fakeStmt) Exec(args []driver.Value) (driver.Result, error) { return nil, nil }
func (fakeStmt) Query(args []driver.Value) (driver.Rows, error)  { return nil, nil }

func newTestSQLTx() (*sql.Tx, *fakeSQLTx, func()) {
	driverName := fmt.Sprintf("fake-sql-%d", atomic.AddInt64(&fakeDriverCounter, 1))
	fake := &fakeSQLTx{}
	sql.Register(driverName, &fakeDriver{tx: fake})
	db, _ := sql.Open(driverName, "")
	tx, _ := db.Begin()
	cleanup := func() { db.Close() }
	return tx, fake, cleanup
}

type stubTxManager struct {
	tx  *sql.Tx
	err error
}

func (s *stubTxManager) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.tx, nil
}

type stubEncrypto struct {
	encryptErr  error
	compareErr  error
	compareResp bool
}

func (s *stubEncrypto) Encrypt(text string) (string, error) {
	if s.encryptErr != nil {
		return "", s.encryptErr
	}
	return "encrypted:" + text, nil
}

func (s *stubEncrypto) Compare(hash string, text string) (bool, error) {
	if s.compareErr != nil {
		return false, s.compareErr
	}
	if s.compareResp {
		return s.compareResp, nil
	}
	return hash == "encrypted:"+text, nil
}

type stubEmailPort struct {
	sent []struct {
		senderEmail string
		senderName  string
		toEmail     string
		toName      string
		subject     string
		body        string
	}
	err error
}

func (s *stubEmailPort) Send(senderEmail string, senderName string, toEmail string, toName string, subject string, body string) error {
	s.sent = append(s.sent, struct {
		senderEmail string
		senderName  string
		toEmail     string
		toName      string
		subject     string
		body        string
	}{
		senderEmail: senderEmail,
		senderName:  senderName,
		toEmail:     toEmail,
		toName:      toName,
		subject:     subject,
		body:        body,
	})
	if s.err != nil {
		return s.err
	}
	return nil
}

type stubUserTokenService struct {
	createdInput input.CreateUserTokenInput
	output       domain.UserToken
	err          error
}

func (s *stubUserTokenService) Create(ctx context.Context, in input.CreateUserTokenInput) (domain.UserToken, error) {
	s.createdInput = in
	if s.err != nil {
		return domain.UserToken{}, s.err
	}
	if s.output.Code == "" {
		return domain.UserToken{Code: "code", Uuid: "uuid"}, nil
	}
	return s.output, nil
}

type stubEmailUseCase struct {
	inviteErr    error
	recoverErr   error
	inviteCalls  int
	recoverCalls int
	lastInvite   struct {
		user domain.User
		code string
		uuid string
	}
	lastRecover struct {
		user domain.User
		code string
		uuid string
	}
	recoverCh chan struct{}
}

func (s *stubEmailUseCase) SendInvite(ctx context.Context, user domain.User, code string, uuid string) error {
	s.inviteCalls++
	s.lastInvite = struct {
		user domain.User
		code string
		uuid string
	}{user: user, code: code, uuid: uuid}
	return s.inviteErr
}

func (s *stubEmailUseCase) SendRecoverPassword(ctx context.Context, user domain.User, code string, uuid string) error {
	s.recoverCalls++
	s.lastRecover = struct {
		user domain.User
		code string
		uuid string
	}{user: user, code: code, uuid: uuid}
	if s.recoverCh != nil {
		s.recoverCh <- struct{}{}
	}
	return s.recoverErr
}

type stubUserTokenRepository struct {
	createInput   domain.UserToken
	createErr     error
	createID      int64
	getById       domain.UserToken
	getByIdErr    error
	lastActive    domain.UserToken
	lastActiveErr error
	setUsedToken  domain.UserToken
	setUsedErr    error
}

func (s *stubUserTokenRepository) Create(ctx context.Context, userToken domain.UserToken) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.createInput = userToken
	if s.createID != 0 {
		return s.createID, nil
	}
	return 1, nil
}

func (s *stubUserTokenRepository) GetLastActiveByUuid(ctx context.Context, uuid string) (domain.UserToken, error) {
	if s.lastActiveErr != nil {
		return domain.UserToken{}, s.lastActiveErr
	}
	return s.lastActive, nil
}

func (s *stubUserTokenRepository) SetUsedToken(ctx context.Context, userToken domain.UserToken) error {
	s.setUsedToken = userToken
	if s.setUsedErr != nil {
		return s.setUsedErr
	}
	return nil
}

func (s *stubUserTokenRepository) GetById(ctx context.Context, id int64) (domain.UserToken, error) {
	return s.getById, s.getByIdErr
}

type stubProductRepository struct {
	created     domain.Product
	createErr   error
	editErr     error
	getById     domain.Product
	getByIdErr  error
	getAll      []output.GetAllProductsOutput
	getAllErr   error
	getAllInput input.GetProductsInput
}

func (s *stubProductRepository) Create(ctx context.Context, product domain.Product) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = product
	return 1, nil
}

func (s *stubProductRepository) Edit(ctx context.Context, product domain.Product, id int64) (int64, error) {
	if s.editErr != nil {
		return 0, s.editErr
	}
	s.created = product
	return id, nil
}

func (s *stubProductRepository) GetById(ctx context.Context, id int64) (domain.Product, error) {
	return s.getById, s.getByIdErr
}

func (s *stubProductRepository) GetAll(ctx context.Context, in input.GetProductsInput) ([]output.GetAllProductsOutput, error) {
	s.getAllInput = in
	return s.getAll, s.getAllErr
}

func (s *stubProductRepository) Inactivate(ctx context.Context, id int64) error {
	s.created = domain.Product{Id: id}
	return nil
}

type stubCategoryRepository struct {
	created      domain.Category
	createErr    error
	getById      domain.Category
	getByIdErr   error
	getByName    domain.Category
	getByNameErr error
	updateErr    error
	deleteErr    error
	getAll       []domain.Category
	getAllErr    error
}

func (s *stubCategoryRepository) Create(ctx context.Context, category domain.Category) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = category
	return 1, nil
}

func (s *stubCategoryRepository) GetById(ctx context.Context, id int64) (domain.Category, error) {
	return s.getById, s.getByIdErr
}

func (s *stubCategoryRepository) GetByName(ctx context.Context, name string) (domain.Category, error) {
	return s.getByName, s.getByNameErr
}

func (s *stubCategoryRepository) Update(ctx context.Context, category domain.Category) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.created = category
	return nil
}

func (s *stubCategoryRepository) Delete(ctx context.Context, id int64) error {
	return s.deleteErr
}

func (s *stubCategoryRepository) GetAll(ctx context.Context) ([]domain.Category, error) {
	return s.getAll, s.getAllErr
}

type stubCustomerRepository struct {
	created       domain.Customer
	createErr     error
	getAll        []domain.Customer
	getAllErr     error
	getById       domain.Customer
	getByIdErr    error
	editErr       error
	inactivateErr error
}

func (s *stubCustomerRepository) Create(ctx context.Context, customer domain.Customer) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = customer
	return 1, nil
}

func (s *stubCustomerRepository) GetAll(ctx context.Context) ([]domain.Customer, error) {
	return s.getAll, s.getAllErr
}

func (s *stubCustomerRepository) GetById(ctx context.Context, id int64) (domain.Customer, error) {
	return s.getById, s.getByIdErr
}

func (s *stubCustomerRepository) Edit(ctx context.Context, customer domain.Customer, id int64) (int64, error) {
	if s.editErr != nil {
		return 0, s.editErr
	}
	s.created = customer
	return id, nil
}

func (s *stubCustomerRepository) Inactivate(ctx context.Context, id int64) error {
	return s.inactivateErr
}

type stubSkuRepository struct {
	created         []domain.Sku
	createErr       error
	createManyErr   error
	updateErr       error
	getById         domain.Sku
	getByIdErr      error
	getByProduct    []domain.Sku
	getByProductErr error
	getAll          []domain.Sku
	getAllErr       error
	inactivateErr   error
	getAllInput     input.GetSkusInput
}

func (s *stubSkuRepository) Create(ctx context.Context, sku domain.Sku, productId int64) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = []domain.Sku{sku}
	return 1, nil
}

func (s *stubSkuRepository) CreateMany(ctx context.Context, skus []domain.Sku, productId int64) ([]int64, error) {
	if s.createManyErr != nil {
		return nil, s.createManyErr
	}
	s.created = append(s.created, skus...)
	return []int64{1}, nil
}

func (s *stubSkuRepository) GetByProductId(ctx context.Context, productId int64) ([]domain.Sku, error) {
	return s.getByProduct, s.getByProductErr
}

func (s *stubSkuRepository) Update(ctx context.Context, sku domain.Sku) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.created = []domain.Sku{sku}
	return nil
}

func (s *stubSkuRepository) GetById(ctx context.Context, id int64) (domain.Sku, error) {
	return s.getById, s.getByIdErr
}

func (s *stubSkuRepository) GetByManyIds(ctx context.Context, ids []int64) ([]domain.Sku, error) {
	return nil, nil
}

func (s *stubSkuRepository) GetAll(ctx context.Context, in input.GetSkusInput) ([]domain.Sku, error) {
	s.getAllInput = in
	return s.getAll, s.getAllErr
}

func (s *stubSkuRepository) Inactivate(ctx context.Context, id int64) error {
	return s.inactivateErr
}

type stubUserRepository struct {
	created           domain.User
	createErr         error
	updateErr         error
	updatePasswordErr error
	getById           domain.User
	getByIdErr        error
	getAll            []domain.User
	getAllErr         error
	getAllInput       input.GetAllUserInput
	inactivateErr     error
	getByUsername     domain.User
	getByUsernameErr  error
	getByEmail        domain.User
	getByEmailErr     error
	getByIdResponses  map[int64]domain.User
	getByIdErrors     map[int64]error
	updatePasswordReq struct {
		user        domain.User
		newPassword string
	}
}

func (s *stubUserRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = user
	return 1, nil
}

func (s *stubUserRepository) Update(ctx context.Context, user domain.User) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.created = user
	return nil
}

func (s *stubUserRepository) GetById(ctx context.Context, id int64) (domain.User, error) {
	if s.getByIdResponses != nil {
		if user, ok := s.getByIdResponses[id]; ok {
			if s.getByIdErrors != nil {
				if err, ok := s.getByIdErrors[id]; ok && err != nil {
					return domain.User{}, err
				}
			}
			return user, nil
		}
	}
	return s.getById, s.getByIdErr
}

func (s *stubUserRepository) GetAll(ctx context.Context, in input.GetAllUserInput) ([]domain.User, error) {
	s.getAllInput = in
	return s.getAll, s.getAllErr
}

func (s *stubUserRepository) Inactivate(ctx context.Context, id int64) error {
	return s.inactivateErr
}

func (s *stubUserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	return s.getByUsername, s.getByUsernameErr
}

func (s *stubUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.getByEmail, s.getByEmailErr
}

func (s *stubUserRepository) UpdatePassword(ctx context.Context, user domain.User, newPassword string) error {
	if s.updatePasswordErr != nil {
		return s.updatePasswordErr
	}
	s.updatePasswordReq = struct {
		user        domain.User
		newPassword string
	}{user: user, newPassword: newPassword}
	return nil
}

type stubInventoryRepository struct {
	createErr         error
	created           domain.Inventory
	getAll            []domain.Inventory
	getAllErr         error
	getByUser         domain.Inventory
	getByUserErr      error
	getById           domain.Inventory
	getByIdErr        error
	getPrimary        domain.Inventory
	getPrimaryErr     error
	getSummary        []output.GetInventorySummaryOutput
	getSummaryErr     error
	getSummaryById    output.GetInventorySummaryByIdOutput
	getSummaryByIdErr error
}

func (s *stubInventoryRepository) Create(ctx context.Context, inventory domain.Inventory) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = inventory
	return 1, nil
}

func (s *stubInventoryRepository) GetById(ctx context.Context, id int64) (domain.Inventory, error) {
	return s.getById, s.getByIdErr
}

func (s *stubInventoryRepository) GetPrimaryInventory(ctx context.Context) (domain.Inventory, error) {
	return s.getPrimary, s.getPrimaryErr
}

func (s *stubInventoryRepository) GetAll(ctx context.Context) ([]domain.Inventory, error) {
	return s.getAll, s.getAllErr
}

func (s *stubInventoryRepository) GetByUserId(ctx context.Context, userId int64) (domain.Inventory, error) {
	return s.getByUser, s.getByUserErr
}

func (s *stubInventoryRepository) GetSummary(ctx context.Context) ([]output.GetInventorySummaryOutput, error) {
	return s.getSummary, s.getSummaryErr
}

func (s *stubInventoryRepository) GetSummaryById(ctx context.Context, id int64) (output.GetInventorySummaryByIdOutput, error) {
	return s.getSummaryById, s.getSummaryByIdErr
}

type stubInventoryItemRepository struct {
	getAll            []output.GetInventoryItemsOutput
	getAllErr         error
	getByInventory    []output.GetInventoryItemsOutput
	getByInventoryErr error
}

func (s *stubInventoryItemRepository) GetAll(ctx context.Context) ([]output.GetInventoryItemsOutput, error) {
	return s.getAll, s.getAllErr
}

func (s *stubInventoryItemRepository) GetByInventoryId(ctx context.Context, id int64) ([]output.GetInventoryItemsOutput, error) {
	return s.getByInventory, s.getByInventoryErr
}

func (s *stubInventoryItemRepository) Create(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) (int64, error) {
	return 0, nil
}

func (s *stubInventoryItemRepository) UpdateQuantity(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) error {
	return nil
}

func (s *stubInventoryItemRepository) GetById(ctx context.Context, id int64) (domain.InventoryItem, error) {
	return domain.InventoryItem{}, nil
}

func (s *stubInventoryItemRepository) GetByIdWithTransaction(ctx context.Context, tx *sql.Tx, id int64) (domain.InventoryItem, error) {
	return domain.InventoryItem{}, nil
}

func (s *stubInventoryItemRepository) GetByManySkuIdsAndInventoryId(ctx context.Context, skuIds []int64, inventoryId int64) ([]domain.InventoryItem, error) {
	return nil, nil
}

type stubInventoryTransactionRepository struct {
	getAll            []output.GetInventoryTransactionsOutput
	getAllErr         error
	getByInventoryId  []output.GetInventoryTransactionsOutput
	getByInventoryErr error
}

func (s *stubInventoryTransactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction domain.InventoryTransaction) (int64, error) {
	return 0, nil
}

func (s *stubInventoryTransactionRepository) GetAll(ctx context.Context) ([]output.GetInventoryTransactionsOutput, error) {
	return s.getAll, s.getAllErr
}

func (s *stubInventoryTransactionRepository) GetByInventoryId(ctx context.Context, inventoryId int64) ([]output.GetInventoryTransactionsOutput, error) {
	if s.getByInventoryErr != nil {
		return nil, s.getByInventoryErr
	}
	if s.getByInventoryId != nil {
		return s.getByInventoryId, nil
	}
	return s.getAll, s.getAllErr
}

type stubRepository struct {
	beginErr error
}

func (s *stubRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return nil, s.beginErr
}

type stubInventoryUseCase struct {
	receivedInput inventory_usecase.DoTransactionInput
	err           error
}

func (s *stubInventoryUseCase) DoTransaction(ctx context.Context, tx *sql.Tx, input inventory_usecase.DoTransactionInput) error {
	s.receivedInput = input
	return s.err
}

type stubSalesUseCase struct {
	receivedInput sales_usecase.DoSaleInput
	err           error
}

func (s *stubSalesUseCase) DoSale(ctx context.Context, input sales_usecase.DoSaleInput) error {
	s.receivedInput = input
	return s.err
}

type stubSalesRepository struct {
	getSalesInput                 input.GetSalesInput
	getSalesOutput                []output.GetSalesItemOutput
	getSalesErr                   error
	saleByIdOutput                output.GetSaleByIdOutput
	saleByIdErr                   error
	paymentsOutput                []output.GetSalesPaymentOutput
	paymentsErr                   error
	itemsOutput                   []output.GetItemsOutput
	itemsErr                      error
	getSaleByIdCalled             bool
	getPaymentsCalled             bool
	getItemsCalled                bool
	changePaymentStatusCalledWith struct {
		id     int64
		status domain.PaymentStatus
	}
	changePaymentStatusErr      error
	changePaymentDateCalledWith struct {
		id   int64
		date *time.Time
	}
	changePaymentDateErr          error
	paymentDateBySaleAndPaymentId domain.SalesPaymentDates
	paymentDateErr                error
}

func (s *stubSalesRepository) CreateSale(ctx context.Context, tx *sql.Tx, sale domain.Sales) (int64, error) {
	return 0, nil
}

func (s *stubSalesRepository) CreateManySaleItem(ctx context.Context, tx *sql.Tx, sale domain.Sales, saleItems []domain.SalesItem) ([]int64, error) {
	return nil, nil
}

func (s *stubSalesRepository) CreatePayment(ctx context.Context, tx *sql.Tx, sale domain.Sales, payment domain.SalesPayment) (int64, error) {
	return 0, nil
}

func (s *stubSalesRepository) CreateManyPaymentDates(ctx context.Context, tx *sql.Tx, payment domain.SalesPayment, paymentDates []domain.SalesPaymentDates) ([]int64, error) {
	return nil, nil
}

func (s *stubSalesRepository) GetSales(ctx context.Context, input input.GetSalesInput) ([]output.GetSalesItemOutput, error) {
	s.getSalesInput = input
	return s.getSalesOutput, s.getSalesErr
}

func (s *stubSalesRepository) GetSaleById(ctx context.Context, id int64) (output.GetSaleByIdOutput, error) {
	s.getSaleByIdCalled = true
	return s.saleByIdOutput, s.saleByIdErr
}

func (s *stubSalesRepository) GetPaymentsBySaleId(ctx context.Context, id int64) ([]output.GetSalesPaymentOutput, error) {
	s.getPaymentsCalled = true
	return s.paymentsOutput, s.paymentsErr
}

func (s *stubSalesRepository) GetItemsBySaleId(ctx context.Context, id int64) ([]output.GetItemsOutput, error) {
	s.getItemsCalled = true
	return s.itemsOutput, s.itemsErr
}

func (s *stubSalesRepository) ChangePaymentStatus(ctx context.Context, id int64, status domain.PaymentStatus) (int64, error) {
	s.changePaymentStatusCalledWith = struct {
		id     int64
		status domain.PaymentStatus
	}{id: id, status: status}
	return id, s.changePaymentStatusErr
}

func (s *stubSalesRepository) ChangePaymentDate(ctx context.Context, id int64, date *time.Time) (int64, error) {
	s.changePaymentDateCalledWith = struct {
		id   int64
		date *time.Time
	}{id: id, date: date}
	return id, s.changePaymentDateErr
}

func (s *stubSalesRepository) GetPaymentDatesBySaleIdAndPaymentDateId(ctx context.Context, id int64, paymentDateId int64) (domain.SalesPaymentDates, error) {
	return s.paymentDateBySaleAndPaymentId, s.paymentDateErr
}
