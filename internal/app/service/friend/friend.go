package friend

import (
	"context"

	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

type storage interface {
	FriendAdd(ctx context.Context, userID1, userID2 models.UserID) error
	FriendDelete(ctx context.Context, userID1, userID2 models.UserID) error
}

type cache interface {
	FriendAdd(userID1, userID2 models.UserID)
	FriendDelete(userID1, userID2 models.UserID)
}

type eventPublisher interface {
	FriendAdd(userID, friendID models.UserID) error
	FriendDelete(userID, friendID models.UserID) error
}

type Service struct {
	storage        storage
	cache          cache
	eventPublisher eventPublisher
}

func NewService(
	storage storage,
	cache cache,
	eventPublisher eventPublisher,
) *Service {
	return &Service{
		storage:        storage,
		cache:          cache,
		eventPublisher: eventPublisher,
	}
}

func (s *Service) FriendDelete(ctx context.Context, userID1, userID2 models.UserID) error {
	err := s.storage.FriendDelete(ctx, userID1, userID2)
	if err != nil {
		return err
	}

	s.cache.FriendDelete(userID1, userID2)
	err = s.eventPublisher.FriendDelete(userID1, userID2)
	if err != nil {
		logger.Error("Failed to notify about friend changes: %v", err)
	}
	return nil
}

func (s *Service) FriendAdd(ctx context.Context, userID1, userID2 models.UserID) error {
	err := s.storage.FriendAdd(ctx, userID1, userID2)
	if err != nil {
		return err
	}

	s.cache.FriendAdd(userID1, userID2)
	err = s.eventPublisher.FriendAdd(userID1, userID2)
	if err != nil {
		logger.Error("Failed to notify about friend changes: %v", err)
	}
	return nil
}
