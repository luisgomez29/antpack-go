package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/luisgomez29/antpack-go/app/auth"
	"github.com/luisgomez29/antpack-go/app/controllers"
	"github.com/luisgomez29/antpack-go/app/repositories"
	apiErrors "github.com/luisgomez29/antpack-go/app/resources/api/errors"
	"github.com/luisgomez29/antpack-go/app/routes"
	"github.com/luisgomez29/antpack-go/app/services"
	"github.com/luisgomez29/antpack-go/pkg/faker"
	"github.com/luisgomez29/antpack-go/pkg/mock/repositories"
	"github.com/luisgomez29/antpack-go/tests"
)

func TestSignUp(t *testing.T) {
	user, password := faker.User()
	testCases := []struct {
		name          string
		body          echo.Map
		buildStubs    func(repo *mockrepo.MockAccountRepository)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: echo.Map{
				"first_name":       user.FirstName,
				"last_name":        user.LastName,
				"email":            user.Email,
				"password":         password,
				"password_confirm": password,
			},
			buildStubs: func(repo *mockrepo.MockAccountRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, user)
			},
		},
		{
			name: "DuplicateEmail",
			body: echo.Map{
				"first_name":       user.FirstName,
				"last_name":        user.LastName,
				"email":            user.Email,
				"password":         password,
				"password_confirm": password,
			},
			buildStubs: func(repo *mockrepo.MockAccountRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Times(1).Return(nil, validation.Errors{})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: echo.Map{
				"first_name":       user.FirstName,
				"last_name":        user.LastName,
				"email":            "invalid_email",
				"password":         password,
				"password_confirm": password,
			},
			buildStubs: func(repo *mockrepo.MockAccountRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "PasswordsMismatch",
			body: echo.Map{
				"first_name":       user.FirstName,
				"last_name":        user.LastName,
				"email":            user.Email,
				"password":         password,
				"password_confirm": "password",
			},
			buildStubs: func(repo *mockrepo.MockAccountRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: echo.Map{
				"first_name":       user.FirstName,
				"last_name":        user.LastName,
				"email":            user.Email,
				"password":         password,
				"password_confirm": "123456",
			},
			buildStubs: func(repo *mockrepo.MockAccountRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NumericPassword",
			body: echo.Map{
				"first_name":       user.FirstName,
				"last_name":        user.LastName,
				"email":            user.Email,
				"password":         "12345678",
				"password_confirm": password,
			},
			buildStubs: func(repo *mockrepo.MockAccountRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "PasswordWithOnlyLetters",
			body: echo.Map{
				"first_name":       user.FirstName,
				"last_name":        user.LastName,
				"email":            user.Email,
				"password":         "abcdefghijk",
				"password_confirm": password,
			},
			buildStubs: func(repo *mockrepo.MockAccountRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockrepo.NewMockAccountRepository(ctrl)
			tc.buildStubs(repo)

			//Repositories
			usersRepo := repositories.NewUserRepository(tests.DB)
			authRepo := repositories.NewAuthRepository(tests.DB, usersRepo)

			// Authentication service
			authn := auth.NewAuth(authRepo)

			// Services
			accountsService := services.NewAccountsService(authn, repo)

			// Routes
			routes.Accounts(tests.Route(tests.E), controllers.NewAccountsController(authn, accountsService))

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewReader(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			tests.E.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestLogin(t *testing.T) {
	user, password := faker.User()
	testCases := []struct {
		name          string
		body          echo.Map
		buildStubs    func(repo *mockrepo.MockAccountRepository)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: echo.Map{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(repo *mockrepo.MockAccountRepository) {
				repo.EXPECT().FindUser(gomock.Eq(user.Email)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, user)
			},
		},
		{
			name: "NonExistentAccount",
			body: echo.Map{
				"email":    "email@example.com",
				"password": password,
			},
			buildStubs: func(repo *mockrepo.MockAccountRepository) {
				repo.EXPECT().FindUser(gomock.Any()).Times(1).Return(nil, apiErrors.NewErrNoRows(""))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: echo.Map{
				"email":    user.Email,
				"password": "pass",
			},
			buildStubs: func(repo *mockrepo.MockAccountRepository) {
				repo.EXPECT().FindUser(gomock.Eq(user.Email)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockrepo.NewMockAccountRepository(ctrl)
			tc.buildStubs(repo)

			//Repositories
			usersRepo := repositories.NewUserRepository(tests.DB)
			authRepo := repositories.NewAuthRepository(tests.DB, usersRepo)

			// Authentication service
			authn := auth.NewAuth(authRepo)

			// Services
			accountsService := services.NewAccountsService(authn, repo)

			// Routes
			routes.Accounts(tests.Route(tests.E), controllers.NewAccountsController(authn, accountsService))

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			tests.E.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}
