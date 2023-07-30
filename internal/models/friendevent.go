package models

const (
	FriendAddedEvent FriendEventType = iota
	FriendDeletedEvent
)

type FriendEventType int

type FriendEvent struct {
	Type     FriendEventType
	UserID   UserID
	FriendID UserID
}
