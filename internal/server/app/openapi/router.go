package openapi

import (
	"net/http"

	"github.com/ghkadim/highload_architect/generated/app/go_server/go"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/server"
)

type router struct {
	router openapi.Router
}

func NewRouter(apiService openapi.DefaultApiServicer) *router {
	return &router{
		router: openapi.NewDefaultApiController(
			apiService,
			openapi.WithDefaultApiErrorHandler(errorHandler),
		),
	}
}

func (r *router) Routes() []server.Route {
	routes := make([]server.Route, 0)
	for _, route := range r.router.Routes() {
		authorize := false
		for _, rt := range openapi.AuthorizeRoutes {
			if rt.Path == route.Pattern && rt.Method == route.Method {
				authorize = true
				break
			}
		}
		routes = append(routes, server.Route{
			Name:        route.Name,
			Method:      route.Method,
			Pattern:     route.Pattern,
			Authorize:   authorize,
			HandlerFunc: route.HandlerFunc,
		})
	}
	return routes
}

func errorHandler(w http.ResponseWriter, r *http.Request, err error, result *openapi.ImplResponse) {
	if err != nil {
		logger.Error("%v", err)
	}
	openapi.DefaultErrorHandler(w, r, err, result)
}
