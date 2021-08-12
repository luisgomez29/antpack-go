package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/luisgomez29/antpack-go/app/models"
)

var (
	mockDB = map[string]*models.User{
		"jon@labstack.com": {
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			Password:  "lg123456",
		},
	}
	userJSON = `{
		"first_name": "Test",
		"last_name": "User",
		"email": "test@example.com",
		"password": "lg123456",
		"password_confirm": "lg123456"
	}`
)

func TestSignUp(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestLogin(t *testing.T) {

	data := `{
		"email": "leidy@gmail.com",
		"password": "lg123456"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}
