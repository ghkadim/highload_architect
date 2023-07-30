package service

import (
	"context"

	"github.com/ghkadim/highload_architect/internal/models"
)

var (
	emptyUserID = models.UserID("")
)

func (s *service) UserRegister(ctx context.Context, user models.User, password string) (models.UserID, error) {
	passwordHash, err := s.session.HashPassword(ctx, password)
	if err != nil {
		return emptyUserID, err
	}

	user.PasswordHash = passwordHash
	id, err := s.master.UserRegister(ctx, user)
	if err != nil {
		return emptyUserID, err
	}
	return id, nil
}

func (s *service) UserGet(ctx context.Context, id models.UserID) (models.User, error) {
	user, err := s.readStorage().UserGet(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *service) UserSearch(ctx context.Context, firstName, secondName string) ([]models.User, error) {
	users, err := s.readStorage().UserSearch(ctx, firstName, secondName)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *service) UserLogin(ctx context.Context, id models.UserID, password string) (string, error) {
	user, err := s.master.UserGet(ctx, id)
	if err != nil {
		return "", err
	}

	equal, err := s.session.CompareHashAndPassword(ctx, user.PasswordHash, password)
	if err != nil {
		return "", err
	}

	if !equal {
		return "", models.ErrUnauthorized
	}

	token, err := s.session.TokenForUser(ctx, user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
