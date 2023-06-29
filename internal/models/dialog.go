package models

type DialogMessageID int64

type DialogMessage struct {
	ID   DialogMessageID
	From UserID
	To   UserID
	Text string
}
