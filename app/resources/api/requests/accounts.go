package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/luisgomez29/antpack-go/app/utils"
)

// Validation of user model fields entered by user.
var (
	firstNameRule = []validation.Rule{
		validation.Required.Error("el nombre es requerido"),
		validation.Match(utils.ReLetters).Error("el nombre debe tener solo letras (A-Z)"),
		validation.Length(2, 40).Error("el nombre debe tener entre 2 y 40 caracteres"),
	}

	lastNameRule = []validation.Rule{
		validation.Required.Error("los apellidos son requeridos"),
		validation.Match(utils.ReLetters).Error("los apellidos deben tener solo letras (A-Z)"),
		validation.Length(2, 40).Error("el nombre debe tener entre 5 y 40 caracteres"),
	}

	emailRule = []validation.Rule{
		validation.Required.Error("el correo electrónico es requerido"),
		is.Email.Error("ingrese una dirección de correo electrónico válida"),
	}

	passwordRule = []validation.Rule{
		validation.Required.Error("la contraseña es requerida"),
		validation.Length(8, 25).Error(
			"la contraseña debe tener entre 8 y 25 caracteres",
		),
		validation.By(utils.PasswordValidator),
	}

	passwordConfirmationRule = []validation.Rule{
		validation.Required.Error("la contraseña de confirmación es requerida"),
	}
)

// Passwords represents the password and the confirmation password.
type Passwords struct {
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func (p Passwords) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Password, passwordRule...),
		validation.Field(&p.PasswordConfirm, passwordConfirmationRule...),
	)
}

// ------ REQUESTS

// SignUpRequest represents the user's request for the creation of an account.
type SignUpRequest struct {
	Passwords

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (rs SignUpRequest) Validate() error {
	return validation.ValidateStruct(&rs,
		validation.Field(&rs.Passwords),
		validation.Field(&rs.FirstName, firstNameRule...),
		validation.Field(&rs.LastName, lastNameRule...),
		validation.Field(&rs.Email, emailRule...),
	)
}

// LoginRequest represents the user's login request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func (rs LoginRequest) Validate() error {
	return validation.ValidateStruct(&rs,
		validation.Field(&rs.Email, emailRule...),
		validation.Field(&rs.Password, validation.Required.Error("la contraseña es requerida")),
	)
}
