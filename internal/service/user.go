package service

import (
	"context"
	"errors"
	openapi "github.com/ghkadim/highload_architect/generated/go_server/go"
	"github.com/ghkadim/highload_architect/internal/models"
)

// UserGetIdGet -
func (s *ApiService) UserGetIdGet(ctx context.Context, id string) (openapi.ImplResponse, error) {
	user, err := s.readStorage().UserGet(ctx, models.UserID(id))
	if err != nil {
		if errors.Is(err, models.UserNotFound) {
			return openapi.Response(404, nil), nil
		}

		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	return openapi.Response(200, openapi.User{
		Id:         string(user.ID),
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Age:        valueOrDefault(user.Age),
		Biography:  valueOrDefault(user.Biography),
		City:       valueOrDefault(user.City),
	}), nil
}

// UserRegisterPost -
func (s *ApiService) UserRegisterPost(ctx context.Context, user openapi.UserRegisterPostRequest) (openapi.ImplResponse, error) {
	password, err := s.session.HashPassword(ctx, user.Password)
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	id, err := s.master.UserRegister(ctx, models.User{
		FirstName:    user.FirstName,
		SecondName:   user.SecondName,
		Age:          &user.Age,
		Biography:    &user.Biography,
		City:         &user.City,
		PasswordHash: password,
	})
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	return openapi.Response(200, openapi.UserRegisterPost200Response{UserId: string(id)}), nil
}

// UserSearchGet -
func (s *ApiService) UserSearchGet(ctx context.Context, firstName string, lastName string) (openapi.ImplResponse, error) {
	if firstName == "" && lastName == "" {
		return openapi.Response(400, "last_name or first_name should not be empty"), nil
	}

	users, err := s.readStorage().UserSearch(ctx, firstName, lastName)
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	apiUsers := make([]openapi.User, 0, len(users))
	for i := range users {
		apiUsers = append(apiUsers, openapi.User{
			Id:         string(users[i].ID),
			FirstName:  users[i].FirstName,
			SecondName: users[i].SecondName,
			Age:        valueOrDefault(users[i].Age),
			Biography:  valueOrDefault(users[i].Biography),
			City:       valueOrDefault(users[i].City),
		})
	}

	return openapi.Response(200, apiUsers), nil
}
