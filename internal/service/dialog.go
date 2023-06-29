package service

import (
	"context"
	"net/http"

	openapi "github.com/ghkadim/highload_architect/generated/go_server/go"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

// DialogUserIdListGet -
func (s *ApiService) DialogUserIdListGet(ctx context.Context, userID2 string) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID1, err := s.session.ParseToken(ctx, token)
	if err != nil {
		logger.Error("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	messages, err := s.master.DialogList(ctx, userID1, models.UserID(userID2))
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	dialogMessages := make([]openapi.DialogMessage, 0, len(messages))
	for _, msg := range messages {
		dialogMessages = append(dialogMessages, openapi.DialogMessage{
			From: string(msg.From),
			To:   string(msg.To),
			Text: msg.Text,
		})
	}
	return openapi.Response(http.StatusOK, dialogMessages), nil
}

// DialogUserIdSendPost -
func (s *ApiService) DialogUserIdSendPost(
	ctx context.Context,
	toUserID string,
	dialogUserIdSendPostRequest openapi.DialogUserIdSendPostRequest,
) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	fromUserID, err := s.session.ParseToken(ctx, token)
	if err != nil {
		logger.Error("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	err = s.master.DialogSend(
		ctx,
		models.DialogMessage{From: fromUserID, To: models.UserID(toUserID), Text: dialogUserIdSendPostRequest.Text})
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}
	return openapi.Response(http.StatusOK, nil), nil
}
