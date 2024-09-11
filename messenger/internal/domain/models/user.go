package models

type User struct {
	ID       int64
	Email    string
	Username string
	PassHash []byte
	Phone    string
	Role     string
	Photo    []byte
	Active   bool
}
