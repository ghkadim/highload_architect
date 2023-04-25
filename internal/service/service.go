package service

import (
	"context"
	"sync/atomic"

	"github.com/ghkadim/highload_architect/internal/models"
)

type Storage interface {
	UserRegister(ctx context.Context, user models.User) (models.UserID, error)
	UserGet(ctx context.Context, id models.UserID) (models.User, error)
	UserSearch(ctx context.Context, firstName, secondName string) ([]models.User, error)
	FriendAdd(ctx context.Context, userID1, userID2 models.UserID) error
	FriendDelete(ctx context.Context, userID1, userID2 models.UserID) error
	PostAdd(ctx context.Context, text string, author models.UserID) (models.Post, error)
	PostUpdate(ctx context.Context, postID models.PostID, text string) error
	PostDelete(ctx context.Context, postID models.PostID) error
	PostGet(ctx context.Context, postID models.PostID) (models.Post, error)
	PostFeed(ctx context.Context, userID models.UserID, offset, limit int) ([]models.Post, error)
}

type Cache interface {
	FriendAdd(userID1, userID2 models.UserID)
	FriendDelete(userID1, userID2 models.UserID)
	PostAdd(post models.Post)
	PostUpdate(postID models.PostID, text string)
	PostDelete(postID models.PostID)
	PostFeed(userID models.UserID, offset, limit int) ([]models.Post, error)
}

type Session interface {
	HashPassword(ctx context.Context, password string) ([]byte, error)
	CompareHashAndPassword(ctx context.Context, hash []byte, password string) (bool, error)
	TokenForUser(ctx context.Context, userID models.UserID) (string, error)
	ParseToken(ctx context.Context, tokenStr string) (models.UserID, error)
}

type ApiService struct {
	master     Storage
	replicas   []Storage
	replicaNum atomic.Int32
	session    Session
	cache      Cache
}

// NewApiService creates an api service
func NewApiService(
	master Storage,
	replicas []Storage,
	cache Cache,
	session Session,
) *ApiService {
	return &ApiService{
		master:   master,
		replicas: replicas,
		cache:    cache,
		session:  session,
	}
}

func valueOrDefault[V any](value *V) V {
	if value == nil {
		return *new(V)
	}
	return *value
}

func bearerToken(ctx context.Context) string {
	val := ctx.Value("BearerToken")
	if val == nil {
		return ""
	}
	return val.(string)
}

func (s *ApiService) readStorage() Storage {
	if len(s.replicas) != 0 {
		return s.replicas[int(s.replicaNum.Add(1))%len(s.replicas)]
	}
	return s.master
}
