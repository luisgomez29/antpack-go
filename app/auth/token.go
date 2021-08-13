package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/luisgomez29/antpack-go/app/models"
	"github.com/luisgomez29/antpack-go/app/utils"
	"github.com/luisgomez29/antpack-go/pkg/config"
)

// AuthorizationTypeBearer Authorization type Bearer
const AuthorizationTypeBearer = "Bearer"

// JWT tokens type
const (
	JWTAccessToken  = "access"
	JWTRefreshToken = "refresh"
)

// Errors
var (
	errJWTMissing    = echo.NewHTTPError(http.StatusBadRequest, "token faltante o tiene un formato incorrecto")
	errJWTInvalid    = echo.NewHTTPError(http.StatusUnauthorized, "token inválido o expirado")
	errJWTimeSetting = echo.NewHTTPError(http.StatusInternalServerError, "Invalid time definition in .env file")
)

// Claims defines the username of the user and the standard claims to generate the JWT token.
type Claims struct {
	jwt.StandardClaims

	TokenType string
	User      *models.User
}

// NewClaims create the claims with values for the Id, IssuedAt and User.
func NewClaims(u *models.User) *Claims {
	return &Claims{
		StandardClaims: jwt.StandardClaims{
			Id:       uuid.NewString(),
			IssuedAt: time.Now().Unix(),
		},
		User: u,
	}
}

// GenerateToken generate a JWT token from the claims.
func GenerateToken(c *Claims) (string, error) {
	claims := jwt.MapClaims{
		"token_type": c.TokenType,
		"email":      c.User.Email,
		"jti":        c.Id,
		"iat":        c.IssuedAt,
		"exp":        c.ExpiresAt,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(
		[]byte(config.Load("JWT_SIGNING_KEY")),
	)

	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return token, nil
}

// VerifyToken verify that the token is valid. If a value is assigned to the `tokenType` parameter,
// the `token_type` of the claim is verified.
func VerifyToken(token string, tokenType ...string) (jwt.MapClaims, error) {
	tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Load("JWT_SIGNING_KEY")), nil
	})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorMalformed:
				return nil, errJWTMissing
			case jwt.ValidationErrorExpired, jwt.ValidationErrorSignatureInvalid:
				return nil, errJWTInvalid
			default:
				return nil, errJWTMissing
			}
		default:
			return nil, errJWTMissing
		}
	}

	claims, ok := tk.Claims.(jwt.MapClaims)
	if !ok || !tk.Valid {
		return nil, errJWTMissing
	}

	// Verify token type
	if tokenType != nil {
		if claims["token_type"] != tokenType[0] {
			return nil, errJWTInvalid
		}
		return claims, nil
	}
	return claims, nil
}

// ExtractToken get the token from the request header.
func ExtractToken(authzHeader string) (string, error) {
	l := len(AuthorizationTypeBearer)
	if len(authzHeader) > l+1 && authzHeader[:l] == AuthorizationTypeBearer {
		return authzHeader[l+1:], nil
	}
	return "", errJWTMissing
}

// newAccessAndRefreshClaims defines the claims of the access JWT token.
func newAccessTokenClaims(u *models.User) (*Claims, error) {
	atTime, err := utils.TimeDuration(config.Load("JWT_ACCESS_TOKEN_EXPIRATION_MINUTES"))
	if err != nil {
		return nil, errJWTimeSetting
	}

	acClaims := NewClaims(u)
	acClaims.ExpiresAt = time.Now().Add(time.Minute * atTime).Unix()
	acClaims.TokenType = JWTAccessToken
	return acClaims, nil
}
