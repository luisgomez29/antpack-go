package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/luisgomez29/antpack-go/app/auth"
	"github.com/luisgomez29/antpack-go/app/controllers"
	"github.com/luisgomez29/antpack-go/app/middlewares"
	"github.com/luisgomez29/antpack-go/app/repositories"
	"github.com/luisgomez29/antpack-go/app/services"
)

// Accounts defines the endpoints for authentication and account management.
func Accounts(g *echo.Group, ctrl controllers.AccountsController) {
	g.POST("/signup", ctrl.SignUp)
	g.POST("/login", ctrl.Login)
}

// Users defines the endpoints for users management.
func Users(g *echo.Group, ctrl controllers.UsersController) {
	g.Use(middlewares.Authentication())
	g.GET("/users", ctrl.All)
	g.GET("/users/:id", ctrl.Get)
}

// Setup sets the available API endpoints.
func Setup(db *gorm.DB, e *echo.Echo) {
	// API V1
	api := e.Group("/api")
	v1 := api.Group("/v1")

	// Repositories
	usersRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db, usersRepo)
	accountsRepo := repositories.NewAccountRepository(db)

	// Authentication service
	authn := auth.NewAuth(authRepo)

	// Services
	usersService := services.NewUsersService(authn, usersRepo)
	accountsService := services.NewAccountsService(authn, accountsRepo)

	// Routes
	Accounts(v1, controllers.NewAccountsController(authn, accountsService))
	Users(v1, controllers.NewUsersController(usersService))
}
