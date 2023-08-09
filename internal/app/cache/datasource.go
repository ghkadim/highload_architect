package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

//go:generate mockery --name DataSource
type DataSource interface {
	UserPosts(ctx context.Context, user models.UserID, limit int) ([]models.Post, error)
	PostFeed(ctx context.Context, user models.UserID, offset, limit int) ([]models.Post, error)
	UserFriends(ctx context.Context, user models.UserID) ([]models.UserID, error)
}

type loadWithRetry struct {
	dataSource DataSource
	tasks      chan struct{}
}

func NewLoadWithRetry(dataSource DataSource) *loadWithRetry {
	tasks := make(chan struct{}, 1)
	tasks <- struct{}{}

	return &loadWithRetry{
		dataSource: dataSource,
		tasks:      tasks,
	}
}

func loadData[T any](s *loadWithRetry, f func() (T, error)) (T, error) {
	<-s.tasks
	defer func() {
		s.tasks <- struct{}{}
	}()
	for {
		res, err := f()
		if err != nil {
			logger.Error("Cache failed to fetch data: %v", err)
			time.Sleep(time.Minute)
			continue
		}
		return res, nil
	}
}

func (s *loadWithRetry) UserPosts(ctx context.Context, user models.UserID, limit int) ([]models.Post, error) {
	return loadData(s, func() ([]models.Post, error) {
		posts, err := s.dataSource.UserPosts(ctx, user, limit)
		if err != nil {
			return nil, fmt.Errorf("UserPosts: %w", err)
		}
		return posts, nil
	})
}

func (s *loadWithRetry) PostFeed(ctx context.Context, user models.UserID, offset, limit int) ([]models.Post, error) {
	return loadData(s, func() ([]models.Post, error) {
		posts, err := s.dataSource.PostFeed(ctx, user, offset, limit)
		if err != nil {
			return nil, fmt.Errorf("PostFeed: %w", err)
		}
		return posts, nil
	})
}

func (s *loadWithRetry) UserFriends(ctx context.Context, user models.UserID) ([]models.UserID, error) {
	return loadData(s, func() ([]models.UserID, error) {
		friends, err := s.dataSource.UserFriends(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("UserFriends: %w", err)
		}
		return friends, nil
	})
}
