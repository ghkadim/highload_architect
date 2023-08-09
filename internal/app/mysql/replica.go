package mysql

import (
	"context"
	"sync/atomic"

	"github.com/ghkadim/highload_architect/internal/models"
)

type StorageWithReplicas struct {
	Storage
	replicas []Storage

	replicaNum atomic.Int32
}

var _ Storage = &StorageWithReplicas{}

func NewStorageWithReplicas(master Storage, replicas []Storage) Storage {
	return &StorageWithReplicas{Storage: master, replicas: replicas}
}

func (s *StorageWithReplicas) readStorage() Storage {
	if len(s.replicas) != 0 {
		return s.replicas[int(s.replicaNum.Add(1))%len(s.replicas)]
	}
	return s.Storage
}

func (s *StorageWithReplicas) UserGet(ctx context.Context, id models.UserID) (models.User, error) {
	return s.readStorage().UserGet(ctx, id)
}

func (s *StorageWithReplicas) UserSearch(ctx context.Context, firstName, secondName string) ([]models.User, error) {
	return s.readStorage().UserSearch(ctx, firstName, secondName)
}

func (s *StorageWithReplicas) PostGet(ctx context.Context, postID models.PostID) (models.Post, error) {
	return s.readStorage().PostGet(ctx, postID)
}

func (s *StorageWithReplicas) PostFeed(ctx context.Context, userID models.UserID, offset, limit int) ([]models.Post, error) {
	return s.readStorage().PostFeed(ctx, userID, offset, limit)
}

func (s *StorageWithReplicas) UserPosts(ctx context.Context, user models.UserID, limit int) ([]models.Post, error) {
	return s.readStorage().UserPosts(ctx, user, limit)
}

func (s *StorageWithReplicas) UserFriends(ctx context.Context, user models.UserID) ([]models.UserID, error) {
	return s.readStorage().UserFriends(ctx, user)
}

func (s *StorageWithReplicas) DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error) {
	return s.readStorage().DialogList(ctx, userID1, userID2)
}

func (s *StorageWithReplicas) DialogMatchingShard(ctx context.Context, matchExpr string, fromID models.DialogMessageID, limit int64) ([]models.DialogMessage, error) {
	return s.readStorage().DialogMatchingShard(ctx, matchExpr, fromID, limit)
}
