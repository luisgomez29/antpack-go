package middlewares

import (
	"github.com/labstack/echo/v4"

	"github.com/luisgomez29/antpack-go/app/auth"
)

// Authentication is the JWT Token authentication middleware.
// If the token is valid, the auth.AccessDetails are stored in the context of the request under the user key.
func Authentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				return err
			}

			authzHeader := c.Request().Header.Get(echo.HeaderAuthorization)

			tokenString, err := auth.ExtractToken(authzHeader)
			if err != nil {
				return err
			}

			claims, err := auth.VerifyToken(tokenString, auth.JWTAccessToken)
			if err != nil {
				return err
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}
