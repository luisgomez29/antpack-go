package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/luisgomez29/antpack-go/app/auth"
	apiErrors "github.com/luisgomez29/antpack-go/app/resources/api/errors"
	"github.com/luisgomez29/antpack-go/app/resources/api/requests"
	"github.com/luisgomez29/antpack-go/app/services"
)

// AccountsController represents endpoints for authentication.
type AccountsController interface {
	SignUp(c echo.Context) error
	Login(c echo.Context) error
}

type accountsController struct {
	auth            auth.Auth
	accountsService services.AccountsService
}

// NewAccountsController create a new accounts' controller.
func NewAccountsController(at auth.Auth, as services.AccountsService) AccountsController {
	return accountsController{auth: at, accountsService: as}
}

func (ct accountsController) SignUp(c echo.Context) error {
	input := new(requests.SignUpRequest)
	if err := c.Bind(input); err != nil {
		return apiErrors.BadRequest("")
	}

	if err := input.Validate(); err != nil {
		return err
	}

	if input.Password != input.PasswordConfirm {
		return apiErrors.PasswordMismatch
	}

	res, err := ct.accountsService.SignUp(input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

func (ct accountsController) Login(c echo.Context) error {
	input := new(requests.LoginRequest)
	if err := c.Bind(input); err != nil {
		return apiErrors.BadRequest("")
	}

	if err := input.Validate(); err != nil {
		return err
	}

	res, err := ct.accountsService.Login(input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}
