package router

import (
	"net/http"
	"os"

	_middleware "github.com/labstack/echo/v4/middleware"

	controller "github.com/bncunha/erp-api/src/api/controllers"
	"github.com/bncunha/erp-api/src/api/middleware"
	"github.com/labstack/echo/v4"
)

type router struct{
  controller *controller.Controller
  echo *echo.Echo
}

func NewRouter(controller *controller.Controller) *router {
  e := echo.New()
	router := &router{
    echo: e,
    controller: controller,
  }
	return router
}

func (r *router) SetupCors(env string) {
  r.echo.Use(_middleware.CORSWithConfig(_middleware.CORSConfig{
    AllowOrigins: []string{"http://localhost:4200", "https://erp-front-0pem.onrender.com"},
    AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
  }))
}

func (r *router) SetupRoutes() {
  r.echo.Use(middleware.TenantMiddleware)
  productGroup := r.echo.Group("/products")
  productGroup.POST("", r.controller.ProductController.Create)
  productGroup.GET("", r.controller.ProductController.GetAll) 
  productGroup.GET("/:id", r.controller.ProductController.GetById)
  productGroup.PUT("/:id", r.controller.ProductController.Edit)
  productGroup.DELETE("/:id", r.controller.ProductController.Inactivate)
  productGroup.GET("/:id/skus", r.controller.ProductController.GetSkus)
  productGroup.POST("/:id/skus", r.controller.SkuController.Create)

  skuGroup := r.echo.Group("/skus")
  skuGroup.GET("", r.controller.SkuController.GetAll)
  skuGroup.PUT("/:id", r.controller.SkuController.Edit)
  skuGroup.GET("/:id", r.controller.SkuController.GetById)
  skuGroup.DELETE("/:id", r.controller.SkuController.Inactivate)

  categoryGroup := r.echo.Group("/categories")
  categoryGroup.POST("", r.controller.CategoryController.Create)
  categoryGroup.GET("", r.controller.CategoryController.GetAll) 
  categoryGroup.GET("/:id", r.controller.CategoryController.GetById)
  categoryGroup.PUT("/:id", r.controller.CategoryController.Edit)
  categoryGroup.DELETE("/:id", r.controller.CategoryController.Inactivate)
  
	r.echo.GET("/health", func (c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}

func (r *router) Start() {
  port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
  r.echo.Logger.Fatal(r.echo.Start(":" + port))
}
