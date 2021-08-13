package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/luisgomez29/antpack-go/app/auth"
	"github.com/luisgomez29/antpack-go/app/models"
	"github.com/luisgomez29/antpack-go/tests"
)

// TestMain runs all the tests within the package.
func TestMain(t *testing.M) {
	tests.Init()
	os.Exit(t.Run())
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, user *models.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var jwtResponse auth.JWTResponse
	err = json.Unmarshal(data, &jwtResponse)
	require.NoError(t, err)

	require.Equal(t, user, jwtResponse.User)
	require.Empty(t, jwtResponse.User.Password)
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user *models.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser *models.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user, gotUser)
	require.Empty(t, gotUser.Password)
}

func requireBodyMatchUsers(t *testing.T, body *bytes.Buffer, user []*models.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUsers map[string][]*models.User
	err = json.Unmarshal(data, &gotUsers)
	require.NoError(t, err)

	require.Equal(t, user, gotUsers["results"])
}
