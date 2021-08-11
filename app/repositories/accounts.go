package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/luisgomez29/antpack-go/app/models"
	apiErrors "github.com/luisgomez29/antpack-go/app/resources/api/errors"
	"github.com/luisgomez29/antpack-go/app/resources/api/requests"
)

// AccountRepository encapsulates the logic to access users from the data source.
type AccountRepository interface {

	// FindUser get the user by email.
	FindUser(email string) (*models.User, error)

	// CreateUser create a user in the database.
	CreateUser(input *requests.SignUpRequest) (*models.User, error)
}

type accountRepository struct {
	conn *gorm.DB
}

// NewAccountRepository creates a new account repository.
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return accountRepository{conn: db}
}

func (r accountRepository) FindUser(email string) (*models.User, error) {
	user := new(models.User)
	if err := r.conn.Debug().Where("email = ?", email).Take(user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apiErrors.NewErrNoRows("usuario o contraseña incorrectos")
	}
	return user, nil
}

func (r accountRepository) CreateUser(input *requests.SignUpRequest) (*models.User, error) {
	user := new(models.User)
	if err := copier.Copy(user, input); err != nil {
		return nil, err
	}

	err := r.conn.Create(user).Error
	if err != nil {
		return nil, user.ValidatePgError(err)
	}

	return user, nil
}
