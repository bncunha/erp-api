package service

import (
	"github.com/bncunha/erp-api/src/application/ports"
	"github.com/bncunha/erp-api/src/application/usecase"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type ApplicationService struct {
	ProductService   ProductService
	SkuService       SkuService
	CategoryService  CategoryService
	AuthService      AuthService
	UserService      UserService
	InventoryService InventoryService
	SalesService     SalesService
	CustomerService  CustomerService
	CompanyService   CompanyService
	UserTokenService UserTokenService
	DashboardService DashboardService
	BillingService   BillingService
	repositories     *repository.Repository
	useCases         *usecase.ApplicationUseCase
	ports            *ports.Ports
}

func NewApplicationService(repositories *repository.Repository, useCases *usecase.ApplicationUseCase, ports *ports.Ports) *ApplicationService {
	return &ApplicationService{repositories: repositories, useCases: useCases, ports: ports}
}

func (s *ApplicationService) SetupServices() {
	s.UserTokenService = NewUserTokenService(s.repositories.UserTokenRepository, s.ports.Encrypto)
	s.ProductService = NewProductService(s.repositories.ProductRepository, s.repositories.CategoryRepository, s.repositories.SkuRepository)
	s.SkuService = NewSkuService(
		s.repositories.SkuRepository,
		s.useCases.InventoryUseCase,
		s.repositories.ProductRepository,
		s.repositories.InventoryItemRepository,
		s.repositories.InventoryTransactionRepository,
		s.repositories,
	)
	s.CategoryService = NewCategoryService(s.repositories.CategoryRepository)
	s.BillingService = NewBillingService(s.repositories.PlanRepository, s.repositories.SubscriptionRepository, s.repositories.BillingPaymentRepository, s.repositories)
	s.AuthService = NewAuthService(s.repositories.UserRepository, s.ports.Encrypto, s.BillingService)
	s.UserService = NewUserService(s.repositories.UserRepository, s.repositories.InventoryRepository, s.ports.Encrypto, s.UserTokenService, s.useCases.EmailUseCase, s.repositories.UserTokenRepository, s.repositories.LegalDocumentRepository, s.repositories.LegalAcceptanceRepository, s.repositories)
	s.InventoryService = NewInventoryService(s.useCases.InventoryUseCase, s.repositories.InventoryItemRepository, s.repositories.InventoryTransactionRepository, s.repositories.InventoryRepository, s.repositories)
	s.SalesService = NewSalesService(s.useCases.SalesUsecase, s.repositories.SalesRepository, s.repositories.InventoryRepository)
	s.CustomerService = NewCustomerService(s.repositories.CustomerRepository)
	s.CompanyService = NewCompanyService(s.repositories.CompanyRepository, s.repositories.AddressRepository, s.repositories.InventoryRepository, s.repositories.UserRepository, s.ports.Encrypto, s.useCases.EmailUseCase, s.repositories.LegalDocumentRepository, s.repositories.LegalAcceptanceRepository, s.repositories)
	s.DashboardService = NewDashboardService(s.repositories.DashboardRepository, s.repositories.UserRepository)
}
