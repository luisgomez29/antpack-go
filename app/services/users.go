package services

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/luisgomez29/antpack-go/app/auth"
	"github.com/luisgomez29/antpack-go/app/models"
	"github.com/luisgomez29/antpack-go/app/repositories"
	apiErrors "github.com/luisgomez29/antpack-go/app/resources/api/errors"
)

// UsersService encapsulates usecase logic for users.
type UsersService interface {
	All(c echo.Context) (map[string][]*models.User, error)
	Get(c echo.Context) (*models.User, error)
}

type usersService struct {
	auth      auth.Auth
	usersRepo repositories.UserRepository
}

// NewUsersService create a new users service.
func NewUsersService(at auth.Auth, u repositories.UserRepository) UsersService {
	return usersService{auth: at, usersRepo: u}
}

func (s usersService) All(c echo.Context) (map[string][]*models.User, error) {
	users, err := s.usersRepo.All()
	if err != nil {
		return nil, err
	}

	return map[string][]*models.User{"results": users}, nil
}

func (s usersService) Get(c echo.Context) (*models.User, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, apiErrors.BadRequest("")
	}

	user, err := s.usersRepo.Get(uint(id))
	if err != nil {
		return nil, err
	}

	return user, nil
}
