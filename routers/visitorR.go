package routers

import (
	"net/http"
	"peer2peer/controllers"

	"github.com/gorilla/mux"
)

var v controllers.Visitor

var routes = Routes{
	Route{
		"Get Visitor",
		"GET",
		"/v1/visitor/{id}",
		v.Get,
	},
	Route{
		"Add Visitor",
		"POST",
		"/v1/visitor",
		v.Create,
	},
}

// AddVisitorRoutes : Add Visitors Routes
func AddVisitorRoutes(r *mux.Router) *mux.Router {
	for _, route := range routes {

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
