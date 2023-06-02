package router

import (
	"net/http"

	"github.com/leomarzochi/facebooklike/cmd/controllers"
)

var UserRoutes = []Route{
	{
		URI:               "/user",
		Method:            http.MethodPost,
		Function:          controllers.UserCreate,
		hasAuthentication: false,
	},
	{
		URI:               "/users",
		Method:            http.MethodGet,
		Function:          controllers.UserListAll,
		hasAuthentication: true,
	},
	{
		URI:               "/user/{id}",
		Method:            http.MethodGet,
		Function:          controllers.UserList,
		hasAuthentication: true,
	},
	{
		URI:               "/user/{id}",
		Method:            http.MethodPut,
		Function:          controllers.UserUpdate,
		hasAuthentication: true,
	},
	{
		URI:               "/user/{id}",
		Method:            http.MethodDelete,
		Function:          controllers.UserDelete,
		hasAuthentication: true,
	},
	{
		URI:               "/user/{followID}/follow",
		Method:            http.MethodGet,
		Function:          controllers.UserFollow,
		hasAuthentication: true,
	},
	{
		URI:               "/user/{followID}/unfollow",
		Method:            http.MethodGet,
		Function:          controllers.UserUnfollow,
		hasAuthentication: true,
	},
	{
		URI:               "/user/{userID}/followed",
		Method:            http.MethodGet,
		Function:          controllers.UserFollows, //Users that i follow
		hasAuthentication: true,
	},
	{
		URI:               "/user/{userID}/followers",
		Method:            http.MethodGet,
		Function:          controllers.UserFollowers, //Users that follow me
		hasAuthentication: true,
	},
	{
		URI:               "/user/change-password",
		Method:            http.MethodPost,
		Function:          controllers.UserChangePassword, //Users that follow me
		hasAuthentication: true,
	},
}
