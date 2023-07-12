package service

import (
	"context"
	"sync/atomic"

	"github.com/ghkadim/highload_architect/internal/models"
	"github.com/ghkadim/highload_architect/internal/result"
	"github.com/ghkadim/highload_architect/internal/utils/closer"
)

type Storage interface {
	UserRegister(ctx context.Context, user models.User) (models.UserID, error)
	UserGet(ctx context.Context, id models.UserID) (models.User, error)
	UserSearch(ctx context.Context, firstName, secondName string) ([]models.User, error)
	PostAdd(ctx context.Context, text string, author models.UserID) (models.Post, error)
	PostUpdate(ctx context.Context, postID models.PostID, text string) error
	PostDelete(ctx context.Context, postID models.PostID) error
	PostGet(ctx context.Context, postID models.PostID) (models.Post, error)
	PostFeed(ctx context.Context, userID models.UserID, offset, limit int) ([]models.Post, error)
	FriendAdd(ctx context.Context, userID1, userID2 models.UserID) error
	FriendDelete(ctx context.Context, userID1, userID2 models.UserID) error
	UserFriends(ctx context.Context, user models.UserID) ([]models.UserID, error)
	DialogSend(ctx context.Context, message models.DialogMessage) error
	DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error)
}

type Cache interface {
	PostAdd(post models.Post)
	PostUpdate(postID models.PostID, text string)
	PostDelete(postID models.PostID)
	PostFeed(userID models.UserID, offset, limit int) ([]models.Post, error)
	FriendAdd(userID1, userID2 models.UserID)
	FriendDelete(userID1, userID2 models.UserID)
}

type Session interface {
	HashPassword(ctx context.Context, password string) ([]byte, error)
	CompareHashAndPassword(ctx context.Context, hash []byte, password string) (bool, error)
	TokenForUser(ctx context.Context, userID models.UserID) (string, error)
}

type EventPublisher interface {
	PostAdd(post models.Post) error
	FriendAdd(userID, friendID models.UserID) error
	FriendDelete(userID, friendID models.UserID) error
}

type EventConsumer interface {
	PostAdded(ctx context.Context, userID models.UserID, friends []models.UserID) (<-chan result.Result[models.Post], closer.Closer)
	FriendUpdated(ctx context.Context, userID models.UserID) (<-chan result.Result[models.FriendEvent], closer.Closer)
}

type service struct {
	master         Storage
	replicas       []Storage
	cache          Cache
	session        Session
	eventPublisher EventPublisher
	eventConsumer  EventConsumer

	replicaNum atomic.Int32
}

func NewService(
	master Storage,
	replicas []Storage,
	cache Cache,
	session Session,
	publisher EventPublisher,
	consumer EventConsumer,
) *service {
	return &service{
		master:         master,
		replicas:       replicas,
		cache:          cache,
		session:        session,
		eventPublisher: publisher,
		eventConsumer:  consumer,
	}
}

func (s *service) readStorage() Storage {
	if len(s.replicas) != 0 {
		return s.replicas[int(s.replicaNum.Add(1))%len(s.replicas)]
	}
	return s.master
}
