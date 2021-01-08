package models

type User struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Username  string
	Password  string // already hashed
}
