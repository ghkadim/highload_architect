package cache

import (
	"context"
	"sync"

	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

type Cache interface {
	FriendAdd(userID1, userID2 models.UserID)
	FriendDelete(userID1, userID2 models.UserID)
	PostAdd(post models.Post)
	PostUpdate(postID models.PostID, text string)
	PostDelete(postID models.PostID)
	PostFeed(userID models.UserID, offset, limit int) ([]models.Post, error)
}

type cache struct {
	sync.RWMutex
	dataSource    DataSource
	feeds         *feedsCache
	subscribers   map[models.UserID]map[models.UserID]struct{}
	feedLimitSize int
}

func NewCache(feedLimit int, dataSource DataSource) Cache {
	c := &cache{
		feeds:         newFeedsCache(feedLimit),
		subscribers:   map[models.UserID]map[models.UserID]struct{}{},
		feedLimitSize: feedLimit,
		dataSource:    dataSource,
	}
	return c
}

func (c *cache) FriendAdd(userID1, userID2 models.UserID) {
	c.Lock()
	defer c.Unlock()

	if !c.feeds.Cached(userID1) {
		return
	}

	if _, ok := c.subscribers[userID2]; !ok {
		c.subscribers[userID2] = make(map[models.UserID]struct{})
	}
	c.subscribers[userID2][userID1] = struct{}{}

	go func() {
		posts, err := c.dataSource.UserPosts(context.Background(), userID2, c.feedLimitSize)
		if err != nil {
			logger.Error("FriendAdd failed to load UserPosts: %v", err)
			return
		}
		for _, p := range posts {
			c.feedPostAdd(userID1, p)
		}
		logger.Info("UsersPosts for %s from %s loaded to cache (%d posts)", userID1, userID2, len(posts))
	}()
}

func (c *cache) FriendDelete(userID1, userID2 models.UserID) {
	c.Lock()
	defer c.Unlock()

	delete(c.subscribers[userID2], userID1)
	if !c.feeds.Cached(userID1) {
		return
	}
	c.feeds.PopForAuthor(userID1, userID2)

	go func() {
		posts, err := c.dataSource.PostFeed(context.Background(), userID1, 0, c.feedLimitSize)
		if err != nil {
			logger.Error("FriendDelete failed to load UserPosts: %v", err)
			return
		}
		c.feedInit(userID1)
		for _, p := range posts {
			c.feedPostAdd(userID1, p)
		}
		logger.Info("Feed for %s loaded to cache (%d posts)", userID1, len(posts))
	}()
}

func (c *cache) PostAdd(post models.Post) {
	c.Lock()
	defer c.Unlock()

	for userID := range c.subscribers[post.AuthorID] {
		c.feedPostAdd(userID, post)
	}
}

func (c *cache) PostUpdate(postID models.PostID, text string) {
	c.Lock()
	defer c.Unlock()
	c.feeds.PostUpdate(postID, text)
}

func (c *cache) PostDelete(postID models.PostID) {
	c.Lock()
	defer c.Unlock()

	for _, userID := range c.feeds.FeedsForPost(postID) {
		go func(userID models.UserID) {
			c.Lock()
			defer c.Unlock()
			c.feeds.pop(postID, userID)
		}(userID)
	}
}

func (c *cache) PostFeed(userID models.UserID, offset, limit int) ([]models.Post, error) {
	c.RLock()
	defer c.RUnlock()
	posts, ok := c.feeds.Feed(userID)
	if !ok {
		go func() {
			friends, err := c.dataSource.UserFriends(context.Background(), userID)
			if err != nil {
				logger.Error("PostFeed failed to load UserFriends: %v", err)
				return
			}
			posts, err := c.dataSource.PostFeed(context.Background(), userID, 0, c.feedLimitSize)
			if err != nil {
				logger.Error("PostFeed failed to load PostFeed: %v", err)
				return
			}
			c.feedInit(userID)
			c.friendsAdd(userID, friends)
			logger.Info("Friends for %s loaded to cache (%d friends)", userID, len(friends))
			for _, p := range posts {
				c.feedPostAdd(userID, p)
			}
			logger.Info("Feed for %s loaded to cache (%d posts)", userID, len(posts))
		}()
		return nil, models.ErrFeedNotFound
	}

	var partialResponseErr error
	if (len(posts) == c.feedLimitSize) && (c.feedLimitSize < offset+limit) {
		partialResponseErr = models.ErrFeedPartial
	}

	if len(posts) <= offset {
		return nil, partialResponseErr
	}
	if len(posts) <= offset+limit {
		return posts[offset:], partialResponseErr
	}
	return posts[offset : offset+limit], nil
}

func (c *cache) feedLimit() int {
	return c.feedLimitSize
}

func (c *cache) feedPostAdd(userID models.UserID, post models.Post) {
	go func() {
		c.Lock()
		defer c.Unlock()
		c.feeds.PostAdd(post, userID)
	}()
}

func (c *cache) feedInit(userID models.UserID) {
	c.Lock()
	defer c.Unlock()
	c.feeds.FeedAdd(userID)
}

func (c *cache) friendsAdd(userID models.UserID, friendIDs []models.UserID) {
	c.Lock()
	defer c.Unlock()
	for _, friend := range friendIDs {
		_, ok := c.subscribers[friend]
		if !ok {
			c.subscribers[friend] = make(map[models.UserID]struct{})
		}
		c.subscribers[friend][userID] = struct{}{}
	}
}
