package service

import (
	"context"
	"errors"

	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

func (s *service) PostAdd(ctx context.Context, text string, author models.UserID) (models.PostID, error) {
	post, err := s.master.PostAdd(ctx, text, author)
	if err != nil {
		return models.PostID(""), err
	}

	s.cache.PostAdd(post)
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
