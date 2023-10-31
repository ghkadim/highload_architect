package counter

import (
	"context"
	"errors"

	"github.com/ghkadim/highload_architect/internal/models"
)

type storage interface {
	CounterAdd(ctx context.Context, userID models.UserID, id string, value int64) (int64, error)
	CounterRead(ctx context.Context, userID models.UserID, id string) (int64, error)
}

type Service interface {
	CounterAdd(ctx context.Context, userID models.UserID, id string, value int64) (int64, error)
	CounterRead(ctx context.Context, userID models.UserID, id string) (int64, error)
}

type service struct {
	storage storage
}

func NewService(
	storage storage,
) Service {
	return &service{
		storage: storage,
	}
}

func (s *service) CounterAdd(ctx context.Context, userID models.UserID, id string, value int64) (int64, error) {
	return s.storage.CounterAdd(ctx, userID, id, value)
}

func (s *service) CounterRead(ctx context.Context, userID models.UserID, id string) (int64, error) {
	val, err := s.storage.CounterRead(ctx, userID, id)
	if err != nil {
		if errors.Is(err, models.ErrCounterNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return val, nil
}
