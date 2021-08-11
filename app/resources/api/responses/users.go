package responses

import "github.com/luisgomez29/antpack-go/app/models"

// UserResponse lists the fields to return for the user types.
func UserResponse(role string, user *models.User) *models.User {
	user.Password = ""
	return user
}

// UserManyResponse lists the fields to return for the slice user types.
func UserManyResponse(role string, users []*models.User) []*models.User {
	for i, user := range users {
		users[i] = UserResponse(role, user)
	}
	return users
}
