package models

type PostID string

type Post struct {
	SequentialID int64
	ID           PostID
	Text         string
	AuthorID     UserID
}
