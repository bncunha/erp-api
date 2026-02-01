package controller

import (
	"errors"
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/api/viewmodel"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/bncunha/erp-api/src/infrastructure/logs"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService}
}

func (c *AuthController) Login(context echo.Context) error {
	var loginRequest request.LoginRequest
	if err := context.Bind(&loginRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	output, err := c.authService.Login(context.Request().Context(), loginRequest)
	if err != nil {
		logs.Logger.Errorf("Erro ao fazer login do usuário %s: %v", loginRequest.Username, err)
		return context.JSON(_http.StatusBadRequest, http.HandleError(errors.New("Usuário ou senha incorretos")))
	}
	return context.JSON(_http.StatusOK, viewmodel.ToLoginViewModel(output))
}
