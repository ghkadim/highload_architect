package service

import (
	"context"
	"github.com/ghkadim/highload_architect/internal/models"
	"sync/atomic"
)

type Storage interface {
	UserRegister(ctx context.Context, user models.User) (string, error)
	UserGet(ctx context.Context, id string) (models.User, error)
	UserSearch(ctx context.Context, firstName, secondName string) ([]models.User, error)
	FriendAdd(ctx context.Context, userID1, userID2 string) error
	FriendDelete(ctx context.Context, userID1, userID2 string) error
}

type Session interface {
	HashPassword(ctx context.Context, password string) ([]byte, error)
	CompareHashAndPassword(ctx context.Context, hash []byte, password string) (bool, error)
	TokenForUser(ctx context.Context, userID string) (string, error)
	ParseToken(ctx context.Context, tokenStr string) (string, error)
}

type ApiService struct {
	master     Storage
	replicas   []Storage
	replicaNum atomic.Int32
	session    Session
}

// NewApiService creates an api service
func NewApiService(
	master Storage,
	replicas []Storage,
	session Session,
) *ApiService {
	return &ApiService{
		master:   master,
		replicas: replicas,
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
