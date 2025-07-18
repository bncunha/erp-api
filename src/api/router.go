package router

import (
	"net/http"
	"os"

	controller "github.com/bncunha/erp-api/src/api/controllers"
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

func (r *router) SetupRoutes() {
  productGroup := r.echo.Group("/products")
  productGroup.POST("/", r.controller.ProductController.Create)
  productGroup.GET("/", r.controller.ProductController.GetAll) 
  productGroup.GET("/:id", r.controller.ProductController.GetById)
  productGroup.PUT("/:id", r.controller.ProductController.Edit)
  productGroup.DELETE("/:id", r.controller.ProductController.Inactivate)

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
