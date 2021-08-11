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
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	FullName  string    `gorm:"-" json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
