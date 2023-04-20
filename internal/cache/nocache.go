package cache

import "github.com/ghkadim/highload_architect/internal/models"

type DisabledCache struct{}

func NewDisabledCache() *DisabledCache {
	return &DisabledCache{}
}

func (c *DisabledCache) FriendAdd(userID1, userID2 models.UserID) {}

func (c *DisabledCache) FriendDelete(userID1, userID2 models.UserID) {}

func (c *DisabledCache) PostAdd(post models.Post) {}

func (c *DisabledCache) PostUpdate(postID models.PostID, text string) {}

func (c *DisabledCache) PostDelete(postID models.PostID) {}

func (c *DisabledCache) PostFeed(userID models.UserID, offset, limit int) ([]models.Post, error) {
	return nil, models.ErrFeedNotFound
}
