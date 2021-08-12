package controllers

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"

	"github.com/luisgomez29/antpack-go/app/middlewares"
	"github.com/luisgomez29/antpack-go/app/routes"
	"github.com/luisgomez29/antpack-go/pkg/database"
)

var db *gorm.DB
var e *echo.Echo

func TestMain(t *testing.M) {
	if err := godotenv.Load("../../.env.test"); err != nil {
		log.Fatal("Error loading .env file")
	}

	db = database.Connect()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	e = echo.New()

	// Middlewares
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.Secure(),
		middlewares.ErrorHandler,
	)

	// Routes
	routes.Setup(db, e)
	code := t.Run()

	db.Exec("DELETE FROM users")
	os.Exit(code)
}
