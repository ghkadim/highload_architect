package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync/atomic"
	"time"

	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
	"github.com/ghkadim/highload_architect/internal/result"
	"github.com/ghkadim/highload_architect/internal/utils/closer"
)

func (s *service) PostAdd(ctx context.Context, text string, author models.UserID) (models.PostID, error) {
	post, err := s.master.PostAdd(ctx, text, author)
	if err != nil {
		return "", err
	}

	s.cache.PostAdd(post)
	err = s.eventPublisher.PostAdd(post)
	if err != nil {
		logger.Error("Failed to send event for post: %v", err)
	}
	return post.ID, nil
}

func (s *service) PostDelete(ctx context.Context, userID models.UserID, postID models.PostID) error {
	post, err := s.master.PostGet(ctx, postID)
	if err != nil {
		if errors.Is(err, models.ErrPostNotFound) {
			logger.Error("Post already deleted: %v", err)
			return nil
		}
	}

	if post.AuthorID != userID {
		return errors.Join(models.ErrUnauthorized, errors.New("post deletion forbidden for user"))
	}

	err = s.master.PostDelete(ctx, postID)
	if err != nil {
		return err
	}

	s.cache.PostDelete(postID)
	return nil
}

func (s *service) PostFeed(ctx context.Context, userID models.UserID, offset, limit int) ([]models.Post, error) {
	cached := true
	posts, err := s.cache.PostFeed(userID, offset, limit)
	if err != nil {
		if !(errors.Is(err, models.ErrFeedNotFound) || errors.Is(err, models.ErrFeedPartial)) {
			return nil, err
		}
		cached = false
	}

	if !cached {
		dbOffset := offset + len(posts)
		dbLimit := limit - len(posts)

		dbPosts, err := s.readStorage().PostFeed(ctx, userID, dbOffset, dbLimit)
		if err != nil {
			return nil, err
		}

		posts = append(posts, dbPosts...)
	}
	return posts, nil
}

func (s *service) PostGet(ctx context.Context, postID models.PostID) (models.Post, error) {
	post, err := s.readStorage().PostGet(ctx, postID)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (s *service) PostUpdate(ctx context.Context, userID models.UserID, postID models.PostID, text string) error {
	post, err := s.master.PostGet(ctx, postID)
	if err != nil {
		return err
	}

	if post.AuthorID != userID {
		return errors.Join(models.ErrUnauthorized, errors.New("post update forbidden for user"))
	}

	err = s.master.PostUpdate(ctx, postID, text)
	if err != nil {
		return err
	}
	s.cache.PostUpdate(postID, text)
	return nil
}

func (s *service) PostFeedPosted(ctx context.Context, subscriber models.UserID) <-chan result.Result[models.Post] {
	resultCh := make(chan result.Result[models.Post])
	errorCh := make(chan error, 2)
	syncFriendsCh := make(chan []models.UserID)
	syncFriendsRequired := atomic.Bool{}
	syncFriendsRequired.Store(true)
	updatedFriendsCh := make(chan []models.UserID)
	ctx, ctxCloser := context.WithCancel(ctx)

	go func() {
		friendUpdatesCh, friendUpdatesCloser := s.eventConsumer.FriendUpdated(ctx, subscriber)
		defer friendUpdatesCloser.Close()
		defer close(updatedFriendsCh)
		var friends []models.UserID
		for {
			select {
			case synced, ok := <-syncFriendsCh:
				if !ok {
					return
				}
				sort.Slice(synced, func(i, j int) bool {
					return synced[i] < synced[j]
				})
				if !reflect.DeepEqual(friends, synced) {
					friends = synced
					updatedFriendsCh <- friends
				}
			case res, ok := <-friendUpdatesCh:
				if !ok {
					return
				}
				update, err := res.Value()
				if err != nil {
					errorCh <- fmt.Errorf("failed to listen friend updates: %w", err)
					return
				}

				logger.Debug("Got FriendUpdate event for user=%s", subscriber)
				switch update.Type {
				case models.FriendAddedEvent:
					syncFriendsRequired.Store(true)
				case models.FriendDeletedEvent:
					syncFriendsRequired.Store(true)
				default:
					logger.Error("Unknown FriendEventType %v", update.Type)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(1)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if syncFriendsRequired.Swap(false) {
					logger.Debug("Syncing friends for user=%s", subscriber)
					friends, err := s.readStorage().UserFriends(ctx, subscriber)
					if err != nil {
						errorCh <- fmt.Errorf("failed to listen friend updates: %w", err)
						close(syncFriendsCh)
						return
					}
					syncFriendsCh <- friends
				}
				ticker.Reset(time.Second)
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer ctxCloser()
		defer close(resultCh)

		var postEventCh <-chan result.Result[models.Post]
		clsr := closer.Nop()
		defer func() {
			_ = clsr.Close()
		}()
		for {
			select {
			case friends, ok := <-updatedFriendsCh:
				if !ok {
					return
				}
				_ = clsr.Close()
				postEventCh, clsr = s.eventConsumer.PostAdded(ctx, subscriber, friends)
			case res, ok := <-postEventCh:
				if !ok {
					return
				}
				_, err := res.Value()
				if err != nil {
					resultCh <- result.ErrorWrap[models.Post](err, "failed to listen for posts")
					return
				}
				resultCh <- res
			case err, ok := <-errorCh:
				if !ok {
					return
				}
				resultCh <- result.Error[models.Post](err)
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	return resultCh
}
