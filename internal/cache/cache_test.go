package cache

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ghkadim/highload_architect/internal/cache/mocks"
	"github.com/ghkadim/highload_architect/internal/models"
)

func TestCache_EmptyFeed(t *testing.T) {
	ds := mocks.NewDataSource(t)
	c := NewCache(1, ds)
	u := newUsersAndPosts().makeUser()

	ds.EXPECT().PostFeed(mock.Anything, u.ID, 0, 1).Return(nil, nil)
	ds.EXPECT().UserFriends(mock.Anything, u.ID).Return(nil, nil)
	feed, err := c.PostFeed(u.ID, 0, 10)
	assert.Empty(t, feed)
	assert.ErrorIs(t, err, models.ErrFeedNotFound)

	assert.Eventually(t,
		func() bool {
			feed, err := c.PostFeed(u.ID, 0, 10)
			if err != nil {
				require.ErrorIs(t, err, models.ErrFeedNotFound)
				return false
			}

			require.Nil(t, err)
			require.ElementsMatch(t, nil, feed)
			return true
		},
		time.Second,
		time.Millisecond*10)
}

func TestCache_FeedCreate(t *testing.T) {
	up := newUsersAndPosts()
	users := up.Users
	u := up.makeUser()

	tests := []struct {
		name    string
		friends friends
		limit   int
		offset  int
	}{
		{
			"empty",
			friends{},
			10,
			0,
		},
		{
			"one friend",
			friends{users[0].ID},
			10,
			0,
		},
		{
			"two friends",
			friends{users[0].ID, users[1].ID},
			10,
			0,
		},
		{
			"out of limit",
			friends{users[0].ID, users[1].ID, users[2].ID},
			10,
			0,
		},
		{
			"with offset",
			friends{users[0].ID, users[1].ID, users[2].ID},
			10,
			5,
		},
		{
			"with offset reduce limit",
			friends{users[0].ID, users[1].ID, users[2].ID},
			10,
			10,
		},
		{
			"with offset out of range",
			friends{users[0].ID, users[1].ID, users[2].ID},
			10,
			15,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewCache(1000, nil).(*cache)
			fullFeed := up.makeFeed(test.friends...)

			testFeedCreated(t, c, u, test.friends, fullFeed, test.offset, test.limit)
		})
	}
}

func TestCache_PostAdd(t *testing.T) {
	up := newUsersAndPosts()
	users := up.Users
	u := up.makeUser()
	limit := 10
	offset := 0

	tests := []struct {
		name       string
		friends    friends
		postAuthor models.UserID
	}{
		{
			name:       "not a friend",
			friends:    friends{users[0].ID, users[1].ID},
			postAuthor: users[2].ID,
		},
		{
			name:       "from friend",
			friends:    friends{users[0].ID},
			postAuthor: users[0].ID,
		},
		{
			name:       "from friend move feed",
			friends:    friends{users[0].ID, users[1].ID},
			postAuthor: users[0].ID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			up := up.copy()
			c := NewCache(1000, nil).(*cache)
			fullFeed := up.makeFeed(test.friends...)

			testFeedCreated(t, c, u, test.friends, fullFeed, offset, limit)

			newPost := up.makePost(test.postAuthor)
			fullFeed = up.makeFeed(test.friends...)
			expectedPosts := limitOffset(fullFeed, limit, offset)

			c.PostAdd(newPost)
			testFeedUpdated(t, func() ([]models.Post, error) {
				return c.PostFeed(u.ID, offset, limit)
			}, expectedPosts)
		})
	}
}

func TestCache_PostUpdate(t *testing.T) {
	up := newUsersAndPosts()
	users := up.Users
	u := up.makeUser()
	limit := 10
	offset := 0

	tests := []struct {
		name    string
		friends friends
		postID  models.PostID
	}{
		{
			name:    "not a friend",
			friends: friends{users[0].ID, users[1].ID},
			postID:  up.makeFeed(users[2].ID)[0].ID,
		},
		{
			name:    "from friend",
			friends: friends{users[0].ID, users[1].ID},
			postID:  up.makeFeed(users[0].ID)[0].ID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			up := up.copy()
			c := NewCache(1000, nil).(*cache)
			fullFeed := up.makeFeed(test.friends...)
			require.Len(t, fullFeed, 10)

			testFeedCreated(t, c, u, test.friends, fullFeed, offset, limit)

			for i := range up.Posts {
				if up.Posts[i].ID == test.postID {
					up.Posts[i].Text = "updated"
				}
			}
			fullFeed = up.makeFeed(test.friends...)
			expectedPosts := limitOffset(fullFeed, limit, offset)

			c.PostUpdate(test.postID, "updated")
			testFeedUpdated(t, func() ([]models.Post, error) {
				return c.PostFeed(u.ID, offset, limit)
			}, expectedPosts)
		})
	}
}

func TestCache_PostDelete(t *testing.T) {
	up := newUsersAndPosts()
	users := up.Users
	u := up.makeUser()
	limit := 10
	offset := 0

	tests := []struct {
		name    string
		friends friends
		postID  models.PostID
	}{
		{
			name:    "not a friend",
			friends: friends{users[0].ID, users[1].ID},
			postID:  up.makeFeed(users[2].ID)[0].ID,
		},
		{
			name:    "from friend",
			friends: friends{users[0].ID, users[1].ID},
			postID:  up.makeFeed(users[0].ID)[0].ID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			up := up.copy()
			c := NewCache(1000, nil).(*cache)
			fullFeed := up.makeFeed(test.friends...)
			require.Len(t, fullFeed, 10)

			testFeedCreated(t, c, u, test.friends, fullFeed, offset, limit)

			newPosts := make([]models.Post, 0)
			for _, post := range up.Posts {
				if post.ID != test.postID {
					newPosts = append(newPosts, post)
				}
			}
			up.Posts = newPosts

			fullFeed = up.makeFeed(test.friends...)
			expectedPosts := limitOffset(fullFeed, limit, offset)
			c.PostDelete(test.postID)
			testFeedUpdated(t, func() ([]models.Post, error) {
				return c.PostFeed(u.ID, offset, limit)
			}, expectedPosts)
		})
	}
}

func TestCache_FriendDelete(t *testing.T) {
	up := newUsersAndPosts()
	users := up.Users
	u := up.makeUser()
	postless := up.makeUser()
	limit := 10
	offset := 0

	tests := []struct {
		name    string
		friends friends
		delete  models.UserID
	}{
		{
			name:    "not a friend",
			friends: friends{users[0].ID, users[1].ID},
			delete:  users[2].ID,
		},
		{
			name:    "from friend",
			friends: friends{users[0].ID, users[1].ID},
			delete:  users[1].ID,
		},
		{
			name:    "from friend without posts",
			friends: friends{users[0].ID, users[1].ID, postless.ID},
			delete:  postless.ID,
		},
		{
			name:    "delete all",
			friends: friends{users[1].ID},
			delete:  users[1].ID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			up := up.copy()
			ds := mocks.NewDataSource(t)
			c := NewCache(1000, ds).(*cache)
			fullFeed := up.makeFeed(test.friends...)

			testFeedCreated(t, c, u, test.friends, fullFeed, offset, limit)

			newFriends := friends{}
			for _, friend := range test.friends {
				if friend == test.delete {
					continue
				}
				newFriends = append(newFriends, friend)
			}
			fullFeed = up.makeFeed(newFriends...)
			expectedPosts := limitOffset(fullFeed, limit, offset)

			ds.EXPECT().PostFeed(mock.Anything, u.ID, 0, c.feedLimit()).Return(fullFeed, nil)

			c.FriendDelete(u.ID, test.delete)
			testFeedUpdated(t, func() ([]models.Post, error) {
				return c.PostFeed(u.ID, offset, limit)
			}, expectedPosts)
		})
	}
}

func TestCache_FriendAdd(t *testing.T) {
	up := newUsersAndPosts()
	users := up.Users
	u := up.makeUser()
	postless := up.makeUser()
	limit := 10
	offset := 0

	tests := []struct {
		name    string
		friends friends
		add     models.UserID
	}{
		{
			name:    "first friend",
			friends: friends{},
			add:     users[0].ID,
		},
		{
			name:    "second friend",
			friends: friends{users[0].ID},
			add:     users[1].ID,
		},
		{
			name:    "repeated add",
			friends: friends{users[0].ID},
			add:     users[0].ID,
		},
		{
			name:    "add postless first",
			friends: friends{},
			add:     postless.ID,
		},
		{
			name:    "add postless second",
			friends: friends{users[1].ID},
			add:     postless.ID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			up := up.copy()
			ds := mocks.NewDataSource(t)
			c := NewCache(1000, ds).(*cache)
			fullFeed := up.makeFeed(test.friends...)

			testFeedCreated(t, c, u, test.friends, fullFeed, offset, limit)

			fullFeed = up.makeFeed(append(test.friends, test.add)...)
			expectedPosts := limitOffset(fullFeed, limit, offset)

			ds.EXPECT().UserPosts(mock.Anything, test.add, c.feedLimit()).
				Return(up.makeFeed(test.add), nil).Maybe()

			c.FriendAdd(u.ID, test.add)
			testFeedUpdated(t, func() ([]models.Post, error) {
				return c.PostFeed(u.ID, offset, limit)
			}, expectedPosts)
		})
	}
}

func testFeedCreated(
	t *testing.T,
	c *cache,
	u models.User,
	friends friends,
	fullFeed []models.Post,
	offset int,
	limit int,
) {
	ds, ok := c.dataSource.(*mocks.DataSource)
	if !ok {
		ds = mocks.NewDataSource(t)
		c.dataSource = ds
	}

	ufCall := ds.EXPECT().
		UserFriends(mock.Anything, u.ID).
		Return(friends, nil)
	defer ufCall.Unset()

	pfCall := ds.EXPECT().
		PostFeed(mock.Anything, u.ID, 0, c.feedLimit()).
		Return(fullFeed, nil)
	defer pfCall.Unset()

	feed, err := c.PostFeed(u.ID, offset, limit)
	require.Empty(t, feed)
	require.ErrorIs(t, err, models.ErrFeedNotFound)

	expectedPosts := limitOffset(fullFeed, limit, offset)
	testFeedUpdated(t, func() ([]models.Post, error) {
		return c.PostFeed(u.ID, offset, limit)
	}, expectedPosts)

	ds.AssertExpectations(t)
}

func testFeedUpdated(
	t *testing.T,
	postFeedF func() ([]models.Post, error),
	expectedPosts []models.Post,
) {
	lm := &lastMsg{}
	ok := assert.Eventually(t,
		func() bool {
			feed, err := postFeedF()
			if err != nil {
				require.ErrorIs(t, err, models.ErrFeedNotFound)
				return false
			}
			if !assert.ElementsMatch(lm, expectedPosts, feed) {
				return false
			}

			require.Nil(t, err)
			require.ElementsMatch(t, expectedPosts, feed)
			return true
		},
		time.Second,
		time.Millisecond*10)

	if !ok {
		t.Errorf("Last error message: %s", lm.Msg)
	}
}

type friends []models.UserID

type usersAndPosts struct {
	Users []models.User
	Posts []models.Post

	postIndex int64
	userIndex int64
}

func newUsersAndPosts() *usersAndPosts {
	up := &usersAndPosts{}
	for i := 0; i < 5; i++ {
		up.makeUser()
	}
	for i := 0; i < 5; i++ {
		for _, u := range up.Users {
			up.makePost(u.ID)
		}
	}
	return up
}

func (up *usersAndPosts) copy() *usersAndPosts {
	newUp := &usersAndPosts{
		postIndex: up.postIndex,
		userIndex: up.userIndex,
	}

	newUp.Users = append(newUp.Users, up.Users...)
	newUp.Posts = append(newUp.Posts, up.Posts...)
	return newUp
}

func (up *usersAndPosts) makeFeed(ids ...models.UserID) []models.Post {
	res := make([]models.Post, 0)
	idsMap := make(map[models.UserID]struct{})
	for _, id := range ids {
		idsMap[id] = struct{}{}
	}
	for id := range idsMap {
		for _, p := range up.Posts {
			if p.AuthorID == id {
				res = append(res, p)
			}
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].SequentialID >= res[j].SequentialID
	})
	return res
}

func (up *usersAndPosts) makeUser() models.User {
	u := models.User{
		ID:         models.UserID(fmt.Sprintf("id_%d", up.userIndex)),
		FirstName:  "first_name",
		SecondName: "second_name",
	}
	up.Users = append(up.Users, u)
	up.userIndex++
	return u
}

func (up *usersAndPosts) makePost(author models.UserID) models.Post {
	up.postIndex++
	p := models.Post{
		ID:           models.PostID(fmt.Sprintf("id_%d", up.postIndex)),
		SequentialID: up.postIndex,
		AuthorID:     author,
		Text:         fmt.Sprintf("post id_%d from %s", up.postIndex, author),
	}
	up.Posts = append(up.Posts, p)
	return p
}

func limitOffset(posts []models.Post, limit, offset int) []models.Post {
	if offset >= len(posts) {
		return nil
	}
	if offset+limit >= len(posts) {
		return posts[offset:]
	}
	return posts[offset:limit]
}

type lastMsg struct {
	Msg string
}

func (lm *lastMsg) Errorf(format string, args ...interface{}) {
	lm.Msg = fmt.Sprintf(format, args...)
}
