package migrations

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/luisgomez29/antpack-go/app/models"
	"github.com/luisgomez29/antpack-go/app/utils"
)

func Load(db *gorm.DB) {
	mdl := []interface{}{&models.User{}}
	err := db.Migrator().DropTable(mdl...)
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(mdl...)
	if err != nil {
		log.Fatal(err)
	}

	// Insert test data
	db.Create(&users)

	// Show inserted data in console
	for _, user := range users {
		utils.Pretty(user)
	}
	fmt.Println("Database Migrated")
}
