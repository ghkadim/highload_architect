package user

import (
	"context"

	"github.com/ghkadim/highload_architect/internal/models"
)

var (
	emptyUserID = models.UserID("")
)

type storage interface {
	UserRegister(ctx context.Context, user models.User) (models.UserID, error)
	UserGet(ctx context.Context, id models.UserID) (models.User, error)
	UserSearch(ctx context.Context, firstName, secondName string) ([]models.User, error)
}

type session interface {
	HashPassword(ctx context.Context, password string) ([]byte, error)
	CompareHashAndPassword(ctx context.Context, hash []byte, password string) (bool, error)
	TokenForUser(ctx context.Context, userID models.UserID) (string, error)
}

type Service struct {
	storage storage
	session session
}

func NewService(
	storage storage,
	session session,
) *Service {
	return &Service{
		storage: storage,
		session: session,
	}
}

func (s *Service) UserRegister(ctx context.Context, user models.User, password string) (models.UserID, error) {
	passwordHash, err := s.session.HashPassword(ctx, password)
	if err != nil {
		return emptyUserID, err
	}

	user.PasswordHash = passwordHash
	id, err := s.storage.UserRegister(ctx, user)
	if err != nil {
		return emptyUserID, err
	}
	return id, nil
}

func (s *Service) UserGet(ctx context.Context, id models.UserID) (models.User, error) {
	user, err := s.storage.UserGet(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *Service) UserSearch(ctx context.Context, firstName, secondName string) ([]models.User, error) {
	users, err := s.storage.UserSearch(ctx, firstName, secondName)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) UserLogin(ctx context.Context, id models.UserID, password string) (string, error) {
	user, err := s.storage.UserGet(ctx, id)
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
