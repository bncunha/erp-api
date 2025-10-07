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
			AllowOrigins: []string{"https://erp-front-production.onrender.com", "https://trinus.app"},
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
	r.echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}

func (r *router) setupPrivateRoutes() {
	private := r.echo.Group("")

	private.Use(middleware.AuthMiddleware)

	productGroup := private.Group("/products")
	productGroup.POST("", r.controller.ProductController.Create, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	productGroup.GET("", r.controller.ProductController.GetAll, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin, domain.UserRoleReseller}))
	productGroup.GET("/:id", r.controller.ProductController.GetById, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	productGroup.PUT("/:id", r.controller.ProductController.Edit, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	productGroup.DELETE("/:id", r.controller.ProductController.Inactivate, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	productGroup.GET("/:id/skus", r.controller.ProductController.GetSkus, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	productGroup.POST("/:id/skus", r.controller.SkuController.Create, middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))

	skuGroup := private.Group("/skus", middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	skuGroup.GET("", r.controller.SkuController.GetAll)
	skuGroup.PUT("/:id", r.controller.SkuController.Edit)
	skuGroup.GET("/:id", r.controller.SkuController.GetById)
	skuGroup.DELETE("/:id", r.controller.SkuController.Inactivate)

	categoryGroup := private.Group("/categories", middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	categoryGroup.POST("", r.controller.CategoryController.Create)
	categoryGroup.GET("", r.controller.CategoryController.GetAll)
	categoryGroup.GET("/:id", r.controller.CategoryController.GetById)
	categoryGroup.PUT("/:id", r.controller.CategoryController.Edit)
	categoryGroup.DELETE("/:id", r.controller.CategoryController.Inactivate)

	userGroup := private.Group("/users", middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	userGroup.POST("", r.controller.UserController.Create)
	userGroup.GET("", r.controller.UserController.GetAll)
	userGroup.GET("/:id", r.controller.UserController.GetById)
	userGroup.PUT("/:id", r.controller.UserController.Edit)
	userGroup.DELETE("/:id", r.controller.UserController.Inactivate)

	inventoryGroup := private.Group("/inventory", middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	inventoryGroup.GET("", r.controller.InventoryController.GetAllInventories)
	inventoryGroup.GET("/:id/items", r.controller.InventoryController.GetInventoryItemsByInventoryId)
	inventoryGroup.GET("/items", r.controller.InventoryController.GetAllInventoryItems)
	inventoryGroup.GET("/transaction", r.controller.InventoryController.GetAllInventoryTransactions)
	inventoryGroup.POST("/transaction", r.controller.InventoryController.DoTransaction)

	salesGroup := private.Group("/sales", middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	salesGroup.POST("", r.controller.SalesController.Create)
	salesGroup.GET("", r.controller.SalesController.GetAll)
	salesGroup.GET("/:id", r.controller.SalesController.GetById)
	salesGroup.PUT("/:id/payments/:payment_id", r.controller.SalesController.ChangePaymentStatus)

	customerGroup := private.Group("/customers", middleware.RoleMiddleware([]domain.Role{domain.UserRoleAdmin}))
	customerGroup.GET("", r.controller.CustomerController.GetAll)
}

func (r *router) Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.echo.Logger.Fatal(r.echo.Start(":" + port))
}
