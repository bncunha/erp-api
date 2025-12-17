package service

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/lib/pq"
)

type stubCompanyRepository struct {
	id  int64
	err error
}

func (s *stubCompanyRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, company domain.Company) (int64, error) {
	return s.id, s.err
}

type stubAddressRepository struct {
	err error
}

func (s *stubAddressRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, address domain.Address) (int64, error) {
	return 1, s.err
}

type stubCompanyInventoryRepository struct {
	input domain.Inventory
	err   error
}

func (s *stubCompanyInventoryRepository) Create(ctx context.Context, inventory domain.Inventory) (int64, error) {
	return 0, nil
}

func (s *stubCompanyInventoryRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, inventory domain.Inventory) (int64, error) {
	s.input = inventory
	return 1, s.err
}

func (s *stubCompanyInventoryRepository) GetById(ctx context.Context, id int64) (domain.Inventory, error) {
	return domain.Inventory{}, nil
}

func (s *stubCompanyInventoryRepository) GetAll(ctx context.Context) ([]domain.Inventory, error) {
	return nil, nil
}

func (s *stubCompanyInventoryRepository) GetByUserId(ctx context.Context, userId int64) (domain.Inventory, error) {
	return domain.Inventory{}, nil
}

func (s *stubCompanyInventoryRepository) GetPrimaryInventory(ctx context.Context) (domain.Inventory, error) {
	return domain.Inventory{}, nil
}

func (s *stubCompanyInventoryRepository) GetSummary(ctx context.Context) ([]output.GetInventorySummaryOutput, error) {
	return nil, nil
}

func (s *stubCompanyInventoryRepository) GetSummaryById(ctx context.Context, id int64) (output.GetInventorySummaryByIdOutput, error) {
	return output.GetInventorySummaryByIdOutput{}, nil
}

type stubCompanyUserRepository struct {
	createdUser domain.User
	id          int64
	err         error
}

func (s *stubCompanyUserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	return domain.User{}, nil
}

func (s *stubCompanyUserRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	s.createdUser = user
	if s.id == 0 {
		s.id = 10
	}
	return s.id, s.err
}

func (s *stubCompanyUserRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, user domain.User) (int64, error) {
	return s.Create(ctx, user)
}

func (s *stubCompanyUserRepository) Update(ctx context.Context, user domain.User) error { return nil }
func (s *stubCompanyUserRepository) Inactivate(ctx context.Context, id int64) error     { return nil }
func (s *stubCompanyUserRepository) GetAll(ctx context.Context, input domain.GetAllUserInput) ([]domain.User, error) {
	return nil, nil
}
func (s *stubCompanyUserRepository) GetById(ctx context.Context, id int64) (domain.User, error) {
	return domain.User{}, nil
}
func (s *stubCompanyUserRepository) UpdatePassword(ctx context.Context, user domain.User, newPassword string) error {
	return nil
}
func (s *stubCompanyUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return domain.User{}, nil
}

type stubCompanyTxManager struct {
	tx  *sql.Tx
	err error
}

func (s *stubCompanyTxManager) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.tx, s.err
}

type stubWelcomeEmailUseCase struct {
	toEmail string
	toName  string
	err     error
}

func (s *stubWelcomeEmailUseCase) SendInvite(ctx context.Context, user domain.User, code string, uuid string) error {
	return nil
}

func (s *stubWelcomeEmailUseCase) SendRecoverPassword(ctx context.Context, user domain.User, code string, uuid string) error {
	return nil
}

func (s *stubWelcomeEmailUseCase) SendWelcome(ctx context.Context, email string, name string) error {
	s.toEmail = email
	s.toName = name
	return s.err
}

func strPtr(s string) *string { return &s }

func newFakeTransaction() (*sql.Tx, *fakeSQLTx) {
	driverName := fmt.Sprintf("fakedriver-%d", atomic.AddInt64(&fakeDriverCounter, 1))
	fake := &fakeSQLTx{}
	sql.Register(driverName, &fakeDriver{tx: fake})
	db, _ := sql.Open(driverName, "")
	tx, _ := db.Begin()
	return tx, fake
}

func TestCompanyServiceCreateSuccess(t *testing.T) {
	tx, fakeTx := newFakeTransaction()
	userRepo := &stubCompanyUserRepository{}
	inventoryRepo := &stubCompanyInventoryRepository{}
	emailUsecase := &stubWelcomeEmailUseCase{}
	service := NewCompanyService(&stubCompanyRepository{id: 5}, &stubAddressRepository{}, inventoryRepo, userRepo, &stubEncrypto{}, emailUsecase, &stubCompanyTxManager{tx: tx})

	err := service.Create(context.Background(), request.CreateCompanyRequest{
		Name:      "Empresa",
		LegalName: "Empresa Ltda",
		Cpf:       "390.533.447-05",
		Cellphone: "11999999999",
		Address: request.CreateCompanyAddress{
			Street:       "Rua A",
			Neighborhood: "Centro",
			Number:       "123",
			City:         "Cidade",
			UF:           "SP",
			Cep:          "00000000",
		},
		User: request.CreateCompanyUserRequest{
			Name:        "Admin",
			Username:    "admin",
			PhoneNumber: strPtr("123"),
			Email:       "admin@test.com",
			Password:    "secret123",
		},
	})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	time.Sleep(10 * time.Millisecond)
	if !fakeTx.committed || fakeTx.rolledBack {
		t.Fatalf("transaction not committed correctly: %+v", fakeTx)
	}
	if userRepo.createdUser.TenantId != 0 || userRepo.createdUser.Username != "admin" {
		t.Fatalf("expected admin user to be created, got %+v", userRepo.createdUser)
	}
	if userRepo.createdUser.Password != "encrypted:secret123" {
		t.Fatalf("expected password to be encrypted, got %s", userRepo.createdUser.Password)
	}
	if inventoryRepo.input.User.Id != userRepo.id || inventoryRepo.input.Type != domain.InventoryTypePrimary {
		t.Fatalf("expected primary inventory to be created for user %d", userRepo.id)
	}
	if emailUsecase.toEmail != "admin@test.com" || emailUsecase.toName != "Admin" {
		t.Fatalf("expected welcome email to be triggered")
	}
}

func TestCompanyServiceCreateDuplicate(t *testing.T) {
	tx, fakeTx := newFakeTransaction()
	service := NewCompanyService(&stubCompanyRepository{err: fmt.Errorf("duplicate key value violates unique constraint")}, &stubAddressRepository{}, &stubCompanyInventoryRepository{}, &stubCompanyUserRepository{}, &stubEncrypto{}, &stubWelcomeEmailUseCase{}, &stubCompanyTxManager{tx: tx})

	err := service.Create(context.Background(), request.CreateCompanyRequest{
		Name:      "Empresa",
		LegalName: "Empresa Ltda",
		Cnpj:      "04.252.011/0001-10",
		Cellphone: "11999999999",
		Address:   request.CreateCompanyAddress{Street: "Rua", Neighborhood: "Centro", Number: "1", City: "Cidade", UF: "SP", Cep: "00000000"},
		User:      request.CreateCompanyUserRequest{Name: "Admin", Username: "admin", Email: "admin@test.com", Password: "secret123"},
	})

	if err == nil || err.Error() != "Empresa j\u00e1 cadastrada" {
		t.Fatalf("expected duplicate error, got %v", err)
	}
	if !fakeTx.rolledBack {
		t.Fatalf("expected transaction rollback")
	}
}

func TestCompanyServiceCreateDuplicateCNPJ(t *testing.T) {
	tx, fakeTx := newFakeTransaction()
	service := NewCompanyService(
		&stubCompanyRepository{err: &pq.Error{Message: "duplicate key value violates unique constraint", Detail: "Key (cnpj)=(04.252.011/0001-10) already exists.", Constraint: "companies_cnpj_unique"}},
		&stubAddressRepository{},
		&stubCompanyInventoryRepository{},
		&stubCompanyUserRepository{},
		&stubEncrypto{},
		&stubWelcomeEmailUseCase{},
		&stubCompanyTxManager{tx: tx},
	)

	err := service.Create(context.Background(), request.CreateCompanyRequest{
		Name:      "Empresa",
		LegalName: "Empresa Ltda",
		Cnpj:      "04.252.011/0001-10",
		Cellphone: "11999999999",
		Address:   request.CreateCompanyAddress{Street: "Rua", Neighborhood: "Centro", Number: "1", City: "Cidade", UF: "SP", Cep: "00000000"},
		User:      request.CreateCompanyUserRequest{Name: "Admin", Username: "admin", Email: "admin@test.com", Password: "secret123"},
	})

	if err == nil || err.Error() != "CNPJ j\u00e1 cadastrado" {
		t.Fatalf("expected CNPJ duplicate error, got %v", err)
	}
	if !fakeTx.rolledBack {
		t.Fatalf("expected transaction rollback")
	}
}

func TestCompanyServiceCreateDuplicateCPF(t *testing.T) {
	tx, fakeTx := newFakeTransaction()
	service := NewCompanyService(
		&stubCompanyRepository{err: &pq.Error{Message: "duplicate key value violates unique constraint", Detail: "Key (cpf)=(39053344705) already exists.", Constraint: "companies_cpf_unique"}},
		&stubAddressRepository{},
		&stubCompanyInventoryRepository{},
		&stubCompanyUserRepository{},
		&stubEncrypto{},
		&stubWelcomeEmailUseCase{},
		&stubCompanyTxManager{tx: tx},
	)

	err := service.Create(context.Background(), request.CreateCompanyRequest{
		Name:      "Empresa",
		LegalName: "Empresa Ltda",
		Cpf:       "39053344705",
		Cellphone: "11999999999",
		Address:   request.CreateCompanyAddress{Street: "Rua", Neighborhood: "Centro", Number: "1", City: "Cidade", UF: "SP", Cep: "00000000"},
		User:      request.CreateCompanyUserRequest{Name: "Admin", Username: "admin", Email: "admin@test.com", Password: "secret123"},
	})

	if err == nil || err.Error() != "CPF j\u00e1 cadastrado" {
		t.Fatalf("expected CPF duplicate error, got %v", err)
	}
	if !fakeTx.rolledBack {
		t.Fatalf("expected transaction rollback")
	}
}

func TestCompanyServiceCreateUserError(t *testing.T) {
	tx, fakeTx := newFakeTransaction()
	userRepo := &stubCompanyUserRepository{err: fmt.Errorf("user error")}
	service := NewCompanyService(&stubCompanyRepository{id: 2}, &stubAddressRepository{}, &stubCompanyInventoryRepository{}, userRepo, &stubEncrypto{}, &stubWelcomeEmailUseCase{}, &stubCompanyTxManager{tx: tx})

	err := service.Create(context.Background(), request.CreateCompanyRequest{
		Name:      "Empresa",
		LegalName: "Empresa Ltda",
		Cpf:       "39053344705",
		Cellphone: "11999999999",
		Address:   request.CreateCompanyAddress{Street: "Rua", Neighborhood: "Centro", Number: "1", City: "Cidade", UF: "SP", Cep: "00000000"},
		User:      request.CreateCompanyUserRequest{Name: "Admin", Username: "admin", Email: "admin@test.com", Password: "secret123"},
	})

	if err == nil || err.Error() != "user error" {
		t.Fatalf("expected user error, got %v", err)
	}
	if !fakeTx.rolledBack {
		t.Fatalf("expected rollback on user error")
	}
}
