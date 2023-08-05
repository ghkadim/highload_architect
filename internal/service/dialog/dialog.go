package dialog

import (
	"context"

	"github.com/ghkadim/highload_architect/internal/models"
)

type storage interface {
	DialogSend(ctx context.Context, message models.DialogMessage) error
	DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error)
}

type Service struct {
	storage storage
}

func NewService(
	storage storage,
) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) DialogSend(ctx context.Context, message models.DialogMessage) error {
	err := s.storage.DialogSend(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error) {
	messages, err := s.storage.DialogList(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
