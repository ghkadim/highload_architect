package websocket

import (
	"context"
	"net/http"

	"nhooyr.io/websocket"

	"github.com/ghkadim/highload_architect/internal/controller"
	"github.com/ghkadim/highload_architect/internal/logger"
)

type router struct {
	routes []controller.Route
}

func NewRouter() *router {
	return &router{}
}

func (r *router) Routes() []controller.Route {
	routes := make([]controller.Route, 0)
	return routes
}

func (r *router) AddRoute(
	name string,
	pattern string,
	authorize bool,
	handler func(ctx context.Context, c *websocket.Conn),
) {
	f := func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
		if err != nil {
			logger.Error("Failed to accept ws: %v", err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "ws error")
		handler(r.Context(), c)
	}

	r.routes = append(r.routes, controller.Route{
		Name:        name,
		Method:      http.MethodGet,
		Pattern:     pattern,
		Authorize:   authorize,
		HandlerFunc: f,
	})
}
