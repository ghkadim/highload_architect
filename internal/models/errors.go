package models

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrPostNotFound = errors.New("post not found")
	ErrFeedNotFound = errors.New("feed not found")
	ErrFeedPartial  = errors.New("feed partial response")
)
