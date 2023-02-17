package service

import (
	"context"
	"errors"
	openapi "github.com/ghkadim/highload_architect/generated/go"
	"github.com/ghkadim/highload_architect/internal/models"
	"sync/atomic"
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
	master     Storage
	replicas   []Storage
	replicaNum atomic.Int32
	session    Session
}

// NewApiService creates an api service
func NewApiService(
	master Storage,
	replicas []Storage,
	session Session,
) *ApiService {
	return &ApiService{
		master:   master,
		replicas: replicas,
		session:  session,
	}
}

func valueOrDefault[V any](value *V) V {
	if value == nil {
		return *new(V)
	}
	return *value
}

func (s *ApiService) readStorage() Storage {
	if len(s.replicas) != 0 {
		return s.replicas[int(s.replicaNum.Add(1))%len(s.replicas)]
	}
	return s.master
}

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

	token, err := s.session.TokenForUser(ctx, user)
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	return openapi.Response(200, openapi.LoginPost200Response{Token: token}), nil
}

// UserGetIdGet -
func (s *ApiService) UserGetIdGet(ctx context.Context, id string) (openapi.ImplResponse, error) {
	user, err := s.readStorage().UserGet(ctx, id)
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

	return openapi.Response(200, openapi.UserRegisterPost200Response{UserId: id}), nil
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
			Id:         users[i].ID,
			FirstName:  users[i].FirstName,
			SecondName: users[i].SecondName,
			Age:        valueOrDefault(users[i].Age),
			Biography:  valueOrDefault(users[i].Biography),
			City:       valueOrDefault(users[i].City),
		})
	}

	return openapi.Response(200, apiUsers), nil
}
