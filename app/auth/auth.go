// Package auth contains the types and functions related to user authentication.
package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/luisgomez29/antpack-go/app/models"
	"github.com/luisgomez29/antpack-go/app/repositories"
)

type Auth interface {
	// HashPassword get the argon2id hash of the password.
	HashPassword(password string) (string, error)

	// VerifyPassword verifies that the password hash matches the password entered by the user.
	VerifyPassword(password, hashedPassword string) (bool, error)

	// GetToken generates the JWT access token.
	GetToken(u *models.User) (map[string]string, error)

	// UserEmailFromContext get user email from request context.
	UserEmailFromContext(c echo.Context) uint

	// IsAuthenticated check if the user is logged in. If the user is logged in, it returns AccessDetails and true.
	IsAuthenticated(c echo.Context) (*AccessDetails, bool)
}

type (
	// AccessDetails represents the user who is logged in.
	AccessDetails struct {
		TokenUuid string
		User      *models.User
	}

	// JWTResponse is the response when the user logs in or register.
	JWTResponse struct {
		AccessToken string       `json:"access_token"`
		User        *models.User `json:"user"`
	}
)

type auth struct {
	authRepo repositories.AuthRepository
}

func NewAuth(at repositories.AuthRepository) Auth {
	return auth{authRepo: at}
}

func (auth) HashPassword(password string) (string, error) {
	return HashPassword(NewPasswordConfig(), password)
}

func (auth) VerifyPassword(password, hashedPassword string) (bool, error) {
	return comparePasswordAndHash(password, hashedPassword)
}

func (auth) GetToken(u *models.User) (map[string]string, error) {
	c, err := newAccessTokenClaims(u)
	if err != nil {
		return nil, err
	}

	accessToken, err := GenerateToken(c)
	if err != nil {
		return nil, err
	}

	tokens := map[string]string{
		"access": accessToken,
	}
	return tokens, nil
}

func (a auth) UserEmailFromContext(c echo.Context) uint {
	user := c.Get("user")
	if user == nil {
		return 0
	}
	cl := user.(jwt.MapClaims)
	return cl["email"].(uint)
}

func (a auth) IsAuthenticated(c echo.Context) (*AccessDetails, bool) {
	id := a.UserEmailFromContext(c)
	if id == 0 {
		return &AccessDetails{}, false
	}

	u := a.authRepo.GetUser(id)
	return &AccessDetails{User: u}, true
}
