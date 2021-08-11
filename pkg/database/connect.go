package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/luisgomez29/antpack-go/pkg/config"
)

// Connect open a connection to the database.
func Connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Bogota",
		config.Load("DB_HOST"), config.Load("DB_USER"), config.Load("DB_PWD"), config.Load("DB_NAME"),
		config.Load("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		QueryFields: true,
	})

	if err != nil {
		log.Fatal(err)
	}
	return db
}
