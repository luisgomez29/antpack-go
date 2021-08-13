package faker

import (
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/luisgomez29/antpack-go/app/auth"
	"github.com/luisgomez29/antpack-go/app/models"
)

// User generate a fake user for testing
func User() (user *models.User, password string) {
	password = "lg123456"
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()

	hashedPassword, err := auth.HashPassword(auth.NewPasswordConfig(), password)
	if err != nil {
		log.Fatal(err)
	}

	user = &models.User{
		FirstName: firstName,
		LastName:  lastName,
		FullName:  fmt.Sprintf("%s %s", firstName, lastName),
		Email:     "email@example.com",
		Password:  hashedPassword,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	return
}
