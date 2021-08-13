package services

import (
	"github.com/luisgomez29/antpack-go/app/auth"
	"github.com/luisgomez29/antpack-go/app/models"
	"github.com/luisgomez29/antpack-go/app/repositories"
	apiErrors "github.com/luisgomez29/antpack-go/app/resources/api/errors"
	"github.com/luisgomez29/antpack-go/app/resources/api/requests"
)

// AccountsService encapsulates usecase logic for account management.
type AccountsService interface {
	SignUp(input *requests.SignUpRequest) (*auth.JWTResponse, error)
	Login(input *requests.LoginRequest) (*auth.JWTResponse, error)
}

type accountsService struct {
	auth         auth.Auth
	accountsRepo repositories.AccountRepository
}

// NewAccountsService create a new accounts service.
func NewAccountsService(at auth.Auth, u repositories.AccountRepository) AccountsService {
	return accountsService{auth: at, accountsRepo: u}
}

func (s accountsService) SignUp(input *requests.SignUpRequest) (*auth.JWTResponse, error) {
	// Generating password Hash
	hash, err := s.auth.HashPassword(input.Password)
	if err != nil {
		return &auth.JWTResponse{}, err
	}

	input.Password = hash
	user, err := s.accountsRepo.CreateUser(input)
	if err != nil {
		return &auth.JWTResponse{}, err
	}

	return s.tokenAndUser(user)
}

func (s accountsService) Login(input *requests.LoginRequest) (*auth.JWTResponse, error) {
	user, err := s.accountsRepo.FindUser(input.Email)
	if err != nil {
		return &auth.JWTResponse{}, err
	}

	match, err := s.auth.VerifyPassword(input.Password, user.Password)
	if !match || err != nil {
		return &auth.JWTResponse{}, apiErrors.Unauthorized("la contraseña ingresada es incorrecta")
	}

	return s.tokenAndUser(user)
}

// tokenAndUser returns the access JWT token and the user.
func (s accountsService) tokenAndUser(user *models.User) (*auth.JWTResponse, error) {
	user.Password = ""
	tokens, err := s.auth.GetToken(user)
	if err != nil {
		return &auth.JWTResponse{}, err
	}

	return &auth.JWTResponse{
		AccessToken: tokens["access"],
		User:        user,
	}, nil
}
