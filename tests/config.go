package tests

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"

	"github.com/luisgomez29/antpack-go/app/auth"
	"github.com/luisgomez29/antpack-go/app/middlewares"
	"github.com/luisgomez29/antpack-go/app/models"
	"github.com/luisgomez29/antpack-go/pkg/database"
)

var DB *gorm.DB
var E *echo.Echo

// Init Initializes the connection to the database and the server.
func Init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	DB = database.Connect()
	sqlDB, _ := DB.DB()
	defer sqlDB.Close()

	E = echo.New()

	// Middlewares
	E.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.Secure(),
		middlewares.ErrorHandler,
	)
}

// Route define the route group.
func Route(e *echo.Echo) *echo.Group {
	// Routes
	api := e.Group("/api")
	v1 := api.Group("/v1")
	return v1
}

// AddAuthorization add authorization header and JWT token to request.
func AddAuthorization(request *http.Request, user *models.User, duration time.Duration) {
	claims := auth.NewClaims(user)
	claims.ExpiresAt = time.Now().Add(time.Minute * duration).Unix()
	claims.TokenType = auth.JWTAccessToken

	token, err := auth.GenerateToken(claims)
	if err != nil {
		log.Fatal(err)
	}

	authorizationHeader := fmt.Sprintf("%s %s", auth.AuthorizationTypeBearer, token)
	request.Header.Set(echo.HeaderAuthorization, authorizationHeader)
}
