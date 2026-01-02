package service

import (
	"context"
	"strings"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	emailusecase "github.com/bncunha/erp-api/src/application/usecase/email_usecase"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/logs"
	"github.com/lib/pq"
)

type CompanyService interface {
	Create(ctx context.Context, request request.CreateCompanyRequest) error
}

type companyService struct {
	companyRepository   domain.CompanyRepository
	addressRepository   domain.AddressRepository
	inventoryRepo       domain.InventoryRepository
	userRepository      domain.UserRepository
	encrypto            domain.Encrypto
	emailUsecase        emailusecase.EmailUseCase
	legalDocumentRepo   domain.LegalDocumentRepository
	legalAcceptanceRepo domain.LegalAcceptanceRepository
	txManager           transactionManager
}

func NewCompanyService(companyRepository domain.CompanyRepository, addressRepository domain.AddressRepository, inventoryRepo domain.InventoryRepository, userRepository domain.UserRepository, encrypto domain.Encrypto, emailUsecase emailusecase.EmailUseCase, legalDocumentRepo domain.LegalDocumentRepository, legalAcceptanceRepo domain.LegalAcceptanceRepository, txManager transactionManager) CompanyService {
	return &companyService{companyRepository, addressRepository, inventoryRepo, userRepository, encrypto, emailUsecase, legalDocumentRepo, legalAcceptanceRepo, txManager}
}

func (s *companyService) Create(ctx context.Context, req request.CreateCompanyRequest) (err error) {
	if err = req.Validate(); err != nil {
		return err
	}

	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	companyId, err := s.companyRepository.CreateWithTx(ctx, tx, domain.Company{
		Name:      req.Name,
		LegalName: req.LegalName,
		Cnpj:      req.Cnpj,
		Cpf:       req.Cpf,
		Cellphone: req.Cellphone,
	})
	if err != nil {
		if errors.IsDuplicated(err) {
			if pqErr, ok := err.(*pq.Error); ok {
				detail := strings.ToLower(pqErr.Detail)
				constraint := strings.ToLower(pqErr.Constraint)
				if strings.Contains(detail, "cnpj") || strings.Contains(constraint, "cnpj") {
					return errors.New("CNPJ já cadastrado")
				}
				if strings.Contains(detail, "cpf") || strings.Contains(constraint, "cpf") {
					return errors.New("CPF já cadastrado")
				}
				if strings.Contains(detail, "cellphone") || strings.Contains(constraint, "cellphone") {
					return errors.New("Telefone já cadastrado na empresa")
				}
			}
			return errors.New("Empresa já cadastrada")
		}
		return err
	}

	_, err = s.addressRepository.CreateWithTx(ctx, tx, domain.Address{
		Street:       req.Address.Street,
		Neighborhood: req.Address.Neighborhood,
		Number:       req.Address.Number,
		City:         req.Address.City,
		UF:           req.Address.UF,
		Cep:          req.Address.Cep,
		TenantId:     companyId,
	})
	if err != nil {
		return err
	}

	ctxWithTenant := context.WithValue(ctx, constants.TENANT_KEY, companyId)

	adminUser := domain.NewUser(domain.CreateUserParams{
		Username:    req.User.Username,
		Name:        req.User.Name,
		PhoneNumber: req.User.PhoneNumber,
		Role:        string(domain.UserRoleAdmin),
		Email:       req.User.Email,
	})

	encryptedPassword, err := s.encrypto.Encrypt(req.User.Password)
	if err != nil {
		return err
	}
	adminUser.Password = encryptedPassword

	adminId, err := s.userRepository.CreateWithTx(ctxWithTenant, tx, adminUser)
	if err != nil {
		if errors.IsDuplicated(err) {
			return errors.ParseDuplicatedMessage("Usuário", err)
		}
		return err
	}

	adminUser.Id = adminId
	adminUser.TenantId = companyId

	termsDoc, err := s.legalDocumentRepo.GetLastActiveByType(ctx, domain.LegalDocumentTypeTerms)
	if err != nil {
		return err
	}

	privacyDoc, err := s.legalDocumentRepo.GetLastActiveByType(ctx, domain.LegalDocumentTypePrivacy)
	if err != nil {
		return err
	}

	_, err = s.legalAcceptanceRepo.CreateWithTx(ctxWithTenant, tx, domain.LegalAcceptance{
		UserId:          adminId,
		TenantId:        companyId,
		LegalDocumentId: termsDoc.Id,
		Accepted:        true,
	})
	if err != nil {
		return err
	}

	_, err = s.legalAcceptanceRepo.CreateWithTx(ctxWithTenant, tx, domain.LegalAcceptance{
		UserId:          adminId,
		TenantId:        companyId,
		LegalDocumentId: privacyDoc.Id,
		Accepted:        true,
	})
	if err != nil {
		return err
	}

	_, err = s.inventoryRepo.CreateWithTx(ctxWithTenant, tx, domain.Inventory{TenantId: companyId, User: domain.User{Id: adminUser.Id}, Type: domain.InventoryTypePrimary})
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		welcomeErr := s.emailUsecase.SendWelcome(ctx, req.User.Email, req.User.Name)
		if welcomeErr != nil {
			logs.Logger.Errorf("Erro ao enviar email de boas vindas para o usuário %s: %v", req.User.Username, welcomeErr)
		}
	}()

	return nil
}
