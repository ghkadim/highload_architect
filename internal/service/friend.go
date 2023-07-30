package service

import (
	"context"

	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

func (s *service) FriendDelete(ctx context.Context, userID1, userID2 models.UserID) error {
	err := s.master.FriendDelete(ctx, userID1, userID2)
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

func (s *service) FriendAdd(ctx context.Context, userID1, userID2 models.UserID) error {
	err := s.master.FriendAdd(ctx, userID1, userID2)
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
