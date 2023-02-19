package service

import (
	"context"
	"errors"
	openapi "github.com/ghkadim/highload_architect/generated/go_server/go"
	"github.com/ghkadim/highload_architect/internal/models"
)

// LoginPost -
func (s *ApiService) LoginPost(ctx context.Context, loginPostRequest openapi.LoginPostRequest) (openapi.ImplResponse, error) {
	user, err := s.master.UserGet(ctx, loginPostRequest.Id)
	if err != nil {
		if errors.Is(err, models.UserNotFound) {
			return openapi.Response(404, nil), nil
		}
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	equal, err := s.session.CompareHashAndPassword(ctx, user.PasswordHash, loginPostRequest.Password)
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	if !equal {
		return openapi.Response(404, nil), err
	}

	token, err := s.session.TokenForUser(ctx, user.ID)
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	return openapi.Response(200, openapi.LoginPost200Response{Token: token}), nil
}
