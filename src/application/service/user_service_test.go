package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
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

func TestUserServiceGetActiveLegalTermsSuccess(t *testing.T) {
	legalRepo := &stubUserLegalDocumentRepository{
		activeByUser: []domain.LegalTermStatus{
			{DocType: domain.LegalDocumentTypeTerms, DocVersion: "v1", Accepted: true},
		},
	}
	service := &userService{legalDocumentRepo: legalRepo}

	terms, err := service.GetActiveLegalTerms(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(terms) != 1 || terms[0].DocVersion != "v1" {
		t.Fatalf("unexpected legal terms: %+v", terms)
	}
}

func TestUserServiceGetActiveLegalTermsNotConfigured(t *testing.T) {
	service := &userService{}
	if _, err := service.GetActiveLegalTerms(context.Background(), 5); err == nil {
		t.Fatalf("expected not configured error")
	}
}

func TestNormalizePhone(t *testing.T) {
	if normalizePhone(nil) != nil {
		t.Fatalf("expected nil phone to stay nil")
	}

	empty := "   "
	if normalizePhone(&empty) != nil {
		t.Fatalf("expected empty phone to return nil")
	}

	value := "  11999998888 "
	normalized := normalizePhone(&value)
	if normalized == nil || *normalized != "11999998888" {
		t.Fatalf("unexpected normalized phone: %v", normalized)
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
func TestUserServiceResetPasswordTokenLookupError(t *testing.T) {
	service, _, _, tokenRepo := newUserServiceTest(&stubUserRepository{}, &stubInventoryRepository{})
	tokenRepo.lastActiveErr = errors.New("lookup fail")

	err := service.ResetPassword(context.Background(), request.ResetPasswordRequest{Code: "reset", Uuid: "uuid", Password: "newpass123"})
	if err == nil || err.Error() != "Código expirado ou inválido! Solicite um novo código para o administrador." {
		t.Fatalf("expected expired code error, got %v", err)
	}
}

func TestUserServiceResetPasswordCompareError(t *testing.T) {
	service, _, _, tokenRepo := newUserServiceTest(&stubUserRepository{}, &stubInventoryRepository{})
	tokenRepo.lastActive = domain.UserToken{
		User:      domain.User{Id: 10},
		CodeHash:  "encrypted:reset",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	service.encrypto = &stubEncrypto{compareErr: errors.New("compare fail")}

	err := service.ResetPassword(context.Background(), request.ResetPasswordRequest{Code: "reset", Uuid: "uuid", Password: "newpass123"})
	if err == nil || err.Error() != "Código expirado ou inválido! Solicite um novo código para o administrador." {
		t.Fatalf("expected expired code error, got %v", err)
	}
}

func TestUserServiceResetPasswordEncryptError(t *testing.T) {
	userRepo := &stubUserRepository{}
	service, _, _, tokenRepo := newUserServiceTest(userRepo, &stubInventoryRepository{})
	tokenRepo.lastActive = domain.UserToken{
		User:      domain.User{Id: 10},
		CodeHash:  "encrypted:reset",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	service.encrypto = &stubEncrypto{encryptErr: errors.New("encrypt fail")}

	err := service.ResetPassword(context.Background(), request.ResetPasswordRequest{Code: "reset", Uuid: "uuid", Password: "newpass123"})
	if err == nil || err.Error() != "encrypt fail" {
		t.Fatalf("expected encrypt error, got %v", err)
	}
}

func TestUserServiceResetPasswordUpdatePasswordError(t *testing.T) {
	userRepo := &stubUserRepository{updatePasswordErr: errors.New("update fail")}
	service, _, _, tokenRepo := newUserServiceTest(userRepo, &stubInventoryRepository{})
	tokenRepo.lastActive = domain.UserToken{
		User:      domain.User{Id: 10},
		CodeHash:  "encrypted:reset",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err := service.ResetPassword(context.Background(), request.ResetPasswordRequest{Code: "reset", Uuid: "uuid", Password: "newpass123"})
	if err == nil || err.Error() != "update fail" {
		t.Fatalf("expected update error, got %v", err)
	}
}

func TestUserServiceResetPasswordSetUsedError(t *testing.T) {
	userRepo := &stubUserRepository{}
	service, _, _, tokenRepo := newUserServiceTest(userRepo, &stubInventoryRepository{})
	tokenRepo.lastActive = domain.UserToken{
		User:      domain.User{Id: 10},
		CodeHash:  "encrypted:reset",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	tokenRepo.setUsedErr = errors.New("set used fail")

	err := service.ResetPassword(context.Background(), request.ResetPasswordRequest{Code: "reset", Uuid: "uuid", Password: "newpass123"})
	if err == nil || err.Error() != "set used fail" {
		t.Fatalf("expected set used error, got %v", err)
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
func TestUserServiceForgotPasswordTokenCreateError(t *testing.T) {
	user := domain.User{Id: 5, Email: "user@test.com", Name: "User"}
	userRepo := &stubUserRepository{getByEmail: user}
	service, _, emailUsecase, tokenRepo := newUserServiceTest(userRepo, &stubInventoryRepository{})
	tokenRepo.createErr = errors.New("create fail")

	if err := service.ForgotPassword(context.Background(), request.ForgotPasswordRequest{Email: user.Email}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if tokenRepo.createInput.User.Id != 0 || emailUsecase.recoverCalls != 0 {
		t.Fatalf("expected no token creation or email send")
	}
}

func TestUserServiceForgotPasswordGetTokenError(t *testing.T) {
	user := domain.User{Id: 5, Email: "user@test.com", Name: "User"}
	userRepo := &stubUserRepository{getByEmail: user}
	service, _, emailUsecase, tokenRepo := newUserServiceTest(userRepo, &stubInventoryRepository{})
	tokenRepo.createID = 1
	tokenRepo.getByIdErr = errors.New("get token fail")

	if err := service.ForgotPassword(context.Background(), request.ForgotPasswordRequest{Email: user.Email}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if emailUsecase.recoverCalls != 0 {
		t.Fatalf("expected no recover email on token error")
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

func TestUserServiceAcceptLegalTermsSuccess(t *testing.T) {
	tx, fakeTx, cleanup := newTestSQLTx()
	defer cleanup()

	docRepo := &stubAcceptLegalDocumentRepository{
		documents: map[string]domain.LegalDocument{
			"TERMS:v1":   {Id: 1, DocType: domain.LegalDocumentTypeTerms, DocVersion: "v1"},
			"PRIVACY:v2": {Id: 2, DocType: domain.LegalDocumentTypePrivacy, DocVersion: "v2"},
		},
	}
	acceptanceRepo := &stubAcceptLegalAcceptanceRepository{}
	service := &userService{
		legalDocumentRepo:   docRepo,
		legalAcceptanceRepo: acceptanceRepo,
		txManager:           &stubTxManager{tx: tx},
	}

	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(10))
	err := service.AcceptLegalTerms(ctx, 99, []request.AcceptLegalTermRequest{
		{DocType: "TERMS", DocVersion: "v1", Accepted: true},
		{DocType: "PRIVACY", DocVersion: "v2", Accepted: true},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if len(acceptanceRepo.created) != 2 {
		t.Fatalf("expected 2 acceptances, got %d", len(acceptanceRepo.created))
	}
	if !fakeTx.committed || fakeTx.rolledBack {
		t.Fatalf("transaction not committed correctly: %+v", fakeTx)
	}
}

func TestUserServiceAcceptLegalTermsValidationError(t *testing.T) {
	service := &userService{
		legalDocumentRepo:   &stubAcceptLegalDocumentRepository{},
		legalAcceptanceRepo: &stubAcceptLegalAcceptanceRepository{},
		txManager:           &stubTxManager{tx: &sql.Tx{}},
	}

	err := service.AcceptLegalTerms(context.Background(), 1, nil)
	if err == nil || err.Error() != "Envie ao menos um termo." {
		t.Fatalf("expected empty terms error, got %v", err)
	}

	err = service.AcceptLegalTerms(context.Background(), 1, []request.AcceptLegalTermRequest{
		{DocType: "TERMS", DocVersion: "v1", Accepted: false},
	})
	if err == nil || !strings.Contains(err.Error(), "Accepted") {
		t.Fatalf("expected accept validation error, got %v", err)
	}
}

func TestUserServiceAcceptLegalTermsDuplicate(t *testing.T) {
	tx, _, cleanup := newTestSQLTx()
	defer cleanup()

	docRepo := &stubAcceptLegalDocumentRepository{
		documents: map[string]domain.LegalDocument{
			"TERMS:v1": {Id: 1, DocType: domain.LegalDocumentTypeTerms, DocVersion: "v1"},
		},
	}
	acceptanceRepo := &stubAcceptLegalAcceptanceRepository{
		err: errors.New("duplicate key value violates unique constraint"),
	}
	service := &userService{
		legalDocumentRepo:   docRepo,
		legalAcceptanceRepo: acceptanceRepo,
		txManager:           &stubTxManager{tx: tx},
	}

	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(10))
	err := service.AcceptLegalTerms(ctx, 1, []request.AcceptLegalTermRequest{
		{DocType: "TERMS", DocVersion: "v1", Accepted: true},
	})
	if err == nil || !strings.Contains(err.Error(), "Termo") {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestUserServiceAcceptLegalTermsInvalidTenant(t *testing.T) {
	tx, _, cleanup := newTestSQLTx()
	defer cleanup()

	service := &userService{
		legalDocumentRepo:   &stubAcceptLegalDocumentRepository{},
		legalAcceptanceRepo: &stubAcceptLegalAcceptanceRepository{},
		txManager:           &stubTxManager{tx: tx},
	}

	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, "invalid")
	err := service.AcceptLegalTerms(ctx, 1, []request.AcceptLegalTermRequest{
		{DocType: "TERMS", DocVersion: "v1", Accepted: true},
	})
	if err == nil || err.Error() != "tenant id invalido" {
		t.Fatalf("expected tenant error, got %v", err)
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

type stubAcceptLegalDocumentRepository struct {
	documents map[string]domain.LegalDocument
	err       error
}

func (s *stubAcceptLegalDocumentRepository) GetLastActiveByType(ctx context.Context, docType domain.LegalDocumentType) (domain.LegalDocument, error) {
	return domain.LegalDocument{}, nil
}

func (s *stubAcceptLegalDocumentRepository) GetByTypeAndVersion(ctx context.Context, docType domain.LegalDocumentType, version string) (domain.LegalDocument, error) {
	if s.err != nil {
		return domain.LegalDocument{}, s.err
	}
	key := string(docType) + ":" + version
	if doc, ok := s.documents[key]; ok {
		return doc, nil
	}
	return domain.LegalDocument{}, errors.New("document not found")
}

func (s *stubAcceptLegalDocumentRepository) GetActiveByUser(ctx context.Context, userId int64) ([]domain.LegalTermStatus, error) {
	return nil, nil
}

type stubAcceptLegalAcceptanceRepository struct {
	err     error
	created []domain.LegalAcceptance
}

func (s *stubAcceptLegalAcceptanceRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, acceptance domain.LegalAcceptance) (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	s.created = append(s.created, acceptance)
	return int64(len(s.created)), nil
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
