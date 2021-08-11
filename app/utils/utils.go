package utils

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

// Regular expressions
var (
	ReDigit   = regexp.MustCompile("^[0-9]+$")
	ReLetters = regexp.MustCompile("^[a-zA-ZÁ-ÿ ?]*$")
)

// PasswordValidator verify that the password has letters and numbers.
func PasswordValidator(value interface{}) error {
	s, _ := value.(string)
	if ReDigit.Match([]byte(s)) || ReLetters.Match([]byte(s)) {
		return errors.New("la contraseña debe tener letras y números")
	}
	return nil
}

// TimeDuration convert string to time.Duration.
func TimeDuration(t string) (time.Duration, error) {
	tc, err := strconv.Atoi(t)
	if err != nil {
		return 0, err
	}
	return time.Duration(tc), nil
}
