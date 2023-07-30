package rabbitmq

import "github.com/ghkadim/highload_architect/internal/models"

type post struct {
	ID       models.PostID `json:"id"`
	Text     string        `json:"text"`
	AuthorID models.UserID `json:"authorId"`
}

type friendUpdate struct {
	Type     models.FriendEventType `json:"type"`
	UserID   models.UserID          `json:"userId"`
	FriendID models.UserID          `json:"friendId"`
}
