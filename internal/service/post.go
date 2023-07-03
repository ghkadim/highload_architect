package service

import (
	"context"
	"errors"

	openapi "github.com/ghkadim/highload_architect/generated/go_server/go"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

// PostCreatePost -
func (s *ApiService) PostCreatePost(ctx context.Context, postCreatePostRequest openapi.PostCreatePostRequest) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := s.session.ParseToken(ctx, token)
	if err != nil {
		logger.Error("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	post, err := s.master.PostAdd(ctx, postCreatePostRequest.Text, userID)
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	s.cache.PostAdd(post)
	return openapi.Response(200, post.ID), nil
}

// PostDeleteIdPut -
func (s *ApiService) PostDeleteIdPut(ctx context.Context, id string) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := s.session.ParseToken(ctx, token)
	if err != nil {
		logger.Error("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	post, err := s.master.PostGet(ctx, models.PostID(id))
	if err != nil {
		if errors.Is(err, models.ErrPostNotFound) {
			logger.Error("Post already deleted: %v", err)
			return openapi.Response(200, nil), nil
		}
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	if post.AuthorID != userID {
		return openapi.Response(403, nil), errors.New("post deletion forbidden for user")
	}

	err = s.master.PostDelete(ctx, models.PostID(id))
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	s.cache.PostDelete(models.PostID(id))
	return openapi.Response(200, nil), nil
}

// PostFeedGet -
func (s *ApiService) PostFeedGet(ctx context.Context, offset, limit float32) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := s.session.ParseToken(ctx, token)
	if err != nil {
		logger.Error("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	cached := true
	posts, err := s.cache.PostFeed(userID, int(offset), int(limit))
	if err != nil {
		if !(errors.Is(err, models.ErrFeedNotFound) || errors.Is(err, models.ErrFeedPartial)) {
			return openapi.Response(500, openapi.LoginPost500Response{}), err
		}
		cached = false
	}

	if !cached {
		dbOffset := int(offset) + len(posts)
		dbLimit := int(limit) - len(posts)

		dbPosts, err := s.master.PostFeed(ctx, userID, dbOffset, dbLimit)
		if err != nil {
			return openapi.Response(500, openapi.LoginPost500Response{}), err
		}

		posts = append(posts, dbPosts...)
	}

	postsResp := make([]openapi.Post, 0, len(posts))
	for _, post := range posts {
		postsResp = append(postsResp, openapi.Post{
			Id:           string(post.ID),
			Text:         post.Text,
			AuthorUserId: string(post.AuthorID),
		})
	}

	return openapi.Response(200, postsResp), nil
}

// PostGetIdGet -
func (s *ApiService) PostGetIdGet(ctx context.Context, id string) (openapi.ImplResponse, error) {
	post, err := s.master.PostGet(ctx, models.PostID(id))
	if err != nil {
		if errors.Is(err, models.ErrPostNotFound) {
			logger.Error("Post already deleted: %v", err)
			return openapi.Response(404, nil), nil
		}
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	return openapi.Response(200, openapi.Post{Id: string(post.ID), Text: post.Text, AuthorUserId: string(post.AuthorID)}), nil
}

// PostUpdatePut -
func (s *ApiService) PostUpdatePut(ctx context.Context, postUpdatePutRequest openapi.PostUpdatePutRequest) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := s.session.ParseToken(ctx, token)
	if err != nil {
		logger.Error("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	post, err := s.master.PostGet(ctx, models.PostID(postUpdatePutRequest.Id))
	if err != nil {
		if errors.Is(err, models.ErrPostNotFound) {
			logger.Error("Post already deleted: %v", err)
			return openapi.Response(200, nil), nil
		}
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	if post.AuthorID != userID {
		return openapi.Response(403, nil), errors.New("post deletion forbidden for user")
	}

	err = s.master.PostUpdate(ctx, models.PostID(postUpdatePutRequest.Id), postUpdatePutRequest.Text)
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	return openapi.Response(200, nil), nil
}
