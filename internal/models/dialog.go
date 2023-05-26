package models

type DialogMessage struct {
	From UserID
	To   UserID
	Text string
}
