package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/logs"
)

func init() {
	if logs.Logger == nil {
		logs.Logger = stubLogs{}
	}
}

func TestUserServiceCreateReseller(t *testing.T) {
	userRepo := &stubUserRepository{
		getByIdResponses: map[int64]domain.User{
			1:  {Id: 1, Email: "user@test.com", Name: "User"},
			99: {Id: 99, Name: "Admin"},
		},
	}
	inventoryRepo := &stubInventoryRepository{}
	service, _, _, _ := newUserServiceTest(userRepo, inventoryRepo)

	req := request.CreateUserRequest{Username: "user", Name: "User", Role: string(domain.InventoryTypeReseller), Email: "user@test.com"}
	if err := service.Create(adminContext(99), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inventoryRepo.created.Type != domain.InventoryTypeReseller {
		t.Fatalf("expected reseller inventory to be created")
	}
}

func TestUserServiceCreateDuplicated(t *testing.T) {
	userRepo := &stubUserRepository{createErr: errors.New("duplicate key value violates unique constraint")}
	service, _, _, _ := newUserServiceTest(userRepo, &stubInventoryRepository{})

	req := request.CreateUserRequest{Username: "user", Name: "User", Role: "ADMIN", Email: "user@test.com"}
	err := service.Create(adminContext(99), req)
	if err == nil || err.Error() != "Usuário já cadastrado!" {
		t.Fatalf("expected duplicated error")
	}
}

func TestUserServiceCreateRepositoryError(t *testing.T) {
	userRepo := &stubUserRepository{createErr: errors.New("fail")}
	service, _, _, _ := newUserServiceTest(userRepo, &stubInventoryRepository{})
	req := request.CreateUserRequest{Username: "user", Name: "User", Role: "ADMIN", Email: "user@test.com"}
	if err := service.Create(adminContext(99), req); err == nil || err.Error() != "fail" {
		t.Fatalf("expected repository error")
	}
}

func TestUserServiceUpdateResellerCreatesInventory(t *testing.T) {
	userRepo := &stubUserRepository{}
	inventoryRepo := &stubInventoryRepository{getByUserErr: domain.ErrInventoryNotFound}
	service, _, _, _ := newUserServiceTest(userRepo, inventoryRepo)

	req := request.EditUserRequest{Username: "user", Name: "User", Role: string(domain.InventoryTypeReseller), Email: "user@test.com"}
	if err := service.Update(context.Background(), req, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inventoryRepo.created.Type != domain.InventoryTypeReseller {
		t.Fatalf("expected inventory creation on update")
	}
}

func TestUserServiceUpdateDuplicated(t *testing.T) {
	userRepo := &stubUserRepository{updateErr: errors.New("duplicate key value violates unique constraint")}
	service, _, _, _ := newUserServiceTest(userRepo, &stubInventoryRepository{})
	req := request.EditUserRequest{Username: "user", Name: "User", Role: "ADMIN", Email: "user@test.com"}

	err := service.Update(context.Background(), req, 1)
	if err == nil || err.Error() != "Usuário já cadastrado!" {
		t.Fatalf("expected duplicated error")
	}
}

func TestUserServiceUpdateExistingInventory(t *testing.T) {
	inventoryRepo := &stubInventoryRepository{getByUser: domain.Inventory{Id: 1}}
	service, _, _, _ := newUserServiceTest(&stubUserRepository{}, inventoryRepo)
	req := request.EditUserRequest{Username: "user", Name: "User", Role: string(domain.InventoryTypeReseller), Email: "user@test.com"}
	if err := service.Update(context.Background(), req, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inventoryRepo.created.Id != 0 {
		t.Fatalf("expected no new inventory creation")
	}
}

func TestUserServiceCreateNonReseller(t *testing.T) {
	userRepo := &stubUserRepository{
		getByIdResponses: map[int64]domain.User{
			1:  {Id: 1, Email: "user@test.com"},
			42: {Id: 42, Name: "Admin"},
		},
	}
	inventoryRepo := &stubInventoryRepository{}
	service, _, _, _ := newUserServiceTest(userRepo, inventoryRepo)
	req := request.CreateUserRequest{Username: "user", Name: "User", Role: "ADMIN", Email: "user@test.com"}
	if err := service.Create(adminContext(42), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inventoryRepo.created.Id != 0 {
		t.Fatalf("expected no inventory creation")
	}
}

func TestUserServiceGetters(t *testing.T) {
	userRepo := &stubUserRepository{getById: domain.User{Id: 1}, getAll: []domain.User{{Id: 2}}}
	service, _, _, _ := newUserServiceTest(userRepo, &stubInventoryRepository{})

	user, err := service.GetById(context.Background(), 1)
	if err != nil || user.Id != 1 {
		t.Fatalf("unexpected get by id result")
	}

	users, err := service.GetAll(context.Background(), request.GetAllUserRequest{})
	if err != nil || len(users) != 1 {
		t.Fatalf("unexpected get all result")
	}
}

func TestUserServiceGetByIdError(t *testing.T) {
	userRepo := &stubUserRepository{getByIdErr: errors.New("fail")}
	service, _, _, _ := newUserServiceTest(userRepo, &stubInventoryRepository{})
	if _, err := service.GetById(context.Background(), 1); err == nil || err.Error() != "fail" {
		t.Fatalf("expected repository error")
	}
}

func TestUserServiceGetAllWithRoleFilter(t *testing.T) {
	filterRole := domain.Role("ADMIN")
	userRepo := &stubUserRepository{getAll: []domain.User{{Id: 1}}}
	service, _, _, _ := newUserServiceTest(userRepo, &stubInventoryRepository{})

	if _, err := service.GetAll(context.Background(), request.GetAllUserRequest{Role: filterRole}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userRepo.getAllInput.Role == nil || *userRepo.getAllInput.Role != filterRole {
		t.Fatalf("expected role filter to be forwarded, got %+v", userRepo.getAllInput.Role)
	}
}

func TestUserServiceGetAllError(t *testing.T) {
	userRepo := &stubUserRepository{getAllErr: errors.New("fail")}
	service, _, _, _ := newUserServiceTest(userRepo, &stubInventoryRepository{})
	if _, err := service.GetAll(context.Background(), request.GetAllUserRequest{}); err == nil || err.Error() != "fail" {
		t.Fatalf("expected repository error")
	}
}

func TestUserServiceInactivate(t *testing.T) {
	userRepo := &stubUserRepository{}
	service, _, _, _ := newUserServiceTest(userRepo, &stubInventoryRepository{})

	if err := service.Inactivate(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserServiceCreateValidationError(t *testing.T) {
	service := &userService{}
	if err := service.Create(context.Background(), request.CreateUserRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestUserServiceCreateInventoryError(t *testing.T) {
	inventoryRepo := &stubInventoryRepository{createErr: errors.New("fail")}
	userRepo := &stubUserRepository{
		getByIdResponses: map[int64]domain.User{
			1:  {Id: 1, Email: "user@test.com"},
			77: {Id: 77},
		},
	}
	service, _, _, _ := newUserServiceTest(userRepo, inventoryRepo)
	req := request.CreateUserRequest{Username: "user", Name: "User", Role: string(domain.InventoryTypeReseller), Email: "user@test.com"}

	if err := service.Create(adminContext(77), req); err == nil || err.Error() != "fail" {
		t.Fatalf("expected inventory error")
	}
}

func TestUserServiceUpdateValidationError(t *testing.T) {
	service := &userService{}
	if err := service.Update(context.Background(), request.EditUserRequest{}, 1); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestUserServiceUpdateInventoryError(t *testing.T) {
	inventoryRepo := &stubInventoryRepository{getByUserErr: errors.New("fail")}
	service, _, _, _ := newUserServiceTest(&stubUserRepository{}, inventoryRepo)
	req := request.EditUserRequest{Username: "user", Name: "User", Role: string(domain.InventoryTypeReseller), Email: "user@test.com"}

	if err := service.Update(context.Background(), req, 1); err == nil || err.Error() != "fail" {
		t.Fatalf("expected inventory error")
	}
}

func TestUserServiceResetPasswordSuccess(t *testing.T) {
	userRepo := &stubUserRepository{}
	service, _, _, tokenRepo := newUserServiceTest(userRepo, &stubInventoryRepository{})
	tokenRepo.lastActive = domain.UserToken{
		User:      domain.User{Id: 10},
		CodeHash:  "encrypted:reset",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	req := request.ResetPasswordRequest{Code: "reset", Uuid: "uuid", Password: "newpass123"}
	if err := service.ResetPassword(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userRepo.updatePasswordReq.user.Id != 10 {
		t.Fatalf("expected password update for user 10, got %d", userRepo.updatePasswordReq.user.Id)
	}
	if userRepo.updatePasswordReq.newPassword == "" {
		t.Fatalf("expected encrypted password to be set")
	}
	if tokenRepo.setUsedToken.User.Id != 10 {
		t.Fatalf("expected token to be marked as used")
	}
}

func TestUserServiceResetPasswordValidationError(t *testing.T) {
	service, _, _, _ := newUserServiceTest(&stubUserRepository{}, &stubInventoryRepository{})
	if err := service.ResetPassword(context.Background(), request.ResetPasswordRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestUserServiceResetPasswordInvalidToken(t *testing.T) {
	service, _, _, tokenRepo := newUserServiceTest(&stubUserRepository{}, &stubInventoryRepository{})
	tokenRepo.lastActive = domain.UserToken{
		User:      domain.User{Id: 10},
		CodeHash:  "encrypted:different",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err := service.ResetPassword(context.Background(), request.ResetPasswordRequest{Code: "reset", Uuid: "uuid", Password: "newpass123"})
	if err == nil || err.Error() != "Código expirado ou inválido! Solicite um novo código para o administrador." {
		t.Fatalf("expected invalid token error, got %v", err)
	}
}

func TestUserServiceForgotPasswordValidationError(t *testing.T) {
	service, _, _, _ := newUserServiceTest(&stubUserRepository{}, &stubInventoryRepository{})
	if err := service.ForgotPassword(context.Background(), request.ForgotPasswordRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestUserServiceForgotPasswordLookupError(t *testing.T) {
	userRepo := &stubUserRepository{getByEmailErr: errors.New("fail")}
	service, _, emailUsecase, tokenRepo := newUserServiceTest(userRepo, &stubInventoryRepository{})
	if err := service.ForgotPassword(context.Background(), request.ForgotPasswordRequest{Email: "user@test.com"}); err != nil {
		t.Fatalf("expected nil error even when lookup fails, got %v", err)
	}
	if tokenRepo.createInput.User.Id != 0 || emailUsecase.recoverCalls != 0 {
		t.Fatalf("expected no token creation or email send")
	}
}

func TestUserServiceForgotPasswordSuccess(t *testing.T) {
	user := domain.User{Id: 5, Email: "user@test.com", Name: "User"}
	userRepo := &stubUserRepository{getByEmail: user}
	service, _, emailUsecase, tokenRepo := newUserServiceTest(userRepo, &stubInventoryRepository{})
	tokenRepo.createID = 1
	tokenRepo.getById = domain.UserToken{Uuid: "uuid-123"}
	emailUsecase.recoverCh = make(chan struct{}, 1)
	if err := service.ForgotPassword(context.Background(), request.ForgotPasswordRequest{Email: user.Email}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	select {
	case <-emailUsecase.recoverCh:
	case <-time.After(50 * time.Millisecond):
		t.Fatalf("expected recover email to be sent")
	}
	if tokenRepo.createInput.Type != domain.UserTokenTypeResetPass {
		t.Fatalf("expected reset pass token, got %s", tokenRepo.createInput.Type)
	}
	if emailUsecase.recoverCalls != 1 {
		t.Fatalf("expected recover email to be sent")
	}
	if emailUsecase.lastRecover.user.Id != user.Id || emailUsecase.lastRecover.uuid != "uuid-123" {
		t.Fatalf("unexpected recover payload: %+v", emailUsecase.lastRecover)
	}
	if emailUsecase.lastRecover.code == "" {
		t.Fatalf("expected non-empty recovery code")
	}
}

type stubUserLegalDocumentRepository struct {
	activeByUser []domain.LegalTermStatus
	err          error
}

func (s *stubUserLegalDocumentRepository) GetLastActiveByType(ctx context.Context, docType domain.LegalDocumentType) (domain.LegalDocument, error) {
	return domain.LegalDocument{}, nil
}

func (s *stubUserLegalDocumentRepository) GetByTypeAndVersion(ctx context.Context, docType domain.LegalDocumentType, version string) (domain.LegalDocument, error) {
	return domain.LegalDocument{}, nil
}

func (s *stubUserLegalDocumentRepository) GetActiveByUser(ctx context.Context, userId int64) ([]domain.LegalTermStatus, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.activeByUser, nil
}

type stubUserLegalAcceptanceRepository struct{}

func (s *stubUserLegalAcceptanceRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, acceptance domain.LegalAcceptance) (int64, error) {
	return 1, nil
}

type stubUserTxManager struct {
	tx  *sql.Tx
	err error
}

func (s *stubUserTxManager) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.tx, s.err
}

func newUserServiceTest(userRepo *stubUserRepository, inventoryRepo *stubInventoryRepository) (*userService, *stubUserTokenService, *stubEmailUseCase, *stubUserTokenRepository) {
	if userRepo == nil {
		userRepo = &stubUserRepository{}
	}
	if inventoryRepo == nil {
		inventoryRepo = &stubInventoryRepository{}
	}
	userTokenService := &stubUserTokenService{output: domain.UserToken{Code: "code", Uuid: "uuid"}}
	emailUsecase := &stubEmailUseCase{}
	userTokenRepository := &stubUserTokenRepository{}
	legalDocumentRepo := &stubUserLegalDocumentRepository{}
	legalAcceptanceRepo := &stubUserLegalAcceptanceRepository{}
	service := &userService{
		userRepository:      userRepo,
		inventoryRepository: inventoryRepo,
		encrypto:            &stubEncrypto{},
		userTokenService:    userTokenService,
		emailUsecase:        emailUsecase,
		userTokenRepository: userTokenRepository,
		legalDocumentRepo:   legalDocumentRepo,
		legalAcceptanceRepo: legalAcceptanceRepo,
		txManager:           &stubUserTxManager{},
	}
	return service, userTokenService, emailUsecase, userTokenRepository
}

func adminContext(id float64) context.Context {
	return context.WithValue(context.Background(), constants.USERID_KEY, id)
}

type stubLogs struct{}

func (stubLogs) Infof(string, ...interface{})  {}
func (stubLogs) Printf(string, ...interface{}) {}
func (stubLogs) Warnf(string, ...interface{})  {}
func (stubLogs) Errorf(string, ...interface{}) {}
func (stubLogs) Fatalf(string, ...interface{}) {}
func (stubLogs) Panicf(string, ...interface{}) {}
func (stubLogs) AddHook(logs.Hook)             {}
func (stubLogs) With(map[string]any) logs.Logs { return stubLogs{} }
