package models

type User struct {
	ID           string
	FirstName    string
	SecondName   string
	Age          int32
	Biography    string
	City         string
	PasswordHash []byte
}
