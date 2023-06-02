package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leomarzochi/facebooklike/cmd/middlewares"
)

type Route struct {
	URI               string
	Method            string
	Function          func(http.ResponseWriter, *http.Request)
	hasAuthentication bool
}

// Generate a new Gorilla/Mux router
func Generate() *mux.Router {
	r := mux.NewRouter()
	return Configure(r)
}

func Configure(r *mux.Router) *mux.Router {
	routes := UserRoutes
	routes = append(routes, authRoutes)

	for _, route := range routes {
		if route.hasAuthentication {
			r.
				HandleFunc(route.URI, middlewares.Logger(
					middlewares.Authenticate(route.Function),
				)).
				Methods(route.Method)
		} else {
			r.
				HandleFunc(route.URI, middlewares.Logger(route.Function)).
				Methods(route.Method)
		}
	}

	return r
}
