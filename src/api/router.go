package router

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type router struct{
  echo *echo.Echo
}

func NewRouter() *router {
  e := echo.New()
	router := &router{
    echo: e,
  }
	return router
}

func (r *router) SetupRoutes() {
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
