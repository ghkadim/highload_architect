package service

import (
	"context"
	"errors"
	openapi "github.com/ghkadim/highload_architect/generated/go"
	"github.com/ghkadim/highload_architect/internal/models"
)

type Storage interface {
	UserRegister(ctx context.Context, user models.User) (string, error)
	UserGet(ctx context.Context, id string) (models.User, error)
	UserSearch(ctx context.Context, firstName, secondName string) ([]models.User, error)
}

type Session interface {
	HashPassword(ctx context.Context, password string) ([]byte, error)
	CompareHashAndPassword(ctx context.Context, hash []byte, password string) (bool, error)
	TokenForUser(ctx context.Context, user models.User) (string, error)
}

type ApiService struct {
	storage Storage
	session Session
}

// NewApiService creates an api service
func NewApiService(
	storage Storage,
	session Session,
) *ApiService {
	return &ApiService{
		storage: storage,
		session: session,
	}
}

// LoginPost -
func (s *ApiService) LoginPost(ctx context.Context, loginPostRequest openapi.LoginPostRequest) (openapi.ImplResponse, error) {
	user, err := s.storage.UserGet(ctx, loginPostRequest.Id)
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

	token, err := s.session.TokenForUser(ctx, user)
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	return openapi.Response(200, openapi.LoginPost200Response{Token: token}), nil
}

// UserGetIdGet -
func (s *ApiService) UserGetIdGet(ctx context.Context, id string) (openapi.ImplResponse, error) {
	user, err := s.storage.UserGet(ctx, id)
	if err != nil {
		if errors.Is(err, models.UserNotFound) {
			return openapi.Response(404, nil), nil
		}

		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	return openapi.Response(200, openapi.User{
		Id:         user.ID,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Age:        user.Age,
		Biography:  user.Biography,
		City:       user.City,
	}), nil
}

// UserRegisterPost -
func (s *ApiService) UserRegisterPost(ctx context.Context, user openapi.UserRegisterPostRequest) (openapi.ImplResponse, error) {
	password, err := s.session.HashPassword(ctx, user.Password)
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	id, err := s.storage.UserRegister(ctx, models.User{
		FirstName:    user.FirstName,
		SecondName:   user.SecondName,
		Age:          user.Age,
		Biography:    user.Biography,
		City:         user.City,
		PasswordHash: password,
	})
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	return openapi.Response(200, openapi.UserRegisterPost200Response{UserId: id}), nil
}

// UserSearchGet -
func (s *ApiService) UserSearchGet(ctx context.Context, firstName string, lastName string) (openapi.ImplResponse, error) {
	users, err := s.storage.UserSearch(ctx, firstName, lastName)
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	apiUsers := make([]openapi.User, 0, len(users))
	for i := range users {
		apiUsers = append(apiUsers, openapi.User{
			Id:         users[i].ID,
			FirstName:  users[i].FirstName,
			SecondName: users[i].SecondName,
			Age:        users[i].Age,
			Biography:  users[i].Biography,
			City:       users[i].City,
		})
	}

	return openapi.Response(200, apiUsers), nil
}
