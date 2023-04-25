package cache

import (
	"sort"

	"github.com/ghkadim/highload_architect/internal/models"
)

type cachedPost struct {
	*models.Post
	feeds map[models.UserID]struct{}
}

type feedsCache struct {
	posts     map[models.PostID]cachedPost
	feeds     map[models.UserID][]*models.Post
	feedLimit int
}

func newFeedsCache(feedLimit int) *feedsCache {
	return &feedsCache{
		posts:     make(map[models.PostID]cachedPost),
		feeds:     make(map[models.UserID][]*models.Post),
		feedLimit: feedLimit,
	}
}

func (c *feedsCache) PostAdd(post models.Post, feed models.UserID) {
	posts, ok := c.feeds[feed]
	if !ok {
		return
	}

	p, ok := c.posts[post.ID]
	if !ok {
		p = cachedPost{
			Post: &post,
			feeds: map[models.UserID]struct{}{
				feed: {},
			},
		}
		c.posts[post.ID] = p
	} else {
		if _, ok := p.feeds[feed]; ok {
			return
		}
	}
	p.feeds[feed] = struct{}{}

	posts = append(posts, &post)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].SequentialID >= posts[j].SequentialID
	})
	c.feeds[feed] = posts
	if len(posts) > c.feedLimit {
		c.pop(posts[c.feedLimit].ID, feed)
	}
}

func (c *feedsCache) pop(id models.PostID, feed models.UserID) {
	posts := c.feeds[feed]
	if len(posts) == 0 {
		return
	}

	postsCopy := make([]*models.Post, 0, len(posts))
	for _, post := range posts {
		if post.ID == id {
			p := c.posts[post.ID]
			delete(p.feeds, feed)
			if len(p.feeds) == 0 {
				delete(c.posts, post.ID)
			}
			continue
		}
		postsCopy = append(postsCopy, post)
	}

	c.feeds[feed] = postsCopy
}

func (c *feedsCache) PopForAuthor(feed models.UserID, author models.UserID) {
	posts := c.feeds[feed]
	if len(posts) == 0 {
		return
	}

	postsCopy := make([]*models.Post, 0, len(posts))
	for _, post := range posts {
		if post.AuthorID == author {
			p := c.posts[post.ID]
			delete(p.feeds, feed)
			if len(p.feeds) == 0 {
				delete(c.posts, post.ID)
			}
			continue
		}
		postsCopy = append(postsCopy, post)
	}

	if len(postsCopy) == 0 {
		delete(c.feeds, feed)
	} else {
		c.feeds[feed] = postsCopy
	}
}

func (c *feedsCache) PostUpdate(id models.PostID, text string) {
	p, ok := c.posts[id]
	if !ok {
		return
	}
	p.Text = text
}

func (c *feedsCache) Cached(feed models.UserID) bool {
	_, ok := c.feeds[feed]
	return ok
}

func (c *feedsCache) FeedAdd(feed models.UserID) {
	c.FeedDelete(feed)
	c.feeds[feed] = make([]*models.Post, 0)
}

func (c *feedsCache) FeedDelete(feed models.UserID) {
	posts := c.feeds[feed]
	delete(c.feeds, feed)
	if len(posts) == 0 {
		return
	}

	for _, post := range posts {
		p := c.posts[post.ID]
		delete(p.feeds, feed)
		if len(p.feeds) == 0 {
			delete(c.posts, post.ID)
		}
	}
}

func (c *feedsCache) Feed(feed models.UserID) ([]models.Post, bool) {
	posts, ok := c.feeds[feed]
	res := make([]models.Post, 0, len(posts))
	for i := range posts {
		res = append(res, *posts[i])
	}
	return res, ok
}

func (c *feedsCache) FeedsForPost(id models.PostID) []models.UserID {
	res := make([]models.UserID, 0)
	for userID := range c.posts[id].feeds {
		res = append(res, userID)
	}
	return res
}
