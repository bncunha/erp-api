package inventory_usecase

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/bncunha/erp-api/src/application/service/input"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type stubInventoryItemRepository struct {
	createdItems          []domain.InventoryItem
	returnedInventoryItem domain.InventoryItem
	createErr             error
	getErr                error
	updateErr             error
	getManyErr            error
	items                 map[int64][]domain.InventoryItem
}

func (s *stubInventoryItemRepository) Create(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) (int64, error) {
	s.createdItems = append(s.createdItems, inventoryItem)
	if s.createErr != nil {
		return 0, s.createErr
	}
	if s.items != nil {
		s.items[inventoryItem.InventoryId] = append(s.items[inventoryItem.InventoryId], inventoryItem)
	}
	return 1, nil
}

func (s *stubInventoryItemRepository) UpdateQuantity(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.returnedInventoryItem = inventoryItem
	return nil
}

func (s *stubInventoryItemRepository) GetById(ctx context.Context, id int64) (domain.InventoryItem, error) {
	return domain.InventoryItem{}, nil
}

func (s *stubInventoryItemRepository) GetByIdWithTransaction(ctx context.Context, tx *sql.Tx, id int64) (domain.InventoryItem, error) {
	if s.getErr != nil {
		return domain.InventoryItem{}, s.getErr
	}
	return s.returnedInventoryItem, nil
}

func (s *stubInventoryItemRepository) GetByManySkuIdsAndInventoryId(ctx context.Context, skuIds []int64, inventoryId int64) ([]domain.InventoryItem, error) {
	if s.getManyErr != nil {
		return nil, s.getManyErr
	}
	if s.items != nil {
		return append([]domain.InventoryItem{}, s.items[inventoryId]...), nil
	}
	return []domain.InventoryItem{}, nil
}

func (s *stubInventoryItemRepository) GetAll(ctx context.Context) ([]output.GetInventoryItemsOutput, error) {
	return nil, nil
}

func (s *stubInventoryItemRepository) GetByInventoryId(ctx context.Context, id int64) ([]output.GetInventoryItemsOutput, error) {
	return nil, nil
}

func (s *stubInventoryItemRepository) GetBySkuId(ctx context.Context, skuId int64) ([]domain.GetSkuInventoryOutput, error) {
	return nil, nil
}

type stubInventoryTransactionRepository struct {
	created []domain.InventoryTransaction
	err     error
}

func (s *stubInventoryTransactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction domain.InventoryTransaction) (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	s.created = append(s.created, transaction)
	return int64(len(s.created)), nil
}

func (s *stubInventoryTransactionRepository) GetAll(ctx context.Context) ([]output.GetInventoryTransactionsOutput, error) {
	return nil, nil
}

func (s *stubInventoryTransactionRepository) GetByInventoryId(ctx context.Context, inventoryId int64) ([]output.GetInventoryTransactionsOutput, error) {
	return nil, nil
}

func (s *stubInventoryTransactionRepository) GetBySkuId(ctx context.Context, skuId int64) ([]output.GetInventoryTransactionsOutput, error) {
	return nil, nil
}

func TestNewInventoryUseCase(t *testing.T) {
	repo := &repository.Repository{}
	uc := NewInventoryUseCase(repo, nil, nil, nil, nil)
	if uc == nil {
		t.Fatalf("expected inventory use case instance")
	}
}

func TestDetachIds(t *testing.T) {
	uc := &inventoryUseCase{}
	ids := uc.detachIds([]DoTransactionSkusInput{{SkuId: 1}, {SkuId: 2}})
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestFindInventoryItem(t *testing.T) {
	uc := &inventoryUseCase{}
	item := uc.findInventoryItem([]domain.InventoryItem{{Sku: domain.Sku{Id: 1}}}, 1)
	if item == nil || item.Sku.Id != 1 {
		t.Fatalf("expected to find inventory item")
	}
	if uc.findInventoryItem([]domain.InventoryItem{{Sku: domain.Sku{Id: 2}}}, 1) != nil {
		t.Fatalf("expected nil when not found")
	}
}

func TestCreateInventoryItemInIfNotExists(t *testing.T) {
	repo := &stubInventoryItemRepository{}
	uc := &inventoryUseCase{inventoryItemRepository: repo}
	var tx *sql.Tx
	inventory := domain.Inventory{Id: 10}
	sku := domain.Sku{Id: 1}

	repo.returnedInventoryItem = domain.InventoryItem{Id: 2, Sku: domain.Sku{Id: sku.Id}, InventoryId: inventory.Id}
	items, err := uc.createInventoryItemInIfNotExists(context.Background(), tx, []domain.Sku{sku}, nil, domain.InventoryTransactionTypeIn, inventory)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0].Sku.Id != sku.Id {
		t.Fatalf("expected created inventory item, got %v", items)
	}
}

func TestCreateInventoryItemInIfNotExistsAlreadyExists(t *testing.T) {
	repo := &stubInventoryItemRepository{}
	uc := &inventoryUseCase{inventoryItemRepository: repo}
	var tx *sql.Tx
	inventory := domain.Inventory{Id: 10}
	sku := domain.Sku{Id: 1}
	existing := []domain.InventoryItem{{Sku: domain.Sku{Id: sku.Id}}}

	items, err := uc.createInventoryItemInIfNotExists(context.Background(), tx, []domain.Sku{sku}, existing, domain.InventoryTransactionTypeTransfer, inventory)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || repo.createdItems != nil {
		t.Fatalf("expected existing items only")
	}
}

func TestValidateInventoryTransaction(t *testing.T) {
	uc := &inventoryUseCase{}
	items := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 5}}
	skus := []domain.Sku{{Id: 1, Quantity: 3}}

	if err := uc.validateInventoryTransaction(domain.Inventory{Id: 1}, domain.Inventory{Id: 1}, items, nil, skus, domain.InventoryTransactionTypeTransfer); err == nil {
		t.Fatalf("expected error when inventories equal")
	}

	err := uc.validateInventoryTransaction(domain.Inventory{Id: 1}, domain.Inventory{Id: 2}, items, nil, skus, domain.InventoryTransactionTypeTransfer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateInventoryTransactionOutType(t *testing.T) {
	uc := &inventoryUseCase{}
	items := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 5}}
	skus := []domain.Sku{{Id: 1, Quantity: 2}}
	if err := uc.validateInventoryTransaction(domain.Inventory{}, domain.Inventory{Id: 1}, items, nil, skus, domain.InventoryTransactionTypeOut); err != nil {
		t.Fatalf("unexpected error validating out transaction: %v", err)
	}
}

func TestValidateExistingInventoryItemOut(t *testing.T) {
	uc := &inventoryUseCase{}
	err := uc.validateExistingInventoryItemOut([]domain.InventoryItem{}, []domain.Sku{{Id: 1}})
	if err == nil {
		t.Fatalf("expected error when sku missing")
	}
	err = uc.validateExistingInventoryItemOut([]domain.InventoryItem{{Sku: domain.Sku{Id: 1}}}, []domain.Sku{{Id: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateInventoryItemOutQuantities(t *testing.T) {
	uc := &inventoryUseCase{}
	err := uc.validateIInventotyItemOutQuantities([]domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 1}}, []domain.Sku{{Id: 1, Quantity: 2}})
	if err == nil {
		t.Fatalf("expected quantity error")
	}
	err = uc.validateIInventotyItemOutQuantities([]domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 2}}, []domain.Sku{{Id: 1, Quantity: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateExistsSkus(t *testing.T) {
	uc := &inventoryUseCase{}
	err := uc.validateExistsSkus([]domain.Sku{{Id: 1}}, []int64{1, 2})
	if err == nil {
		t.Fatalf("expected error for missing sku")
	}
	if err := uc.validateExistsSkus([]domain.Sku{{Id: 1}, {Id: 2}}, []int64{1, 2}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateDuplicatedSkus(t *testing.T) {
	uc := &inventoryUseCase{}
	err := uc.validateDuplicatedSkus([]domain.Sku{{Id: 1, Code: "A", Product: domain.Product{Name: "P"}}}, []int64{1, 1})
	if err == nil {
		t.Fatalf("expected duplicate error")
	}
	if err := uc.validateDuplicatedSkus([]domain.Sku{{Id: 1}}, []int64{1, 2}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransferQuantity(t *testing.T) {
	repo := &stubInventoryItemRepository{}
	uc := &inventoryUseCase{inventoryItemRepository: repo}
	var tx *sql.Tx
	out := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 5}}
	in := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 1}}
	skus := []domain.Sku{{Id: 1, Quantity: 2}}

	if err := uc.transferQuantity(context.Background(), tx, out, in, skus); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransferQuantityErrorOnSubQuantity(t *testing.T) {
	repo := &stubInventoryItemRepository{}
	uc := &inventoryUseCase{inventoryItemRepository: repo}
	out := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 1}}
	in := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 0}}
	skus := []domain.Sku{{Id: 1, Quantity: 5}}

	if err := uc.transferQuantity(context.Background(), nil, out, in, skus); err == nil || err != ErrQuantityInsufficient {
		t.Fatalf("expected quantity insufficient error, got %v", err)
	}
}

func TestTransferQuantityErrorOnAddQuantity(t *testing.T) {
	repo := &stubInventoryItemRepository{}
	uc := &inventoryUseCase{inventoryItemRepository: repo}
	out := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 5}}
	in := []domain.InventoryItem{} // missing destination item triggers add error
	skus := []domain.Sku{{Id: 1, Quantity: 2}}

	if err := uc.transferQuantity(context.Background(), nil, out, in, skus); err == nil || err.Error() == "" {
		t.Fatalf("expected error when destination inventory item missing")
	}
}

func TestAddQuantity(t *testing.T) {
	repo := &stubInventoryItemRepository{}
	uc := &inventoryUseCase{inventoryItemRepository: repo}
	var tx *sql.Tx
	items := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 1, Id: 10}}
	skus := []domain.Sku{{Id: 1, Quantity: 2}}

	if err := uc.addQuantity(context.Background(), tx, items, skus); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.returnedInventoryItem.Quantity != 3 {
		t.Fatalf("expected quantity 3, got %v", repo.returnedInventoryItem.Quantity)
	}

	if err := uc.addQuantity(context.Background(), tx, nil, skus); err == nil {
		t.Fatalf("expected error for missing inventory item")
	}
}

func TestAddQuantityUpdateError(t *testing.T) {
	repo := &stubInventoryItemRepository{updateErr: errors.New("fail")}
	uc := &inventoryUseCase{inventoryItemRepository: repo}
	items := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 1}}
	skus := []domain.Sku{{Id: 1, Quantity: 1}}
	if err := uc.addQuantity(context.Background(), nil, items, skus); err == nil || err.Error() != "fail" {
		t.Fatalf("expected update error")
	}
}

func TestSubQuantity(t *testing.T) {
	repo := &stubInventoryItemRepository{}
	uc := &inventoryUseCase{inventoryItemRepository: repo}
	var tx *sql.Tx
	items := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 3, Id: 10}}
	skus := []domain.Sku{{Id: 1, Quantity: 2}}

	if err := uc.subQuantity(context.Background(), tx, items, skus); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.returnedInventoryItem.Quantity != 1 {
		t.Fatalf("expected quantity 1, got %v", repo.returnedInventoryItem.Quantity)
	}

	if err := uc.subQuantity(context.Background(), tx, []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Quantity: 1}}, []domain.Sku{{Id: 1, Quantity: 2}}); err == nil {
		t.Fatalf("expected insufficient quantity error")
	}
}

func TestCreateTransactions(t *testing.T) {
	repo := &stubInventoryTransactionRepository{}
	uc := &inventoryUseCase{inventoryTransactionRepo: repo}
	var tx *sql.Tx
	outItems := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Id: 1}}
	inItems := []domain.InventoryItem{{Sku: domain.Sku{Id: 1}, Id: 2}}
	skus := []domain.Sku{{Id: 1, Quantity: 2}}

	if err := uc.createTransactions(context.Background(), tx, outItems, inItems, domain.Inventory{}, domain.Inventory{}, skus, domain.InventoryTransactionTypeIn, "just", domain.Sales{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.created) != 1 {
		t.Fatalf("expected transaction to be created")
	}
	if repo.created[0].Justification != "just" {
		t.Fatalf("unexpected justification")
	}
}

type fakeTx struct {
	committed  bool
	rolledBack bool
}

func (f *fakeTx) Commit() error {
	f.committed = true
	return nil
}

func (f *fakeTx) Rollback() error {
	f.rolledBack = true
	return nil
}

type fakeDriver struct {
	tx *fakeTx
}

func (d *fakeDriver) Open(name string) (driver.Conn, error) {
	return &fakeConn{tx: d.tx}, nil
}

type fakeConn struct {
	tx *fakeTx
}

func (c *fakeConn) Prepare(query string) (driver.Stmt, error) { return fakeStmt{}, nil }
func (c *fakeConn) Close() error                              { return nil }
func (c *fakeConn) Begin() (driver.Tx, error)                 { return c.tx, nil }
func (c *fakeConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.tx, nil
}

type fakeStmt struct{}

func (fakeStmt) Close() error  { return nil }
func (fakeStmt) NumInput() int { return 0 }
func (fakeStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.New("not implemented")
}
func (fakeStmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, errors.New("not implemented")
}

type doTxInventoryRepository struct {
	inventories map[int64]domain.Inventory
}

func (r *doTxInventoryRepository) Create(ctx context.Context, inventory domain.Inventory) (int64, error) {
	return 0, nil
}

func (r *doTxInventoryRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, inventory domain.Inventory) (int64, error) {
	return r.Create(ctx, inventory)
}

func (r *doTxInventoryRepository) GetById(ctx context.Context, id int64) (domain.Inventory, error) {
	if inv, ok := r.inventories[id]; ok {
		return inv, nil
	}
	return domain.Inventory{}, errors.New("not found")
}

func (r *doTxInventoryRepository) GetAll(ctx context.Context) ([]domain.Inventory, error) {
	return nil, nil
}
func (r *doTxInventoryRepository) GetByUserId(ctx context.Context, userId int64) (domain.Inventory, error) {
	return domain.Inventory{}, nil
}

func (r *doTxInventoryRepository) GetPrimaryInventory(ctx context.Context) (domain.Inventory, error) {
	return domain.Inventory{}, nil
}

func (r *doTxInventoryRepository) GetSummary(ctx context.Context) ([]output.GetInventorySummaryOutput, error) {
	return nil, nil
}

func (r *doTxInventoryRepository) GetSummaryById(ctx context.Context, id int64) (output.GetInventorySummaryByIdOutput, error) {
	return output.GetInventorySummaryByIdOutput{}, nil
}

type doTxInventoryItemRepository struct {
	items          map[int64][]domain.InventoryItem
	createdItems   []domain.InventoryItem
	updatedItems   []domain.InventoryItem
	nextID         int64
	lastCreatedFor map[int64]domain.InventoryItem
}

func newDoTxInventoryItemRepository() *doTxInventoryItemRepository {
	return &doTxInventoryItemRepository{items: make(map[int64][]domain.InventoryItem), lastCreatedFor: make(map[int64]domain.InventoryItem)}
}

func (r *doTxInventoryItemRepository) Create(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) (int64, error) {
	r.nextID++
	inventoryItem.Id = r.nextID
	r.createdItems = append(r.createdItems, inventoryItem)
	r.items[inventoryItem.InventoryId] = append(r.items[inventoryItem.InventoryId], inventoryItem)
	r.lastCreatedFor[inventoryItem.Id] = inventoryItem
	return inventoryItem.Id, nil
}

func (r *doTxInventoryItemRepository) UpdateQuantity(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem) error {
	r.updatedItems = append(r.updatedItems, inventoryItem)
	return nil
}

func (r *doTxInventoryItemRepository) GetById(ctx context.Context, id int64) (domain.InventoryItem, error) {
	return domain.InventoryItem{}, nil
}

func (r *doTxInventoryItemRepository) GetByIdWithTransaction(ctx context.Context, tx *sql.Tx, id int64) (domain.InventoryItem, error) {
	return r.lastCreatedFor[id], nil
}

func (r *doTxInventoryItemRepository) GetByManySkuIdsAndInventoryId(ctx context.Context, skuIds []int64, inventoryId int64) ([]domain.InventoryItem, error) {
	return append([]domain.InventoryItem{}, r.items[inventoryId]...), nil
}

func (r *doTxInventoryItemRepository) GetAll(ctx context.Context) ([]output.GetInventoryItemsOutput, error) {
	return nil, nil
}

func (r *doTxInventoryItemRepository) GetByInventoryId(ctx context.Context, id int64) ([]output.GetInventoryItemsOutput, error) {
	return nil, nil
}

func (r *doTxInventoryItemRepository) GetBySkuId(ctx context.Context, skuId int64) ([]domain.GetSkuInventoryOutput, error) {
	return nil, nil
}

type doTxInventoryTransactionRepository struct {
	transactions []domain.InventoryTransaction
}

func (r *doTxInventoryTransactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction domain.InventoryTransaction) (int64, error) {
	r.transactions = append(r.transactions, transaction)
	return int64(len(r.transactions)), nil
}

func (r *doTxInventoryTransactionRepository) GetAll(ctx context.Context) ([]output.GetInventoryTransactionsOutput, error) {
	return nil, nil
}

func (r *doTxInventoryTransactionRepository) GetBySkuId(ctx context.Context, skuId int64) ([]output.GetInventoryTransactionsOutput, error) {
	return nil, nil
}

func (r *doTxInventoryTransactionRepository) GetByInventoryId(ctx context.Context, inventoryId int64) ([]output.GetInventoryTransactionsOutput, error) {
	return nil, nil
}

type doTxSkuRepository struct {
	skus []domain.Sku
}

func (r *doTxSkuRepository) GetByManyIds(ctx context.Context, ids []int64) ([]domain.Sku, error) {
	return append([]domain.Sku{}, r.skus...), nil
}

func (r *doTxSkuRepository) Create(ctx context.Context, sku domain.Sku, productId int64) (int64, error) {
	return 0, nil
}
func (r *doTxSkuRepository) CreateMany(ctx context.Context, skus []domain.Sku, productId int64) ([]int64, error) {
	return nil, nil
}
func (r *doTxSkuRepository) GetByProductId(ctx context.Context, productId int64) ([]domain.Sku, error) {
	return nil, nil
}
func (r *doTxSkuRepository) Update(ctx context.Context, sku domain.Sku) error { return nil }
func (r *doTxSkuRepository) GetById(ctx context.Context, id int64) (domain.Sku, error) {
	return domain.Sku{}, nil
}
func (r *doTxSkuRepository) GetAll(ctx context.Context, _ input.GetSkusInput) ([]domain.Sku, error) {
	return nil, nil
}
func (r *doTxSkuRepository) Inactivate(ctx context.Context, id int64) error { return nil }

var fakeDriverCounter int64

func newFakeDB(tx *fakeTx) *sql.DB {
	name := fmt.Sprintf("fake-driver-%d", atomic.AddInt64(&fakeDriverCounter, 1))
	sql.Register(name, &fakeDriver{tx: tx})
	db, _ := sql.Open(name, "")
	return db
}

func TestInventoryUseCaseDoTransaction(t *testing.T) {
	fakeTx := &fakeTx{}
	db := newFakeDB(fakeTx)
	repos := repository.NewRepository(db)
	tx, _ := db.BeginTx(context.Background(), nil)

	inventoryRepo := &doTxInventoryRepository{inventories: map[int64]domain.Inventory{1: {Id: 1}, 2: {Id: 2}}}
	itemRepo := newDoTxInventoryItemRepository()
	itemRepo.items[1] = []domain.InventoryItem{{Id: 1, InventoryId: 1, Sku: domain.Sku{Id: 1}, Quantity: 5}}
	txRepo := &doTxInventoryTransactionRepository{}
	skuRepo := &doTxSkuRepository{skus: []domain.Sku{{Id: 1, Code: "A", Product: domain.Product{Name: "Prod"}}}}

	uc := &inventoryUseCase{
		repository:               repos,
		inventoryRepository:      inventoryRepo,
		inventoryItemRepository:  itemRepo,
		inventoryTransactionRepo: txRepo,
		skuRepository:            skuRepo,
	}

	err := uc.DoTransaction(context.Background(), tx, DoTransactionInput{
		Type:                   domain.InventoryTransactionTypeTransfer,
		InventoryOriginId:      1,
		InventoryDestinationId: 2,
		Justification:          "move",
		Skus:                   []DoTransactionSkusInput{{SkuId: 1, Quantity: 2}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(itemRepo.createdItems) != 1 || len(itemRepo.updatedItems) != 2 {
		t.Fatalf("expected inventory item updates")
	}
	if len(txRepo.transactions) != 1 {
		t.Fatalf("expected transaction record")
	}
}

func TestInventoryUseCaseDoTransactionDuplicatedSkus(t *testing.T) {
	uc := &inventoryUseCase{skuRepository: &doTxSkuRepository{skus: []domain.Sku{{Id: 1, Code: "A", Product: domain.Product{Name: "Prod"}}}}}
	err := uc.DoTransaction(context.Background(), nil, DoTransactionInput{Skus: []DoTransactionSkusInput{{SkuId: 1, Quantity: 1}, {SkuId: 1, Quantity: 1}}})
	if err == nil || !strings.Contains(err.Error(), ErrSkusDuplicated.Error()) {
		t.Fatalf("expected duplicated skus error")
	}
}

func TestInventoryUseCaseDoTransactionSkuNotFound(t *testing.T) {
	uc := &inventoryUseCase{skuRepository: &doTxSkuRepository{skus: []domain.Sku{{Id: 2, Code: "A", Product: domain.Product{Name: "Prod"}}}}}
	err := uc.DoTransaction(context.Background(), nil, DoTransactionInput{Skus: []DoTransactionSkusInput{{SkuId: 1, Quantity: 1}}})
	if err == nil || !strings.Contains(err.Error(), ErrSkusNotFound.Error()) {
		t.Fatalf("expected sku not found error")
	}
}

func TestInventoryUseCaseDoTransactionInventoryItemError(t *testing.T) {
	inventoryRepo := &doTxInventoryRepository{inventories: map[int64]domain.Inventory{1: {Id: 1}}}
	itemRepo := &stubInventoryItemRepository{getManyErr: errors.New("fail")}
	skuRepo := &doTxSkuRepository{skus: []domain.Sku{{Id: 1, Code: "A", Product: domain.Product{Name: "Prod"}}}}
	uc := &inventoryUseCase{inventoryRepository: inventoryRepo, inventoryItemRepository: itemRepo, skuRepository: skuRepo}

	err := uc.DoTransaction(context.Background(), nil, DoTransactionInput{InventoryOriginId: 1, Skus: []DoTransactionSkusInput{{SkuId: 1, Quantity: 1}}})
	if err == nil || err.Error() != "fail" {
		t.Fatalf("expected inventory item error")
	}
}

func TestInventoryUseCaseDoTransactionTypeIn(t *testing.T) {
	fakeTx := &fakeTx{}
	db := newFakeDB(fakeTx)
	repos := repository.NewRepository(db)
	tx, _ := db.BeginTx(context.Background(), nil)
	inventoryRepo := &doTxInventoryRepository{inventories: map[int64]domain.Inventory{2: {Id: 2}}}
	itemRepo := newDoTxInventoryItemRepository()
	txRepo := &doTxInventoryTransactionRepository{}
	skuRepo := &doTxSkuRepository{skus: []domain.Sku{{Id: 1, Code: "A", Product: domain.Product{Name: "Prod"}}}}

	uc := &inventoryUseCase{repository: repos, inventoryRepository: inventoryRepo, inventoryItemRepository: itemRepo, inventoryTransactionRepo: txRepo, skuRepository: skuRepo}

	err := uc.DoTransaction(context.Background(), tx, DoTransactionInput{Type: domain.InventoryTransactionTypeIn, InventoryDestinationId: 2, Skus: []DoTransactionSkusInput{{SkuId: 1, Quantity: 3}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(itemRepo.createdItems) != 1 {
		t.Fatalf("expected inventory item created for type in")
	}
}

func TestInventoryUseCaseDoTransactionTypeOut(t *testing.T) {
	fakeTx := &fakeTx{}
	db := newFakeDB(fakeTx)
	repos := repository.NewRepository(db)
	tx, _ := db.BeginTx(context.Background(), nil)
	inventoryRepo := &doTxInventoryRepository{inventories: map[int64]domain.Inventory{1: {Id: 1}}}
	itemRepo := newDoTxInventoryItemRepository()
	itemRepo.items[1] = []domain.InventoryItem{{Id: 1, InventoryId: 1, Sku: domain.Sku{Id: 1}, Quantity: 5}}
	txRepo := &doTxInventoryTransactionRepository{}
	skuRepo := &doTxSkuRepository{skus: []domain.Sku{{Id: 1, Code: "A", Product: domain.Product{Name: "Prod"}}}}

	uc := &inventoryUseCase{repository: repos, inventoryRepository: inventoryRepo, inventoryItemRepository: itemRepo, inventoryTransactionRepo: txRepo, skuRepository: skuRepo}

	err := uc.DoTransaction(context.Background(), tx, DoTransactionInput{Type: domain.InventoryTransactionTypeOut, InventoryOriginId: 1, Skus: []DoTransactionSkusInput{{SkuId: 1, Quantity: 2}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(itemRepo.updatedItems) != 1 {
		t.Fatalf("expected inventory item updated for type out")
	}
}

func TestInventoryUseCaseDoTransactionInventoryNotFound(t *testing.T) {
	inventoryRepo := &doTxInventoryRepository{inventories: map[int64]domain.Inventory{}}
	skuRepo := &doTxSkuRepository{skus: []domain.Sku{{Id: 1, Code: "A", Product: domain.Product{Name: "Prod"}}}}
	uc := &inventoryUseCase{inventoryRepository: inventoryRepo, skuRepository: skuRepo}
	err := uc.DoTransaction(context.Background(), nil, DoTransactionInput{InventoryOriginId: 1, Skus: []DoTransactionSkusInput{{SkuId: 1, Quantity: 1}}})
	if err == nil {
		t.Fatalf("expected error when inventory not found")
	}
}

func TestInventoryUseCaseDoTransactionCreateInventoryItemError(t *testing.T) {
	fakeTx := &fakeTx{}
	db := newFakeDB(fakeTx)
	repos := repository.NewRepository(db)
	tx, _ := db.BeginTx(context.Background(), nil)
	inventoryRepo := &doTxInventoryRepository{inventories: map[int64]domain.Inventory{2: {Id: 2}}}
	itemRepo := &stubInventoryItemRepository{createErr: errors.New("fail"), items: make(map[int64][]domain.InventoryItem)}
	txRepo := &doTxInventoryTransactionRepository{}
	skuRepo := &doTxSkuRepository{skus: []domain.Sku{{Id: 1, Code: "A", Product: domain.Product{Name: "Prod"}}}}
	uc := &inventoryUseCase{repository: repos, inventoryRepository: inventoryRepo, inventoryItemRepository: itemRepo, inventoryTransactionRepo: txRepo, skuRepository: skuRepo}

	err := uc.DoTransaction(context.Background(), tx, DoTransactionInput{Type: domain.InventoryTransactionTypeIn, InventoryDestinationId: 2, Skus: []DoTransactionSkusInput{{SkuId: 1, Quantity: 1}}})
	if err == nil || err.Error() != "fail" {
		t.Fatalf("expected error from inventory item creation")
	}
}

func TestInventoryUseCaseDoTransactionTransactionError(t *testing.T) {
	fakeTx := &fakeTx{}
	db := newFakeDB(fakeTx)
	repos := repository.NewRepository(db)
	tx, _ := db.BeginTx(context.Background(), nil)
	inventoryRepo := &doTxInventoryRepository{inventories: map[int64]domain.Inventory{1: {Id: 1}, 2: {Id: 2}}}
	itemRepo := newDoTxInventoryItemRepository()
	itemRepo.items[1] = []domain.InventoryItem{{Id: 1, InventoryId: 1, Sku: domain.Sku{Id: 1}, Quantity: 5}}
	txRepo := &stubInventoryTransactionRepository{err: errors.New("fail")}
	skuRepo := &doTxSkuRepository{skus: []domain.Sku{{Id: 1, Code: "A", Product: domain.Product{Name: "Prod"}}}}
	uc := &inventoryUseCase{repository: repos, inventoryRepository: inventoryRepo, inventoryItemRepository: itemRepo, inventoryTransactionRepo: txRepo, skuRepository: skuRepo}

	err := uc.DoTransaction(context.Background(), tx, DoTransactionInput{Type: domain.InventoryTransactionTypeTransfer, InventoryOriginId: 1, InventoryDestinationId: 2, Skus: []DoTransactionSkusInput{{SkuId: 1, Quantity: 1}}})
	if err == nil || err.Error() != "fail" {
		t.Fatalf("expected transaction error")
	}
}

func TestInventoryUseCaseDoTransactionUpdateQuantityError(t *testing.T) {
	fakeTx := &fakeTx{}
	db := newFakeDB(fakeTx)
	repos := repository.NewRepository(db)
	tx, _ := db.BeginTx(context.Background(), nil)
	inventoryRepo := &doTxInventoryRepository{inventories: map[int64]domain.Inventory{1: {Id: 1}, 2: {Id: 2}}}
	itemRepo := &stubInventoryItemRepository{updateErr: errors.New("fail"), items: map[int64][]domain.InventoryItem{1: {{Id: 1, InventoryId: 1, Sku: domain.Sku{Id: 1}, Quantity: 5}}, 2: {{Id: 2, InventoryId: 2, Sku: domain.Sku{Id: 1}, Quantity: 0}}}, returnedInventoryItem: domain.InventoryItem{Id: 2, InventoryId: 2, Sku: domain.Sku{Id: 1}, Quantity: 0}}
	txRepo := &doTxInventoryTransactionRepository{}
	skuRepo := &doTxSkuRepository{skus: []domain.Sku{{Id: 1, Code: "A", Product: domain.Product{Name: "Prod"}}}}
	uc := &inventoryUseCase{repository: repos, inventoryRepository: inventoryRepo, inventoryItemRepository: itemRepo, inventoryTransactionRepo: txRepo, skuRepository: skuRepo}

	err := uc.DoTransaction(context.Background(), tx, DoTransactionInput{Type: domain.InventoryTransactionTypeTransfer, InventoryOriginId: 1, InventoryDestinationId: 2, Skus: []DoTransactionSkusInput{{SkuId: 1, Quantity: 1}}})
	if err == nil || err.Error() != "fail" {
		t.Fatalf("expected update quantity error")
	}
}
