// Package errors contains the types and functions related to errors
package errors

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ErrNoRows customize the error message when the error of type gorm.ErrRecordNotFound occurs.
// Used as an error return value for utils.ValidateErrNoRows.
type ErrNoRows struct {
	msg string
}

func (e *ErrNoRows) Error() string {
	return e.msg
}

func NewErrNoRows(msg string) *ErrNoRows {
	return &ErrNoRows{msg}
}

// PasswordMismatch occurs when passwords do not match
var PasswordMismatch = validation.Errors{"password": fmt.Errorf("las contraseñas ingresadas no coinciden")}
