package models

import "errors"

var (
	UserNotFound = errors.New("user not found")
	PostNotFound = errors.New("post not found")
)
