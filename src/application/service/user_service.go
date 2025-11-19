package service

import (
	"context"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/service/input"
	emailusecase "github.com/bncunha/erp-api/src/application/usecase/email_usecase"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/logs"
)

type UserService interface {
	Create(ctx context.Context, request request.CreateUserRequest) error
	Update(ctx context.Context, request request.EditUserRequest, userId int64) error
	GetById(ctx context.Context, userId int64) (domain.User, error)
	GetAll(ctx context.Context, request request.GetAllUserRequest) ([]domain.User, error)
	Inactivate(ctx context.Context, id int64) error
	ResetPassword(ctx context.Context, request request.ResetPasswordRequest) error
	ForgotPassword(ctx context.Context, request request.ForgotPasswordRequest) error
}

type userService struct {
	userRepository      domain.UserRepository
	inventoryRepository domain.InventoryRepository
	encrypto            domain.Encrypto
	userTokenService    UserTokenService
	emailUsecase        emailusecase.EmailUseCase
	userTokenRepository domain.UserTokenRepository
}

func NewUserService(userRepository domain.UserRepository, inventoryRepository domain.InventoryRepository, encrypto domain.Encrypto, userTokenService UserTokenService, emailUsecase emailusecase.EmailUseCase, userTokenRepository domain.UserTokenRepository) UserService {
	return &userService{userRepository, inventoryRepository, encrypto, userTokenService, emailUsecase, userTokenRepository}
}

func (s *userService) Create(ctx context.Context, request request.CreateUserRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	user := domain.NewUser(domain.CreateUserParams{
		Username:    request.Username,
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		Role:        request.Role,
		Email:       request.Email,
	})

	userId, err := s.userRepository.Create(ctx, user)
	if err != nil {
		if errors.IsDuplicated(err) {
			return errors.New("Usuário já cadastrado!")
		}
		return err
	}

	user, err = s.userRepository.GetById(ctx, userId)
	if err != nil {
		return err
	}

	adminUser, err := s.userRepository.GetById(ctx, int64(ctx.Value(constants.USERID_KEY).(float64)))
	if err != nil {
		return err
	}

	userToken, err := s.userTokenService.Create(ctx, input.CreateUserTokenInput{
		User:      user,
		CreatedBy: adminUser,
		Type:      domain.UserTokenTypeInvite,
	})
	if err != nil {
		return err
	}
	code := userToken.Code

	if request.Role == string(domain.InventoryTypeReseller) {
		inventory := domain.Inventory{
			User: domain.User{Id: userId},
			Type: domain.InventoryTypeReseller,
		}
		_, err = s.inventoryRepository.Create(ctx, inventory)
		if err != nil {
			return err
		}
	}

	go func() {
		err := s.emailUsecase.SendInvite(ctx, user, code, userToken.Uuid)
		if err != nil {
			logs.Logger.Errorf("Erro ao enviar email de convite para o usuário %d: %v", userId, err)
		}
	}()

	return nil
}

func (s *userService) Update(ctx context.Context, request request.EditUserRequest, userId int64) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	user := domain.User{
		Id:          userId,
		Username:    request.Username,
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		Role:        request.Role,
	}

	err = s.userRepository.Update(ctx, user)
	if err != nil {
		if errors.IsDuplicated(err) {
			return errors.New("Usuário já cadastrado!")
		}
		return err
	}

	if request.Role == string(domain.InventoryTypeReseller) {
		_, err := s.inventoryRepository.GetByUserId(ctx, userId)
		if err != nil && !errors.Is(err, domain.ErrInventoryNotFound) {
			return err
		}

		if err != nil && errors.Is(err, domain.ErrInventoryNotFound) {
			inventory := domain.Inventory{
				User: domain.User{Id: userId},
				Type: domain.InventoryTypeReseller,
			}
			_, err = s.inventoryRepository.Create(ctx, inventory)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *userService) GetById(ctx context.Context, userId int64) (domain.User, error) {
	user, err := s.userRepository.GetById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *userService) GetAll(ctx context.Context, request request.GetAllUserRequest) ([]domain.User, error) {
	var role *domain.Role

	if request.Role != "" {
		role = &request.Role
	}
	users, err := s.userRepository.GetAll(ctx, input.GetAllUserInput{Role: role})
	if err != nil {
		return users, err
	}
	return users, nil
}

func (s *userService) Inactivate(ctx context.Context, id int64) error {
	return s.userRepository.Inactivate(ctx, id)
}

func (s *userService) ResetPassword(ctx context.Context, request request.ResetPasswordRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	userToken, err := s.userTokenRepository.GetLastActiveByUuid(ctx, request.Uuid)
	if err != nil {
		logs.Logger.Errorf("Erro ao recuperar token de reset: %v", err)
		return errors.New("Código expirado ou inválido! Solicite um novo código para o administrador.")
	}
	if isValid, err := userToken.IsValid(s.encrypto, request.Code); err != nil || !isValid {
		logs.Logger.Errorf("Erro ao recuperar token de reset: %v", err)
		return errors.New("Código expirado ou inválido! Solicite um novo código para o administrador.")
	}

	newPassword, err := s.encrypto.Encrypt(request.Password)
	if err != nil {
		return err
	}

	err = s.userRepository.UpdatePassword(ctx, userToken.User, newPassword)
	if err != nil {
		return err
	}

	userToken.SetUsedAt()
	err = s.userTokenRepository.SetUsedToken(ctx, userToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) ForgotPassword(ctx context.Context, request request.ForgotPasswordRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	user, err := s.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		logs.Logger.Errorf("Erro ao recuperar usuário por email: %v", err)
		return nil
	}

	userToken := domain.NewUserToken(domain.CreateUserTokenParams{
		User:      user,
		CreatedBy: user,
		Type:      domain.UserTokenTypeResetPass,
	}, s.encrypto)
	code := userToken.Code
	userTokenId, err := s.userTokenRepository.Create(ctx, userToken)
	if err != nil {
		logs.Logger.Errorf("Erro ao criar token de recuperação de senha: %v", err)
		return nil
	}

	userToken, err = s.userTokenRepository.GetById(ctx, userTokenId)
	if err != nil {
		logs.Logger.Errorf("Erro ao recuperar token de recuperação de senha: %v", err)
		return nil
	}

	go func() {
		err := s.emailUsecase.SendRecoverPassword(ctx, user, code, userToken.Uuid)
		if err != nil {
			logs.Logger.Errorf("Erro ao enviar email de recuperação de senha para o usuário %d: %v", user.Id, err)
		}
	}()

	return nil
}
