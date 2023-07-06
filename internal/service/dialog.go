package service

import (
	"context"

	"github.com/ghkadim/highload_architect/internal/models"
)

func (s *service) DialogSend(ctx context.Context, message models.DialogMessage) error {
	err := s.master.DialogSend(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error) {
	messages, err := s.readStorage().DialogList(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
