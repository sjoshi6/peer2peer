package routers

import (
	"net/http"
	"peer2peer/controllers"

	"github.com/gorilla/mux"
)

var auth controllers.Auth

var authroutes = Routes{

	Route{
		"SignUp",
		"POST",
		"/v1/signup",
		auth.SignUpHandler,
	},
	Route{
		"Login",
		"POST",
		"/v1/login",
		auth.LoginHandler,
	},
	// Needs to be here as it does not need auth tokens for access
	Route{
		"Add Visitor",
		"POST",
		"/v1/visitor",
		v.Create,
	},
}

// AddAuthRoutes : Add Auth Routes
func AddAuthRoutes(r *mux.Router) *mux.Router {
	for _, route := range authroutes {

		var handler http.Handler
		handler = route.HandlerFunc
		handler = APILogger(handler, route.Name)

		r.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return r
}
