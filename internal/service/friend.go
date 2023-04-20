package service

import (
	"context"
	"net/http"

	openapi "github.com/ghkadim/highload_architect/generated/go_server/go"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

// FriendDeleteUserIdPut -
func (s *ApiService) FriendDeleteUserIdPut(ctx context.Context, friendUserID string) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := s.session.ParseToken(ctx, token)
	if err != nil {
		logger.Error("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	err = s.master.FriendDelete(ctx, userID, models.UserID(friendUserID))
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	s.cache.FriendDelete(userID, models.UserID(friendUserID))
	return openapi.Response(http.StatusOK, nil), nil
}

// FriendSetUserIdPut -
func (s *ApiService) FriendSetUserIdPut(ctx context.Context, friendUserID string) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := s.session.ParseToken(ctx, token)
	if err != nil {
		logger.Error("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	err = s.master.FriendAdd(ctx, userID, models.UserID(friendUserID))
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	s.cache.FriendAdd(userID, models.UserID(friendUserID))
	return openapi.Response(http.StatusOK, nil), nil
}
