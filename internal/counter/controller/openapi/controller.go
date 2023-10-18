package openapi

import (
	"context"
	"errors"

	openapi "github.com/ghkadim/highload_architect/generated/counter/go_server/go"
	"github.com/ghkadim/highload_architect/internal/models"
)

type service interface {
	CounterRead(ctx context.Context, userID models.UserID, id string) (int64, error)
}

type session interface {
	ParseToken(ctx context.Context, tokenStr string) (models.UserID, error)
}

var _ openapi.DefaultApiServicer = &apiController{}

type apiController struct {
	service service
	session session
}

func NewController(
	service service,
	session session,
) *apiController {
	return &apiController{
		service: service,
		session: session,
	}
}

func (c *apiController) CounterCounterIdGet(ctx context.Context, id string) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}

	value, err := c.service.CounterRead(ctx, userID, id)
	if err != nil {
		return errorResponse(err)
	}

	return successResponse(openapi.CounterCounterIdGet200Response{
		Value: value,
	})
}

func errorResponse(err error) (openapi.ImplResponse, error) {
	switch {
	case errors.Is(err, models.ErrUserNotFound):
		return openapi.Response(404, nil), err
	case errors.Is(err, models.ErrPostNotFound):
		return openapi.Response(404, nil), err
	default:
		return openapi.Response(500, nil), err
	}
}

func successResponse(body interface{}) (openapi.ImplResponse, error) {
	return openapi.Response(200, body), nil
}

func bearerToken(ctx context.Context) string {
	val := ctx.Value(models.BearerTokenCtxKey)
	if val == nil {
		return ""
	}
	return val.(string)
}
