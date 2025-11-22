package controller

import (
	_http "net/http"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/api/viewmodel"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService,
	}
}

func (c *UserController) Create(context echo.Context) error {
	var userRequest request.CreateUserRequest
	if err := context.Bind(&userRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	err := c.userService.Create(context.Request().Context(), userRequest)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusCreated, nil)
}

func (c *UserController) Edit(context echo.Context) error {
	userId := helper.ParseInt64(context.Param("id"))
	var userRequest request.EditUserRequest

	if err := context.Bind(&userRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	err := c.userService.Update(context.Request().Context(), userRequest, userId)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, nil)
}

func (c *UserController) GetById(context echo.Context) error {
	userId := helper.ParseInt64(context.Param("id"))

	user, err := c.userService.GetById(context.Request().Context(), userId)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, viewmodel.ToUserViewModel(user))
}

func (c *UserController) GetAll(context echo.Context) error {
	users, err := c.userService.GetAll(context.Request().Context(), request.GetAllUserRequest{
		Role: domain.Role(context.QueryParam("role")),
	})
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	var userViewModels []viewmodel.UserViewModel = make([]viewmodel.UserViewModel, 0)
	for _, user := range users {
		userViewModels = append(userViewModels, viewmodel.ToUserViewModel(user))
	}

	return context.JSON(_http.StatusOK, userViewModels)
}

func (c *UserController) Inactivate(context echo.Context) error {
	id := helper.ParseInt64(context.Param("id"))

	err := c.userService.Inactivate(context.Request().Context(), id)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, nil)
}

func (c *UserController) ChangePassword(context echo.Context) error {
	var resetPasswordRequest request.ResetPasswordRequest
	if err := context.Bind(&resetPasswordRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	err := c.userService.ResetPassword(context.Request().Context(), resetPasswordRequest)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, nil)
}

func (c *UserController) ForgotPassword(context echo.Context) error {
	var forgotPasswordRequest request.ForgotPasswordRequest
	if err := context.Bind(&forgotPasswordRequest); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	err := c.userService.ForgotPassword(context.Request().Context(), forgotPasswordRequest)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}

	return context.JSON(_http.StatusOK, nil)
}
