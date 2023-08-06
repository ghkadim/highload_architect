package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	Routes() []Route
}

type server struct {
	router *mux.Router
}

func NewServer(routers ...Router) *server {
	routes := make([]Route, 0)
	for _, router := range routers {
		routes = append(routes, router.Routes()...)
	}

	return &server{
		router: newRouter(routes),
	}
}

func (ac *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ac.router.ServeHTTP(w, req)
}

func newRouter(routes []Route) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		if route.Authorize {
			handler = authorize(handler)
		}
		handler = logger(handler)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
