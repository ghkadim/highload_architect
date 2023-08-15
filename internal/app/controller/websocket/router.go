package websocket

import (
	"context"
	"net/http"

	"nhooyr.io/websocket"

	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/server"
)

type router struct {
	routes []server.Route
}

func NewRouter(controller *wsController) *router {
	r := &router{}
	r.AddRoute(
		"postFeedPosted",
		"/post/feed/posted",
		true,
		controller.PostFeedPosted)
	return r
}

func (r *router) Routes() []server.Route {
	return r.routes
}

func (r *router) AddRoute(
	name string,
	pattern string,
	authorize bool,
	handler func(ctx context.Context, r *http.Request, c *websocket.Conn) error,
) {
	f := func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
		if err != nil {
			logger.Errorf("Failed to accept ws: %v", err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "ws error")

		err = handler(r.Context(), r, c)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				return
			}
			logger.Errorf("Websocket failed: %v", err)
		}
		c.Close(websocket.StatusNormalClosure, "done")
	}

	r.routes = append(r.routes, server.Route{
		Name:        name,
		Method:      http.MethodGet,
		Pattern:     pattern,
		Authorize:   authorize,
		HandlerFunc: f,
	})
}
