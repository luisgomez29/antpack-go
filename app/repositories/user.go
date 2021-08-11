package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/luisgomez29/antpack-go/app/models"
	apiErrors "github.com/luisgomez29/antpack-go/app/resources/api/errors"
)

// UserRepository encapsulates the logic to access users from the data source.
type UserRepository interface {
	All() ([]*models.User, error)
	Get(id uint) (*models.User, error)
}

// userRepository persists users in database.
type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) All() ([]*models.User, error) {
	var users []*models.User
	r.conn.Omit("password").Find(&users).Limit(50)
	return users, nil
}

func (r *userRepository) Get(id uint) (*models.User, error) {
	u := new(models.User)
	if err := r.conn.Omit("password").Take(u, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apiErrors.NewErrNoRows("usuario no encontrado")
	}
	return u, nil
}
