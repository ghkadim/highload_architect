package models

type UserID string

type User struct {
	ID           UserID
	FirstName    string
	SecondName   string
	Age          *int32
	Biography    *string
	City         *string
	Password     *string
	PasswordHash []byte
}
