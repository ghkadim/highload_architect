package openapi

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	openapi "github.com/ghkadim/highload_architect/generated/dialog/go_server/go"
	"github.com/ghkadim/highload_architect/internal/models"
)

type service interface {
	DialogSend(ctx context.Context, message models.DialogMessage) (models.DialogMessageID, error)
	DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error)
	DialogMessageRead(ctx context.Context, userID models.UserID, messageID models.DialogMessageID) error
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

// DialogUserIdListGet -
func (c *apiController) DialogUserIdListGet(ctx context.Context, userID2 string) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID1, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}

	messages, err := c.service.DialogList(ctx, userID1, models.UserID(userID2))
	if err != nil {
		return errorResponse(err)
	}

	dialogMessages := make([]openapi.DialogMessage, 0, len(messages))
	for _, msg := range messages {
		dialogMessages = append(dialogMessages, openapi.DialogMessage{
			Id:   fmt.Sprintf("%d", msg.ID),
			From: string(msg.From),
			To:   string(msg.To),
			Text: msg.Text,
		})
	}
	return successResponse(dialogMessages)
}

// DialogUserIdSendPost -
func (c *apiController) DialogUserIdSendPost(
	ctx context.Context,
	toUserID string,
	dialogUserIdSendPostRequest openapi.DialogUserIdSendPostRequest,
) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	fromUserID, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}

	id, err := c.service.DialogSend(
		ctx,
		models.DialogMessage{
			From: fromUserID,
			To:   models.UserID(toUserID),
			Text: dialogUserIdSendPostRequest.Text,
		})
	if err != nil {
		return errorResponse(err)
	}
	return successResponse(openapi.DialogUserIdSendPost200Response{MessageId: fmt.Sprintf("%d", id)})
}

func (c *apiController) DialogMessageMessageIdReadPut(
	ctx context.Context,
	messageId string,
) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	fromUserID, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}
	msgId, err := strconv.ParseInt(messageId, 10, 64)
	if err != nil {
		return errorResponse(err)
	}
	err = c.service.DialogMessageRead(ctx, fromUserID, models.DialogMessageID(msgId))
	if err != nil {
		return errorResponse(err)
	}
	return successResponse(nil)
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
