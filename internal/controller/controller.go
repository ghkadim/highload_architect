package controller

import (
	"context"
	openapi "github.com/ghkadim/highload_architect/generated/go_server/go"
	"github.com/ghkadim/highload_architect/internal/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

type ApiController struct {
	router *mux.Router
}

func NewApiController(apiService *service.ApiService) *ApiController {
	apiController := openapi.NewDefaultApiController(apiService)
	routes := apiController.Routes()

	return &ApiController{
		router: newRouter(routes),
	}
}

func (ac *ApiController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ac.router.ServeHTTP(w, req)
}

func authorize(inner http.Handler, pattern, method string) http.Handler {
	enabled := false
	for _, r := range openapi.AuthorizeRoutes {
		if r.Path == pattern && r.Method == method {
			enabled = true
			break
		}
	}

	if !enabled {
		return inner
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := strings.Trim(r.Header.Get("Authorization"), " ")
		if len(auth) == 0 || !strings.HasPrefix(auth, "Bearer ") {
			openapi.EncodeJSONResponse(nil, func(i int) *int { return &i }(http.StatusUnauthorized), w)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		r = r.WithContext(context.WithValue(r.Context(), "BearerToken", token))
		inner.ServeHTTP(w, r)
	})
}

func newRouter(routes openapi.Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = authorize(route.HandlerFunc, route.Pattern, route.Method)
		handler = logger(handler)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
}

func logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writer := &responseWriter{ResponseWriter: w}

		inner.ServeHTTP(writer, r)

		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.RequestURI,
			writer.code,
			time.Since(start),
		)
	})
}
