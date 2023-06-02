package router

import (
	"net/http"

	"github.com/leomarzochi/facebooklike/cmd/controllers"
)

var authRoutes = Route{
	URI:               "/login",
	hasAuthentication: false,
	Method:            http.MethodPost,
	Function:          controllers.Login,
}
