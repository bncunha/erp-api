package router

import (
	"net/http"
	"os"

	_middleware "github.com/labstack/echo/v4/middleware"

	controller "github.com/bncunha/erp-api/src/api/controllers"
	"github.com/bncunha/erp-api/src/api/middleware"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/labstack/echo/v4"
)

type router struct {
	controller *controller.Controller
	echo       *echo.Echo
}

func NewRouter(controller *controller.Controller) *router {
	e := echo.New()
	router := &router{
		echo:       e,
		controller: controller,
	}
	return router
}

func (r *router) SetupCors(env string) {
	if env != "production" {
		r.echo.Use(_middleware.CORSWithConfig(_middleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:4200", "https://erp-front-0pem.onrender.com"},
			AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		}))
	} else {
		r.echo.Use(_middleware.CORSWithConfig(_middleware.CORSConfig{
			AllowOrigins: []string{"https://erp-front-production.onrender.com", "https://trinus.app", "https://erp.trinus.app"},
			AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		}))
	}
}

func (r *router) SetupRoutes() {
	r.echo.Use(middleware.RequestLogger())
	r.echo.Use(middleware.Recover())

	r.setupPublicRoutes()
	r.setupPrivateRoutes()
}

func (r *router) setupPublicRoutes() {
	r.echo.POST("/login", r.controller.AuthController.Login)
	r.echo.POST("/forgot-password", r.controller.UserController.ForgotPassword)
	r.echo.POST("/change-password", r.controller.UserController.ChangePassword)
	r.echo.POST("/signup", r.controller.CompanyController.Create)
	r.echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}

func (r *router) setupPrivateRoutes() {
	private := r.echo.Group("")

	private.Use(middleware.AuthMiddleware)
	private.Use(middleware.BillingWriteGuard())

	productGroup := private.Group("/products")
	productGroup.POST("", r.controller.ProductController.Create, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	productGroup.GET("", r.controller.ProductController.GetAll, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	productGroup.GET("/:id", r.controller.ProductController.GetById, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	productGroup.PUT("/:id", r.controller.ProductController.Edit, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	productGroup.DELETE("/:id", r.controller.ProductController.Inactivate, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	productGroup.GET("/:id/skus", r.controller.ProductController.GetSkus, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	productGroup.POST("/:id/skus", r.controller.SkuController.Create, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))

	skuGroup := private.Group("/skus")
	skuGroup.GET("", r.controller.SkuController.GetAll, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	skuGroup.PUT("/:id", r.controller.SkuController.Edit, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	skuGroup.GET("/:id", r.controller.SkuController.GetById, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	skuGroup.DELETE("/:id", r.controller.SkuController.Inactivate, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	skuGroup.GET("/:id/inventory", r.controller.SkuController.GetInventory, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	skuGroup.GET("/:id/transactions", r.controller.SkuController.GetTransactions, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))

	categoryGroup := private.Group("/categories", middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	categoryGroup.POST("", r.controller.CategoryController.Create)
	categoryGroup.GET("", r.controller.CategoryController.GetAll)
	categoryGroup.GET("/:id", r.controller.CategoryController.GetById)
	categoryGroup.PUT("/:id", r.controller.CategoryController.Edit)
	categoryGroup.DELETE("/:id", r.controller.CategoryController.Inactivate)

	userGroup := private.Group("/users")
	userGroup.POST("", r.controller.UserController.Create, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	userGroup.GET("", r.controller.UserController.GetAll, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	userGroup.GET("/:id", r.controller.UserController.GetById, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	userGroup.GET("/legal-terms", r.controller.UserController.GetLegalTerms, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	userGroup.POST("/legal-terms", r.controller.UserController.AcceptLegalTerms, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	userGroup.PUT("/:id", r.controller.UserController.Edit, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	userGroup.DELETE("/:id", r.controller.UserController.Inactivate, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))

	inventoryGroup := private.Group("/inventory", middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	inventoryGroup.GET("", r.controller.InventoryController.GetAllInventories)
	inventoryGroup.GET("/summary", r.controller.InventoryController.GetInventoriesSummary)
	inventoryGroup.GET("/:id/summary", r.controller.InventoryController.GetInventorySummary)
	inventoryGroup.GET("/:id/items", r.controller.InventoryController.GetInventoryItemsByInventoryId)
	inventoryGroup.GET("/items", r.controller.InventoryController.GetAllInventoryItems)
	inventoryGroup.GET("/:id/transaction", r.controller.InventoryController.GetInventoryTransactionsByInventoryId)
	inventoryGroup.POST("/transaction", r.controller.InventoryController.DoTransaction)

	salesGroup := private.Group("/sales")
	salesGroup.POST("", r.controller.SalesController.Create, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	salesGroup.POST("/:id/returns", r.controller.SalesController.CreateReturn, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	salesGroup.GET("", r.controller.SalesController.GetAll, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	salesGroup.GET("/:id", r.controller.SalesController.GetById, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	salesGroup.PUT("/:id/payments/:payment_id", r.controller.SalesController.ChangePaymentStatus, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))

	customerGroup := private.Group("/customers")
	customerGroup.POST("", r.controller.CustomerController.Create, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	customerGroup.GET("", r.controller.CustomerController.GetAll, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	customerGroup.GET("/:id", r.controller.CustomerController.GetById, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	customerGroup.PUT("/:id", r.controller.CustomerController.Edit, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	customerGroup.DELETE("/:id", r.controller.CustomerController.Inactivate, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))

	dashboardGroup := private.Group("/dashboard")
	dashboardGroup.GET("/widgets", r.controller.DashboardController.GetWidgets)
	dashboardGroup.POST("/widgets/data", r.controller.DashboardController.GetWidgetData)

	billingGroup := private.Group("/billing")
	billingGroup.GET("", r.controller.BillingController.Summary, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	billingGroup.GET("/status", r.controller.BillingController.Status)
	billingGroup.GET("/payments", r.controller.BillingController.Payments, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))

	newsGroup := private.Group("/news")
	newsGroup.GET("/latest", r.controller.NewsController.GetLatest)
}

func (r *router) Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.echo.Logger.Fatal(r.echo.Start(":" + port))
}
