package models

import (
	"errors"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// User represents the data about a user.
type User struct {
	ID        uint       `json:"id"`
	FirstName string     `gorm:"not null;size:40" json:"first_name,omitempty"`
	LastName  string     `gorm:"not null;size:40;" json:"last_name,omitempty"`
	FullName  string     `gorm:"-" json:"full_name,omitempty"`
	Email     string     `gorm:"not null;size:60;unique" json:"email,omitempty"`
	Password  string     `gorm:"not null;size:128" json:"password,omitempty"`
	CreatedAt *time.Time `gorm:"not null;default:now()" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"not null;default:now()" json:"updated_at,omitempty"`
}

func (u *User) AfterSave(*gorm.DB) error {
	u.FullName = u.FirstName + " " + u.LastName
	return nil
}

func (u *User) AfterFind(*gorm.DB) error {
	u.FullName = u.FirstName + " " + u.LastName
	return nil
}

func (*User) ValidatePgError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			switch pgErr.ConstraintName {
			case "users_email_key":
				e := errors.New("ya existe un usuario con este correo electrónico")
				return validation.Errors{"username": e}
			}
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err)
}
