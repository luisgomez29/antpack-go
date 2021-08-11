package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/luisgomez29/antpack-go/app/middlewares"
	"github.com/luisgomez29/antpack-go/app/resources/migrations"
	"github.com/luisgomez29/antpack-go/app/routes"
	"github.com/luisgomez29/antpack-go/pkg/config"
	"github.com/luisgomez29/antpack-go/pkg/database"
)

func main() {
	db := database.Connect()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Migrate database

	migrations.Load(db)

	e := echo.New()

	// Middlewares
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.Secure(),
		middlewares.ErrorHandler,
	)

	// Routes
	routes.Setup(db, e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Load("SERVER_PORT"))))
}
