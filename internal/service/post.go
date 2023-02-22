package service

import (
	"context"
	"errors"
	openapi "github.com/ghkadim/highload_architect/generated/go_server/go"
	"github.com/ghkadim/highload_architect/internal/models"
	"log"
)

// PostCreatePost -
func (s *ApiService) PostCreatePost(ctx context.Context, postCreatePostRequest openapi.PostCreatePostRequest) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := s.session.ParseToken(ctx, token)
	if err != nil {
		log.Printf("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	postID, err := s.master.PostAdd(ctx, postCreatePostRequest.Text, userID)
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}

	return openapi.Response(200, postID), nil
}

// PostDeleteIdPut -
func (s *ApiService) PostDeleteIdPut(ctx context.Context, id string) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := s.session.ParseToken(ctx, token)
	if err != nil {
		log.Printf("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	post, err := s.master.PostGet(ctx, models.PostID(id))
	if err != nil {
		if errors.Is(err, models.PostNotFound) {
			log.Printf("Post already deleted: %v", err)
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

	return openapi.Response(200, nil), nil
}

// PostFeedGet -
func (s *ApiService) PostFeedGet(ctx context.Context, offset int32, limit int32) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := s.session.ParseToken(ctx, token)
	if err != nil {
		log.Printf("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	posts, err := s.readStorage().PostFeed(ctx, userID, int(offset), int(limit))
	if err != nil {
		return openapi.Response(500, openapi.LoginPost500Response{}), err
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
		if errors.Is(err, models.PostNotFound) {
			log.Printf("Post already deleted: %v", err)
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
		log.Printf("Bad token: %v", err)
		return openapi.Response(401, nil), nil
	}

	post, err := s.master.PostGet(ctx, models.PostID(postUpdatePutRequest.Id))
	if err != nil {
		if errors.Is(err, models.PostNotFound) {
			log.Printf("Post already deleted: %v", err)
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
