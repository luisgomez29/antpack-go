package repositories

import (
	"gorm.io/gorm"

	"github.com/luisgomez29/antpack-go/app/models"
)

type AuthRepository interface {
	// GetUser get user from UserRepository.Get.
	GetUser(id uint) *models.User
}

type authRepository struct {
	conn *gorm.DB
	user UserRepository
}

// NewAuthRepository creates a new auth repository.
func NewAuthRepository(db *gorm.DB, u UserRepository) AuthRepository {
	return authRepository{conn: db, user: u}
}

func (r authRepository) GetUser(id uint) *models.User {
	user, _ := r.user.Get(id)
	return user
}
