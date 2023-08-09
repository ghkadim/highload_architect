package websocket

import (
	"context"
	"encoding/json"
	"net/http"

	"nhooyr.io/websocket"

	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
	"github.com/ghkadim/highload_architect/internal/result"
)

type service interface {
	PostFeedPosted(ctx context.Context, subscriber models.UserID) <-chan result.Result[models.Post]
}

type session interface {
	ParseToken(ctx context.Context, tokenStr string) (models.UserID, error)
}

type wsController struct {
	service service
	session session
}

func NewController(service service, session session) *wsController {
	return &wsController{
		service: service,
		session: session,
	}
}

func (c *wsController) PostFeedPosted(ctx context.Context, r *http.Request, w *websocket.Conn) error {
	subscriberID, err := c.userFromToken(ctx)
	if err != nil {
		return models.ErrUnauthorized
	}

	logger.Info("PostFeedPosted: starting new subscription for userID=%s", subscriberID)
	ctx = w.CloseRead(ctx)
	postCh := c.service.PostFeedPosted(ctx, subscriberID)
	for {
		select {
		case res, ok := <-postCh:
			if !ok {
				return nil
			}
			newPost, err := res.Value()
			if err != nil {
				return err
			}

			logger.Debug("PostFeedPosted: new post for userID=%s", subscriberID)
			writer, err := w.Writer(ctx, websocket.MessageText)
			if err != nil {
				return err
			}
			err = json.NewEncoder(writer).Encode(post{
				Id:           string(newPost.ID),
				Text:         newPost.Text,
				AuthorUserId: string(newPost.AuthorID),
			})
			if err != nil {
				writer.Close()
				return err
			}
			writer.Close()
		case <-ctx.Done():
			logger.Info("PostFeedPosted: connection closed for userID=%s", subscriberID)
			return nil
		}
	}
}

func (c *wsController) userFromToken(ctx context.Context) (models.UserID, error) {
	token := ctx.Value(models.BearerTokenCtxKey)
	if token == nil {
		return "", models.ErrUnauthorized
	}
	userID, err := c.session.ParseToken(ctx, token.(string))
	if err != nil {
		return "", err
	}
	return userID, nil
}
