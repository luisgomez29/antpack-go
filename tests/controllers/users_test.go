package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/luisgomez29/antpack-go/app/auth"
	"github.com/luisgomez29/antpack-go/app/controllers"
	"github.com/luisgomez29/antpack-go/app/models"
	"github.com/luisgomez29/antpack-go/app/repositories"
	apiErrors "github.com/luisgomez29/antpack-go/app/resources/api/errors"
	"github.com/luisgomez29/antpack-go/app/routes"
	"github.com/luisgomez29/antpack-go/app/services"
	"github.com/luisgomez29/antpack-go/pkg/faker"
	"github.com/luisgomez29/antpack-go/pkg/mock/repositories"
	"github.com/luisgomez29/antpack-go/tests"
)

func TestAll(t *testing.T) {
	n := 3
	users := make([]*models.User, n)
	for i := 0; i < n; i++ {
		users[i], _ = faker.User()
		users[i].Password = ""
	}

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request)
		buildStubs    func(repo *mockrepo.MockUserRepository)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request) {
				tests.AddAuthorization(request, users[0], time.Minute)
			},
			buildStubs: func(repo *mockrepo.MockUserRepository) {
				repo.EXPECT().All().Times(1).Return(users, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUsers(t, recorder.Body, users)
			},
		},
		{
			name: "JWTMissing",
			setupAuth: func(request *http.Request) {
			},
			buildStubs: func(repo *mockrepo.MockUserRepository) {
				repo.EXPECT().All().Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(request *http.Request) {
				tests.AddAuthorization(request, users[0], -time.Minute)
			},
			buildStubs: func(repo *mockrepo.MockUserRepository) {
				repo.EXPECT().All().Times(0)
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

			repo := mockrepo.NewMockUserRepository(ctrl)
			tc.buildStubs(repo)

			//Repositories
			authRepo := repositories.NewAuthRepository(tests.DB, repo)

			// Authentication service
			authn := auth.NewAuth(authRepo)

			// Services
			usersService := services.NewUsersService(authn, repo)

			// Routes
			routes.Users(tests.Route(tests.E), controllers.NewUsersController(usersService))

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			tc.setupAuth(req)
			tests.E.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestGet(t *testing.T) {
	user, _ := faker.User()
	user.Password = ""

	testCases := []struct {
		name          string
		userID        uint
		setupAuth     func(request *http.Request)
		buildStubs    func(repo *mockrepo.MockUserRepository)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID,
			setupAuth: func(request *http.Request) {
				tests.AddAuthorization(request, user, time.Minute)
			},
			buildStubs: func(repo *mockrepo.MockUserRepository) {
				repo.EXPECT().Get(gomock.Eq(user.ID)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:   "JWTMissing",
			userID: user.ID,
			setupAuth: func(request *http.Request) {
			},
			buildStubs: func(repo *mockrepo.MockUserRepository) {
				repo.EXPECT().Get(gomock.Eq(user.ID)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "ExpiredToken",
			userID: user.ID,
			setupAuth: func(request *http.Request) {
				tests.AddAuthorization(request, user, -time.Minute)
			},
			buildStubs: func(repo *mockrepo.MockUserRepository) {
				repo.EXPECT().Get(gomock.Eq(user.ID)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			userID: user.ID,
			setupAuth: func(request *http.Request) {
				tests.AddAuthorization(request, user, time.Minute)
			},
			buildStubs: func(repo *mockrepo.MockUserRepository) {
				repo.EXPECT().Get(gomock.Eq(user.ID)).Times(1).Return(nil, apiErrors.NewErrNoRows(""))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockrepo.NewMockUserRepository(ctrl)
			tc.buildStubs(repo)

			//Repositories
			authRepo := repositories.NewAuthRepository(tests.DB, repo)

			// Authentication service
			authn := auth.NewAuth(authRepo)

			// Services
			usersService := services.NewUsersService(authn, repo)

			// Routes
			routes.Users(tests.Route(tests.E), controllers.NewUsersController(usersService))

			url := fmt.Sprintf("/api/v1/users/%d", tc.userID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			tc.setupAuth(req)
			tests.E.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}
