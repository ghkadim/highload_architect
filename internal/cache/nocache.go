package cache

import "github.com/ghkadim/highload_architect/internal/models"

type disabledCache struct{}

func NewDisabledCache() Cache {
	return &disabledCache{}
}

func (c *disabledCache) FriendAdd(userID1, userID2 models.UserID) {}

func (c *disabledCache) FriendDelete(userID1, userID2 models.UserID) {}

func (c *disabledCache) PostAdd(post models.Post) {}

func (c *disabledCache) PostUpdate(postID models.PostID, text string) {}

func (c *disabledCache) PostDelete(postID models.PostID) {}

func (c *disabledCache) PostFeed(userID models.UserID, offset, limit int) ([]models.Post, error) {
	return nil, models.ErrFeedNotFound
}
