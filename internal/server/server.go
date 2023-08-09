package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	log "github.com/ghkadim/highload_architect/internal/logger"
)

type Router interface {
	Routes() []Route
}

type Server struct {
	router *mux.Router
}

func NewServer(routers ...Router) *Server {
	routes := make([]Route, 0)
	for _, router := range routers {
		routes = append(routes, router.Routes()...)
	}

	return &Server{
		router: newRouter(routes),
	}
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

func (s *Server) ListenAndServe(addr string) {
	srv := http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: %s\n", err)
		}
	}()

	<-done
	log.Info("Server stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown failed:%+v", err)
	}
	log.Info("Server exited properly")
}
